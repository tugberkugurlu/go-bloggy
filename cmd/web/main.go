package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Post struct {
	Body               template.HTML
	Images             []string
	PublishedOn        time.Time
	Tags               []TagCountPair
	Metadata           PostMetadata
	PublishedOnDisplay string
}

type Tag struct {
	Name  string
	Count int
}

type PostMetadata struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Abstract  string   `yaml:"abstract"`
	CreatedOn string   `yaml:"created_at"`
	Tags      []string `yaml:"tags"`
	Slugs     []string `yaml:"slugs"`
}

type TagCountPair struct {
	Key   string
	Value *Tag
}

type TagCountPairList []TagCountPair

func (p TagCountPairList) Len() int           { return len(p) }
func (p TagCountPairList) Less(i, j int) bool { return p[i].Value.Count < p[j].Value.Count }
func (p TagCountPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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

var tagsList TagCountPairList
var posts []*Post
var postsBySlug map[string]*Post
var postsByTagSlug map[string][]*Post
var tagsBySlug map[string]*Tag

func main() {
	tagsBySlug = make(map[string]*Tag)
	postsBySlug = make(map[string]*Post)
	postsByTagSlug = make(map[string][]*Post)
	err := filepath.Walk("../../web/posts", func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			func(markdownFilePath string) {
				file, err := os.Open(markdownFilePath)
				if err != nil {
					log.Fatal(err)
				}
				defer func() {
					closeErr := file.Close()
					if closeErr != nil {
						log.Fatal(closeErr)
					}
				}()

				var metadata []byte
				var body []byte
				var metadataStarted bool
				var metadataEnded bool
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					if !metadataStarted {
						if scanner.Text() != "---" {
							log.Fatalf("'%s' doesn't have a valid yaml front matter", markdownFilePath)
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

				var postMetadata PostMetadata
				yamlErr := yaml.Unmarshal(metadata, &postMetadata)
				if yamlErr != nil {
					log.Fatal(yamlErr)
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
					log.Fatal(htmlParseErr)
				}

				// 2010-03-08 23:21:00 +0000 UTC
				const layout = "2006-01-02 15:04:05 -0700 MST"
				publishedOn, parseErr := time.Parse(layout, postMetadata.CreatedOn)
				if parseErr != nil {
					log.Fatal(parseErr)
				}

				post := &Post{
					Body:               template.HTML(string(body)),
					Images:             images,
					Metadata:           postMetadata,
					PublishedOn:        publishedOn,
					PublishedOnDisplay: publishedOn.Format("2006-01-02 15:04:05"),
				}
				posts = append(posts, post)
				for _, tag := range postMetadata.Tags {
					tagSlug := slug.Make(tag)
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
			}(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	tagsList = rankByTagCount(tagsBySlug)
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].PublishedOn.Unix() > posts[j].PublishedOn.Unix()
	})

	for _, post := range posts {
		var tags []TagCountPair
		for _, tag := range post.Metadata.Tags {
			tagSlug := slug.Make(tag)
			tags = append(tags, TagCountPair{
				Key:   tagSlug,
				Value: tagsBySlug[tagSlug],
			})
			postsByTagSlug[tagSlug] = append(postsByTagSlug[tagSlug], post)
		}
		sort.SliceStable(tags, func(i, j int) bool {
			return tags[i].Value.Count > tags[j].Value.Count
		})
		post.Tags = tags
	}

	r := mux.NewRouter()
	r.Host("tugberkugurlu.com").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := r.URL.Scheme
		if scheme == "" {
			scheme = "http"
		}
		url := fmt.Sprintf("%s://%s%s", scheme, "www.tugberkugurlu.com", r.URL.RequestURI())
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	fs := http.FileServer(http.Dir("../../web/static"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	r.PathPrefix("/content/images/").HandlerFunc(legacyBlogImagesRedirector)
	r.HandleFunc("/archive/{slug}", blogPostPageHandler)
	r.HandleFunc("/tags/{tag_slug}", tagsPageHandler)
	r.HandleFunc("/about", staticPage("about"))
	r.Methods("GET").Path("/speaking").HandlerFunc(speakingPageHandler)
	r.HandleFunc("/contact", staticPage("contact"))
	r.HandleFunc("/archive", blogHomeHandler)
	r.HandleFunc("/feeds/rss", rssHandler)
	r.HandleFunc("/", homeHandler)

	rootFs := http.FileServer(http.Dir("../../web/static-root"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", rootFs))

	log.Fatal(http.ListenAndServe(":8080", CaselessMatcher(r)))
}

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type Layout struct {
	Tags    TagCountPairList
	Section string
	Data    interface{}
}

type Home struct {
	Posts              []*Post
	SpeakingActivities []*SpeakingActivity
}

type Blog struct {
	Posts []*Post
}

type TagsPage struct {
	Tag   *Tag
	Posts []*Post
}

func speakingPageHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, r, []string{
		"../../web/template/speaking.html",
		"../../web/template/shared/speaking-activity-card.html",
	}, speakingActivities)
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	author := &feeds.Author{Name: "Tugberk Ugurlu"}
	feed := &feeds.Feed{
		Title:       "Tugberk Ugurlu @ the Heart of Software",
		Link:        &feeds.Link{Href: "http://tugberkugurlu.com"},
		Description: "Software Engineer and Tech Product enthusiast Tugberk Ugurlu's home on the interwebs! Here, you can find out about Tugberk's conference talks, books and blog posts on software development techniques and practices",
		Author:      author,
	}

	for _, post := range posts[:20] {
		postLink := generatePostURL(post)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       post.Metadata.Title,
			Description: string(post.Body),
			Created:     post.PublishedOn,
			Author:      author,
			Link:        &feeds.Link{Href: postLink},
			Id:          postLink,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public,max-age=900")

	_, err = w.Write([]byte(rss))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generatePostURL(post *Post) string {
	postLink := fmt.Sprintf("http://tugberkugurlu.com/archive/%s", post.Metadata.Slugs[0])
	return postLink
}

func tagsPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagSlug := vars["tag_slug"]
	posts, ok := postsByTagSlug[tagSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	tag, ok := tagsBySlug[tagSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	ExecuteTemplate(w, r, []string{
		"../../web/template/tag.html",
		"../../web/template/shared/post-item.html",
	}, TagsPage{
		Posts: posts,
		Tag:   tag,
	})
}

func blogPostPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postSlug := vars["slug"]
	post, ok := postsBySlug[postSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if !strings.EqualFold(post.Metadata.Slugs[0], postSlug) {
		http.Redirect(w, r, generatePostURL(post), http.StatusMovedPermanently)
		return
	}

	ExecuteTemplate(w, r, []string{"../../web/template/post.html"}, post)
}

func legacyBlogImagesRedirector(w http.ResponseWriter, r *http.Request) {
	const legacyImagesRootPath = "/content/images"
	const newUriPrefix = "https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images"
	http.Redirect(w, r, fmt.Sprintf("%s%s", newUriPrefix, strings.ToLower(r.URL.Path[len(legacyImagesRootPath):])), http.StatusMovedPermanently)
}

func staticPage(page string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ExecuteTemplate(w, r, []string{fmt.Sprintf("../../web/template/%s.html", page)}, nil)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, r, []string{
		"../../web/template/home.html",
		"../../web/template/shared/post-item.html",
		"../../web/template/shared/speaking-activity-card.html",
	}, Home{
		Posts:              posts[:3],
		SpeakingActivities: speakingActivities[:4],
	})
}

func blogHomeHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, r, []string{
		"../../web/template/blog.html",
		"../../web/template/shared/post-item.html",
	}, Blog{
		Posts: posts,
	})
}

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, templatePaths []string, data interface{}) {
	t, err := template.ParseFiles(append(templatePaths, "../../web/template/layout.html")...)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	section := r.URL.Path[1:]
	index := strings.Index(section, "/")
	if index != -1 {
		section = r.URL.Path[1 : index+1]
	}

	pageContext := &Layout{
		Tags:    tagsList,
		Data:    data,
		Section: section,
	}
	t.ExecuteTemplate(w, "layout", pageContext)
}
