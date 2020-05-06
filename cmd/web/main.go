package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"golang.org/x/net/html"
	"sort"
	"time"
)

type Post struct {
	Body     []byte
	Images   []string
	Metadata PostMetadata
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
var posts []Post

func main() {
	tags := make(map[string]int)
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
					fmt.Println(htmlParseErr)
				}

				posts = append(posts, Post{
					Body: body,
					Images: images,
					Metadata: postMetadata,
				})
				for _, tag := range postMetadata.Tags {
					tags[tag] = tags[tag]+1
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
		// 2010-03-08 23:21:00 +0000 UTC
		const layout = "2006-01-02 15:04:05 -0700 MST"
		iTime, iErr := time.Parse(layout, posts[i].Metadata.CreatedOn)
		if iErr != nil {
			log.Fatal(iErr)
		}

		jTime, jErr := time.Parse(layout, posts[j].Metadata.CreatedOn)
		if jErr != nil {
			log.Fatal(jErr)
		}
		return iTime.Unix() > jTime.Unix()
	})
	fs := http.FileServer(http.Dir("../../web/static"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	Posts []Post
	Tags PairList
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, err  := template.ParseFiles("../../web/template/home.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.Execute(w, &Page{
		Posts: posts[:20],
		Tags: tagsList,
	})
}