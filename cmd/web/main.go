package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tugberkugurlu/go-bloggy/internal/blog"
)

// Config holds server-level configuration loaded from config.yaml.
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

var config Config
var layoutConfig blog.LayoutConfig
var index *blog.BlogIndex
var carousels []blog.Carousel

const templateRoot = "../../web/template"

func main() {
	var configErr error
	config, configErr = parseConfig()
	if configErr != nil {
		log.Fatal(configErr)
	}

	layoutConfig = blog.LayoutConfig{
		AssetsUrl:                 config.FullAssetsUrl(),
		GoogleAnalyticsTrackingId: os.Getenv("TUGBERKWEB_GoogleAnalytics__TrackingId"),
	}

	var loadErr error
	index, loadErr = blog.LoadIndex("../../web/posts")
	if loadErr != nil {
		log.Fatal(loadErr)
	}

	carousels = []blog.Carousel{blog.GetTopPicksCarousel(index.PostsByID)}

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
	r.Methods("GET").Path("/about").Handler(gziphandler.GzipHandler(http.HandlerFunc(aboutPageHandler)))
	r.Methods("GET").Path("/speaking").Handler(gziphandler.GzipHandler(http.HandlerFunc(speakingPageHandler)))
	r.Handle("/contact", gziphandler.GzipHandler(http.HandlerFunc(contactPageHandler)))
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

// templatePaths builds the full filesystem paths for a set of template file
// names relative to templateRoot, always appending the shared layout.
func templatePaths(files ...string) []string {
	paths := make([]string, len(files)+1)
	for i, f := range files {
		paths[i] = templateRoot + "/" + f
	}
	paths[len(files)] = templateRoot + "/layout.html"
	return paths
}

// executeTemplate is the HTTP-specific wrapper around blog.RenderPage.
func executeTemplate(w http.ResponseWriter, r *http.Request, files []string, data interface{}) {
	section := r.URL.Path[1:]
	if idx := strings.Index(section, "/"); idx != -1 {
		section = r.URL.Path[1 : idx+1]
	}
	if err := blog.RenderPage(w, section, templatePaths(files...), layoutConfig, index.TagsList, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"about.html",
		"shared/carousel.html",
	}, blog.AboutPage{
		GeekTalksCarousel: blog.GetCarouselForTag("geek-talks", "Posts on Speaking", index.PostsByTagSlug),
		TopCarousel:       carousels[0],
	})
}

func contactPageHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"contact.html",
		"shared/carousel.html",
	}, blog.ContactPage{
		ArchitectureCarousel: blog.GetCarouselForTag("architecture", "Top Picks for Architecture", index.PostsByTagSlug),
	})
}

func speakingPageHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"speaking.html",
		"shared/speaking-activity-card.html",
		"shared/carousel.html",
	}, blog.SpeakingPage{
		GeekTalksCarousel:  blog.GetCarouselForTag("geek-talks", "Posts on Speaking", index.PostsByTagSlug),
		SpeakingActivities: blog.SpeakingActivities,
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

	for _, post := range index.Posts[:20] {
		postLink := blog.GeneratePostURL(post)
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
	}
}

func tagsPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagSlug := vars["tag_slug"]
	posts, ok := index.PostsByTagSlug[tagSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}
	tag, ok := index.TagsBySlug[tagSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}
	executeTemplate(w, r, []string{
		"tag.html",
		"shared/post-item.html",
	}, blog.TagsPage{
		Posts: posts,
		Tag:   tag,
	})
}

func blogPostPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postSlug := vars["slug"]
	post, ok := index.PostsBySlug[postSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if !strings.EqualFold(post.Metadata.Slugs[0], postSlug) {
		http.Redirect(w, r, blog.GeneratePostURL(post), http.StatusMovedPermanently)
		return
	}
	executeTemplate(w, r, []string{
		"post.html",
		"shared/carousel.html",
	}, blog.PostPage{
		Post:                 post,
		RelatedPostsCarousel: blog.GetRelatedPostsCarousel(post, index.PostsByTagSlug),
		AdTags:               strings.Join(post.Metadata.Tags, ","),
	})
}

func legacyBlogImagesRedirector(w http.ResponseWriter, r *http.Request) {
	const legacyImagesRootPath = "/content/images"
	const newURIPrefix = "https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images"
	http.Redirect(w, r, fmt.Sprintf("%s%s", newURIPrefix, strings.ToLower(r.URL.Path[len(legacyImagesRootPath):])), http.StatusMovedPermanently)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"home.html",
		"shared/carousel.html",
		"shared/post-item.html",
		"shared/speaking-activity-card.html",
	}, blog.Home{
		TopCarousel:        carousels[0],
		Posts:              index.Posts[:3],
		SpeakingActivities: blog.SpeakingActivities[:4],
	})
}

func blogHomeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"blog.html",
		"shared/post-item.html",
	}, blog.Blog{
		Posts: index.Posts,
	})
}
