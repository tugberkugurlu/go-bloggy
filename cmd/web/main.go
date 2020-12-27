package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday/v2"
	"github.com/spf13/viper"
	"github.com/tugberkugurlu/go-bloggy/internal/htmltextextractor"
	"github.com/tugberkugurlu/go-bloggy/internal/readingtime"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AssetsUrl    string
	AssetsPrefix string
}

func (c Config) IsLocal() bool {
	return strings.Index(c.AssetsUrl, "/") == 0
}

func (c Config) FullAssetsUrl() string {
	if c.AssetsPrefix == "" {
		return c.AssetsUrl
	}
	return fmt.Sprintf("%s/%s",
		strings.TrimSuffix(c.AssetsUrl, "/"),
		strings.TrimSuffix(c.AssetsPrefix, "/"),
	)
}

type Post struct {
	Body      template.HTML
	Highlight string
	Images    []string
	Tags      []TagCountPair
	Metadata  PostMetadata

	ReadingTimeDisplay string
	ReadingTime        *time.Duration

	PublishedOn             time.Time
	PublishedOnDisplay      string
	PublishedOnDisplayBrief string
}

type Tag struct {
	Name  string
	Count int
}

type PostMetadata struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Abstract  string   `yaml:"abstract"`
	Format    string   `yaml:"format"`
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

var config Config
var layoutConfig LayoutConfig
var tagsList TagCountPairList
var posts []*Post
var postsBySlug map[string]*Post
var postsByTagSlug map[string][]*Post
var tagsBySlug map[string]*Tag

func main() {
	var configErr error
	config, configErr = parseConfig()
	if configErr != nil {
		log.Fatal(configErr)
	}

	layoutConfig = LayoutConfig{
		AssetsUrl:                 config.FullAssetsUrl(),
		GoogleAnalyticsTrackingId: os.Getenv("TUGBERKWEB_GoogleAnalytics__TrackingId"),
	}

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
					log.Fatal(htmlParseErr)
				}

				// 2010-03-08 23:21:00 +0000 UTC
				const layout = "2006-01-02 15:04:05 -0700 MST"
				publishedOn, parseErr := time.Parse(layout, postMetadata.CreatedOn)
				if parseErr != nil {
					log.Fatal(parseErr)
				}

				var readingTime *time.Duration
				if texts, textExtractErr := htmltextextractor.Extract(body); textExtractErr == nil {
					rt, rtCalcErr := readingtime.Calculate(texts)
					if rtCalcErr == nil {
						readingTime = &rt
					}
				}

				readingTimeDisplay := func() string {
					if readingTime == nil {
						return ""
					}
					return fmt.Sprintf("%d minutes read", int(readingTime.Minutes()))
				}()
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
	if config.IsLocal() {
		r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	}
	r.PathPrefix("/content/images/").HandlerFunc(legacyBlogImagesRedirector)
	r.Handle("/archive/{slug}", gziphandler.GzipHandler(http.HandlerFunc(blogPostPageHandler)))
	r.Handle("/tags/{tag_slug}", gziphandler.GzipHandler(http.HandlerFunc(tagsPageHandler)))
	r.Handle("/about", gziphandler.GzipHandler(http.HandlerFunc(staticPage("about"))))
	r.Methods("GET").Path("/speaking").Handler(gziphandler.GzipHandler(http.HandlerFunc(speakingPageHandler)))
	r.Handle("/contact", gziphandler.GzipHandler(http.HandlerFunc(staticPage("contact"))))
	r.Handle("/archive", gziphandler.GzipHandler(http.HandlerFunc(blogHomeHandler)))
	r.HandleFunc("/feeds/rss", rssHandler)
	r.Handle("/", gziphandler.GzipHandler(http.HandlerFunc(homeHandler)))

	rootFs := http.FileServer(http.Dir("../../web/static-root"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", rootFs))

	serverPortStr := os.Getenv("SERVER_PORT")
	_, parseErr := strconv.ParseInt(serverPortStr, 10, 16)
	if parseErr != nil {
		serverPortStr = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPortStr), CaselessMatcher(r)))
}

func parseConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../ci/")
	viper.AddConfigPath("../../")
	configErr := viper.ReadInConfig()
	if configErr != nil {
		if _, ok := configErr.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("Fatal error config file: %s \n", configErr)
		}
	}
	var assetsUrl string
	if viperVal := viper.Get("assets_url"); viperVal == nil {
		return Config{}, errors.New("Fatal error: 'assets_url' config value doesn't exist")
	} else {
		assetsUrl = viperVal.(string)
	}
	var assetsPrefix string
	if viperVal := viper.Get("assets_prefix"); viperVal != nil {
		assetsPrefix = viperVal.(string)
	}
	return Config{
		AssetsUrl:    assetsUrl,
		AssetsPrefix: assetsPrefix,
	}, nil
}

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type PageData struct {
	Data   interface{}
	Config LayoutConfig
}

type Layout struct {
	Title       string
	Description string
	Tags        TagCountPairList
	Section     string
	Data        PageData
	AdTags      string
	Config      LayoutConfig
}

type LayoutConfig struct {
	GoogleAnalyticsTrackingId string
	AssetsUrl                 string
}

type Carousel struct {
	Title string
	Posts []*Post
}

type Home struct {
	TopCarousel        Carousel
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

func (t TagsPage) Title() string {
	return t.Tag.Name
}

func (t TagsPage) Description() string {
	baseTitle := fmt.Sprintf("Tugberk's thoughts on the topic of '%s'", t.Tag.Name)
	if t.Tag.Count > 1 {
		return fmt.Sprintf("%s, through %d blog posts", baseTitle, t.Tag.Count)
	}
	return baseTitle
}

type PostPage struct {
	Post   *Post
	AdTags string
}

func (p PostPage) Title() string {
	return p.Post.Metadata.Title
}

func (p PostPage) Description() string {
	return p.Post.Metadata.Abstract
}

type SpeakingPage struct {
	SpeakingActivities []*SpeakingActivity
}

func (s SpeakingPage) Title() string {
	return "Tugberk Ugurlu Public Speaking Engagements"
}

func (s SpeakingPage) Description() string {
	return "Tugberk speaks at conferences on technical leadership, software architecture, lean software development, Microsoft Azure and .NET."
}

var _ Page = (*SpeakingPage)(nil)
var _ Page = (*PostPage)(nil)
var _ Page = (*TagsPage)(nil)

type Page interface {
	Title() string
	Description() string
}

func speakingPageHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, r, layoutConfig, []string{
		"../../web/template/speaking.html",
		"../../web/template/shared/speaking-activity-card.html",
	}, SpeakingPage{
		SpeakingActivities: speakingActivities,
	})
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	author := &feeds.Author{Name: "Tugberk Ugurlu"}
	feed := &feeds.Feed{
		Title:       "Tugberk Ugurlu @ the Heart of Software",
		Link:        &feeds.Link{Href: "https://www.tugberkugurlu.com"},
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
	postLink := fmt.Sprintf("https://www.tugberkugurlu.com/archive/%s", post.Metadata.Slugs[0])
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

	ExecuteTemplate(w, r, layoutConfig, []string{
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

	ExecuteTemplate(w, r, layoutConfig, []string{"../../web/template/post.html"}, PostPage{
		Post:   post,
		AdTags: strings.Join(post.Metadata.Tags, ","),
	})
}

func legacyBlogImagesRedirector(w http.ResponseWriter, r *http.Request) {
	const legacyImagesRootPath = "/content/images"
	const newURIPrefix = "https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images"
	http.Redirect(w, r, fmt.Sprintf("%s%s", newURIPrefix, strings.ToLower(r.URL.Path[len(legacyImagesRootPath):])), http.StatusMovedPermanently)
}

func staticPage(page string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ExecuteTemplate(w, r, layoutConfig, []string{fmt.Sprintf("../../web/template/%s.html", page)}, nil)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, r, layoutConfig, []string{
		"../../web/template/home.html",
		"../../web/template/shared/carousel.html",
		"../../web/template/shared/post-item.html",
		"../../web/template/shared/speaking-activity-card.html",
	}, Home{
		TopCarousel: Carousel{
			Title: "Posts on ASP.NET",
			Posts: func() []*Post {
				input := postsByTagSlug["asp-net"]
				iteration := int(math.Min(float64(20), float64(len(input))))
				result := make([]*Post, 0, iteration)
				for i := 0; i < iteration; i++ {
					p := input[i]
					if len(p.Images) == 0 {
						continue
					}
					result = append(result, p)
				}
				return result
			}(),
		},
		Posts:              posts[:3],
		SpeakingActivities: speakingActivities[:4],
	})
}

func blogHomeHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, r, layoutConfig, []string{
		"../../web/template/blog.html",
		"../../web/template/shared/post-item.html",
	}, Blog{
		Posts: posts,
	})
}

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, config LayoutConfig, templatePaths []string, data interface{}) {
	t := template.New("")
	t = t.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	t, err := t.ParseFiles(append(templatePaths, "../../web/template/layout.html")...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	section := r.URL.Path[1:]
	index := strings.Index(section, "/")
	if index != -1 {
		section = r.URL.Path[1 : index+1]
	}

	pageTitle := "Tugberk @ the Heart of Software"
	pageDescription := "Software Engineer and Tech Product enthusiast Tugberk Ugurlu's home on the interwebs! Here, you can find out about Tugberk's conference talks, books and blog posts on software development techniques and practices."
	page, ok := data.(Page)
	if ok {
		if page.Title() != "" {
			pageTitle = fmt.Sprintf("%s | %s", page.Title(), pageTitle)
		}
		if page.Description() != "" {
			pageDescription = page.Description()
		}
	}

	adTags := "software development, asp.net, aws, azure, sql server, dynamodb, elasticsearch, mongodb, .net"
	postPage, ok := data.(PostPage)
	if ok {
		adTags = postPage.AdTags
	}

	pageContext := &Layout{
		Title:       pageTitle,
		Description: pageDescription,
		Tags:        tagsList,
		Section:     section,
		AdTags:      adTags,
		Config:      config,

		Data: PageData{
			Data:   data,
			Config: config,
		},
	}
	templateErr := t.ExecuteTemplate(w, "layout", pageContext)
	if templateErr != nil {
		http.Error(w, templateErr.Error(), http.StatusInternalServerError)
	}
}
