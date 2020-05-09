package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
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
	Body        template.HTML
	Images      []string
	PublishedOn time.Time
	Metadata    PostMetadata
	PublishedOnDisplay string
}

type PostMetadata struct {
	Title string `yaml:"title"`
	Abstract string `yaml:"abstract"`
	CreatedOn string `yaml:"created_at"`
	Tags []string `yaml:"tags"`
	Slugs []string `yaml:"slugs"`
}

type Pair struct {
	Key string
	Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

func rankByTagCount(tagFrequencies map[string]int) PairList{
	pl := make(PairList, len(tagFrequencies))
	i := 0
	for k, v := range tagFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

var tagsList PairList
var posts []*Post
var postsBySlug map[string]*Post

func main() {
	tags := make(map[string]int)
	postsBySlug = make(map[string]*Post)
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
					Body: template.HTML(string(body)),
					Images: images,
					Metadata: postMetadata,
					PublishedOn: publishedOn,
					PublishedOnDisplay: publishedOn.Format("2006-01-02 15:04:05"),
				}
				posts = append(posts, post)
				for _, tag := range postMetadata.Tags {
					tags[tag] = tags[tag]+1
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

	tagsList = rankByTagCount(tags)
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].PublishedOn.Unix() > posts[j].PublishedOn.Unix()
	})
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("../../web/static"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	r.PathPrefix("/content/images/").HandlerFunc(legacyBlogImagesRedirector)
	r.HandleFunc("/archive/{slug}", blogPostPageHandler)
	r.HandleFunc("/about", staticPage("about"))
	r.HandleFunc("/speaking", staticPage("speaking"))
	r.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", CaselessMatcher(r)))
}

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type Layout struct {
	Tags PairList
	Data interface{}
}

type Home struct {
	Posts []*Post
}

func blogPostPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, ok := postsBySlug[vars["slug"]]
	if !ok {
		http.NotFound(w, r)
		return
	}

	ExecuteTemplate(w, "../../web/template/post.html", &Layout{
		Tags: tagsList,
		Data: post,
	})
}

func legacyBlogImagesRedirector(w http.ResponseWriter, r *http.Request) {
	const legacyImagesRootPath = "/content/images"
	const newUriPrefix = "https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images"
	http.Redirect(w, r, fmt.Sprintf("%s%s", newUriPrefix, strings.ToLower(r.URL.Path[len(legacyImagesRootPath):])), http.StatusMovedPermanently)
}

func staticPage(page string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ExecuteTemplate(w, fmt.Sprintf("../../web/template/%s.html", page), &Layout{
			Tags: tagsList,
		})
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, "../../web/template/home.html", &Layout{
		Tags: tagsList,
		Data: Home{
			Posts: posts,
		},
	})
}

func ExecuteTemplate(w http.ResponseWriter, templatePath string, pageContext *Layout) {
	t, err  := template.ParseFiles(templatePath, "../../web/template/layout.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "layout", pageContext)
}