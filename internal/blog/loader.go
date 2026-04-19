package blog

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/russross/blackfriday/v2"
	"github.com/spf13/viper"
	"github.com/tugberkugurlu/go-bloggy/internal/htmltextextractor"
	"github.com/tugberkugurlu/go-bloggy/internal/readingtime"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

// LoadSiteConfig holds the paths needed to load a site.
type LoadSiteConfig struct {
	PostsDir   string
	// ConfigPath is the directory containing config.yaml. If empty, Viper
	// searches ../../ci/ and ../../ relative to the working directory (legacy
	// behavior from when cmd/web ran from cmd/web/).
	ConfigPath string
}

// LoadSite walks the posts directory, parses all markdown files, builds
// indexes, carousels, and returns a fully populated Site.
func LoadSite(cfg LoadSiteConfig) (*Site, error) {
	siteConfig, err := parseConfig(cfg.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	tagsBySlug := make(map[string]*Tag)
	postsByID := make(map[string]*Post)
	postsBySlug := make(map[string]*Post)
	postsByTagSlug := make(map[string][]*Post)
	var posts []*Post

	err = filepath.Walk(cfg.PostsDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		post, parseErr := parsePostFile(path)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", path, parseErr)
		}

		posts = append(posts, post)

		for _, tag := range post.Metadata.Tags {
			tagSlug := ToSlug(tag)
			t, ok := tagsBySlug[tagSlug]
			if !ok {
				t = &Tag{
					Name: tag,
				}
			}
			t.Count++
			tagsBySlug[tagSlug] = t
		}

		for _, slug := range post.Metadata.Slugs {
			postsBySlug[slug] = post
		}
		postsByID[post.Metadata.ID] = post

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking posts: %w", err)
	}

	tagsList := rankByTagCount(tagsBySlug)
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].PublishedOn.Unix() > posts[j].PublishedOn.Unix()
	})

	for _, post := range posts {
		var tags []TagCountPair
		for _, tag := range post.Metadata.Tags {
			tagSlug := ToSlug(tag)
			tags = append(tags, TagCountPair{
				Key:   tagSlug,
				Value: tagsBySlug[tagSlug],
			})
			postsByTagSlug[tagSlug] = append(postsByTagSlug[tagSlug], post)
		}
		post.Tags = tags
	}

	carousels := GetCarousels(postsByID)

	return &Site{
		Config:         siteConfig,
		Posts:          posts,
		PostsByID:      postsByID,
		PostsBySlug:    postsBySlug,
		PostsByTagSlug: postsByTagSlug,
		TagsBySlug:     tagsBySlug,
		TagsList:       tagsList,
		Carousels:      carousels,
		Speaking:       SpeakingActivities,
	}, nil
}

// parsePostFile reads and parses a single markdown post file.
func parsePostFile(markdownFilePath string) (*Post, error) {
	file, err := os.Open(markdownFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata []byte
	var body []byte
	var metadataStarted bool
	var metadataEnded bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !metadataStarted {
			if scanner.Text() != "---" {
				return nil, fmt.Errorf("'%s' doesn't have a valid yaml front matter", markdownFilePath)
			}
			metadataStarted = true
			continue
		}

		if scanner.Text() == "---" {
			metadataEnded = true
			continue
		}

		if metadataEnded {
			body = append(body, scanner.Bytes()...)
			body = append(body, []byte("\n")...)
		} else {
			metadata = append(metadata, scanner.Bytes()...)
			metadata = append(metadata, []byte("\n")...)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %s: %w", markdownFilePath, err)
	}

	var postMetadata PostMetadata
	if yamlErr := yaml.Unmarshal(metadata, &postMetadata); yamlErr != nil {
		return nil, fmt.Errorf("unmarshaling YAML in %s: %w", markdownFilePath, yamlErr)
	}

	if postMetadata.Format == "md" {
		body = blackfriday.Run(body, blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.TOC,
		})))
	}

	var images []string
	if document, htmlParseErr := html.Parse(bytes.NewReader(body)); htmlParseErr == nil {
		var crawler func(*html.Node)
		crawler = func(node *html.Node) {
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
				crawler(child)
			}
		}
		crawler(document)
	} else {
		return nil, fmt.Errorf("parsing HTML in %s: %w", markdownFilePath, htmlParseErr)
	}

	// 2010-03-08 23:21:00 +0000 UTC
	const layout = "2006-01-02 15:04:05 -0700 MST"
	publishedOn, parseErr := time.Parse(layout, postMetadata.CreatedOn)
	if parseErr != nil {
		return nil, fmt.Errorf("parsing date in %s: %w", markdownFilePath, parseErr)
	}

	var readingTime *time.Duration
	if texts, textExtractErr := htmltextextractor.Extract(body); textExtractErr == nil {
		rt, rtCalcErr := readingtime.Calculate(texts)
		if rtCalcErr == nil {
			readingTime = &rt
		}
	}

	readingTimeDisplay := ""
	if readingTime != nil {
		readingTimeDisplay = fmt.Sprintf("%d minutes read", int(readingTime.Minutes()))
	}

	post := &Post{
		Body:                    template.HTML(string(body)),
		Images:                  images,
		ReadingTime:             readingTime,
		ReadingTimeDisplay:      readingTimeDisplay,
		Highlight:               readingTimeDisplay,
		Metadata:                postMetadata,
		PublishedOn:             publishedOn,
		PublishedOnDisplayBrief: publishedOn.Format("2 January 2006"),
		PublishedOnDisplay:      publishedOn.Format("2006-01-02 15:04:05"),
	}

	return post, nil
}

func parseConfig(configPath string) (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath("../../ci/")
	v.AddConfigPath("../../")
	configErr := v.ReadInConfig()
	if configErr != nil {
		if _, ok := configErr.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("fatal error config file: %w", configErr)
		}
	}
	var assetsUrl string
	if viperVal := v.Get("assets_url"); viperVal == nil {
		return Config{}, fmt.Errorf("fatal error: 'assets_url' config value doesn't exist")
	} else {
		assetsUrl = viperVal.(string)
	}
	var assetsPrefix string
	if viperVal := v.Get("assets_prefix"); viperVal != nil {
		assetsPrefix = viperVal.(string)
	}
	return Config{
		AssetsUrl:    assetsUrl,
		AssetsPrefix: assetsPrefix,
	}, nil
}

// rankByTagCount sorts tags by count (descending).
func rankByTagCount(tagFrequencies map[string]*Tag) TagCountPairList {
	pl := make(TagCountPairList, len(tagFrequencies))
	i := 0
	for k, v := range tagFrequencies {
		pl[i] = TagCountPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
