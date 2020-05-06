package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

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
	fs := http.FileServer(http.Dir("../../web/static"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	Title string
	Tags PairList
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	t, err  := template.ParseFiles("../../web/template/home.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.Execute(w, &Page{
		Title: title,
		Tags: tagsList,
	})
}