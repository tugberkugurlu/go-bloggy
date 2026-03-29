package blog

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/russross/blackfriday/v2"
	"github.com/tugberkugurlu/go-bloggy/internal/htmltextextractor"
	"github.com/tugberkugurlu/go-bloggy/internal/readingtime"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

// preDefinedSlugs maps tag names that would produce an incorrect slug via the
// automatic slugifier (e.g. "C#" would become just "c" without the override).
var preDefinedSlugs = map[string]string{
	"c#":  "c-sharp",
	"c++": "cpp",
}

// ToSlug converts a tag name to a URL-safe slug, applying any predefined
// overrides before running the automatic slugifier.
func ToSlug(tag string) string {
	if v, ok := preDefinedSlugs[strings.ToLower(tag)]; ok {
		return slug.Make(v)
	}
	return slug.Make(tag)
}

// GeneratePostURL returns the canonical absolute URL for a post.
func GeneratePostURL(post *Post) string {
	return fmt.Sprintf("https://www.tugberkugurlu.com/archive/%s", post.Metadata.Slugs[0])
}

// BlogIndex contains all in-memory lookup structures built from the posts
// directory. All fields are populated by LoadIndex.
type BlogIndex struct {
	// Posts is all posts sorted newest-first.
	Posts []*Post
	// PostsByID enables O(1) lookup by the post's unique ID field.
	PostsByID map[string]*Post
	// PostsBySlug enables O(1) lookup by any slug (primary or alternate).
	PostsBySlug map[string]*Post
	// PostsByTagSlug maps each tag slug to the posts that carry that tag,
	// in newest-first order.
	PostsByTagSlug map[string][]*Post
	// TagsBySlug maps each tag slug to its Tag (name + count).
	TagsBySlug map[string]*Tag
	// TagsList is all tags sorted by post count descending, used for the
	// sidebar tag cloud.
	TagsList TagCountPairList
}

// publishedOnLayout is the time format used in post front-matter created_at
// fields (after YAML parsing normalises the value to a Go time string).
const publishedOnLayout = "2006-01-02 15:04:05 -0700 MST"

// LoadIndex walks postsDir, parses every .md file it finds, and returns a
// fully-populated BlogIndex. postsDir may contain arbitrary sub-directories;
// only files with a .md extension are processed.
func LoadIndex(postsDir string) (*BlogIndex, error) {
	tagsBySlug := make(map[string]*Tag)
	var rawPosts []*Post

	err := filepath.Walk(postsDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		post, parseErr := parsePost(path)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", path, parseErr)
		}
		rawPosts = append(rawPosts, post)
		for _, tag := range post.Metadata.Tags {
			ts := ToSlug(tag)
			t, ok := tagsBySlug[ts]
			if !ok {
				t = &Tag{Name: tag}
			}
			t.Count++
			tagsBySlug[ts] = t
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort newest-first before building per-tag slices so that tag pages
	// and carousels are already in the correct order.
	sort.SliceStable(rawPosts, func(i, j int) bool {
		return rawPosts[i].PublishedOn.Unix() > rawPosts[j].PublishedOn.Unix()
	})

	tagsList := rankByTagCount(tagsBySlug)

	postsByID := make(map[string]*Post, len(rawPosts))
	postsBySlug := make(map[string]*Post, len(rawPosts))
	postsByTagSlug := make(map[string][]*Post)

	for _, post := range rawPosts {
		var tags []TagCountPair
		for _, tag := range post.Metadata.Tags {
			ts := ToSlug(tag)
			tags = append(tags, TagCountPair{Key: ts, Value: tagsBySlug[ts]})
			postsByTagSlug[ts] = append(postsByTagSlug[ts], post)
		}
		post.Tags = tags
		postsByID[post.Metadata.ID] = post
		for _, s := range post.Metadata.Slugs {
			postsBySlug[s] = post
		}
	}

	return &BlogIndex{
		Posts:          rawPosts,
		PostsByID:      postsByID,
		PostsBySlug:    postsBySlug,
		PostsByTagSlug: postsByTagSlug,
		TagsBySlug:     tagsBySlug,
		TagsList:       tagsList,
	}, nil
}

func rankByTagCount(tagFrequencies map[string]*Tag) TagCountPairList {
	pl := make(TagCountPairList, 0, len(tagFrequencies))
	for k, v := range tagFrequencies {
		pl = append(pl, TagCountPair{k, v})
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// parsePost reads a single markdown file and returns a populated Post.
func parsePost(markdownFilePath string) (*Post, error) {
	file, err := os.Open(markdownFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadataBytes, bodyBytes []byte
	var metadataStarted, metadataEnded bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !metadataStarted {
			if line != "---" {
				return nil, fmt.Errorf("missing YAML front-matter opening delimiter")
			}
			metadataStarted = true
			continue
		}
		if line == "---" {
			metadataEnded = true
			continue
		}
		if metadataEnded {
			bodyBytes = append(bodyBytes, scanner.Bytes()...)
			bodyBytes = append(bodyBytes, '\n')
		} else {
			metadataBytes = append(metadataBytes, scanner.Bytes()...)
			metadataBytes = append(metadataBytes, '\n')
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return nil, scanErr
	}

	var meta PostMetadata
	if err := yaml.Unmarshal(metadataBytes, &meta); err != nil {
		return nil, fmt.Errorf("unmarshalling front matter: %w", err)
	}

	if meta.Format == "md" {
		bodyBytes = blackfriday.Run(bodyBytes, blackfriday.WithRenderer(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				Flags: blackfriday.TOC,
			}),
		))
	}

	images := extractImages(bodyBytes)

	publishedOn, err := time.Parse(publishedOnLayout, meta.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at %q: %w", meta.CreatedOn, err)
	}

	var readingTime *time.Duration
	if texts, extractErr := htmltextextractor.Extract(bodyBytes); extractErr == nil {
		if rt, calcErr := readingtime.Calculate(texts); calcErr == nil {
			readingTime = &rt
		}
	}

	readingTimeDisplay := ""
	if readingTime != nil {
		readingTimeDisplay = fmt.Sprintf("%d minutes read", int(readingTime.Minutes()))
	}

	return &Post{
		Body:                    template.HTML(string(bodyBytes)),
		Images:                  images,
		ReadingTime:             readingTime,
		ReadingTimeDisplay:      readingTimeDisplay,
		Highlight:               readingTimeDisplay,
		Metadata:                meta,
		PublishedOn:             publishedOn,
		PublishedOnDisplayBrief: publishedOn.Format("2 January 2006"),
		PublishedOnDisplay:      publishedOn.Format("2006-01-02 15:04:05"),
	}, nil
}

// extractImages walks the HTML node tree and returns every <img src="..."> value found.
func extractImages(htmlBody []byte) []string {
	var images []string
	doc, err := html.Parse(bytes.NewReader(htmlBody))
	if err != nil {
		return images
	}
	var crawl func(*html.Node)
	crawl = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "img" {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					images = append(images, attr.Val)
					break
				}
			}
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawl(child)
		}
	}
	crawl(doc)
	return images
}
