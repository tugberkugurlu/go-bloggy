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
	"github.com/tugberkugurlu/go-bloggy/internal/blog"
)

var site *blog.Site

func main() {
	loadedSite, err := blog.LoadSite(blog.LoadSiteConfig{
		PostsDir: "../../web/posts",
	})
	if err != nil {
		log.Fatal(err)
	}
	site = loadedSite

	site.LayoutConfig = blog.LayoutConfig{
		AssetsUrl:                 site.Config.FullAssetsUrl(),
		GoogleAnalyticsTrackingId: os.Getenv("TUGBERKWEB_GoogleAnalytics__TrackingId"),
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
	if site.Config.IsLocal() {
		r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	}
	r.PathPrefix("/content/images/").HandlerFunc(legacyBlogImagesRedirector)
	r.Handle("/archive/{slug}", gziphandler.GzipHandler(http.HandlerFunc(blogPostPageHandler)))
	r.Handle("/tags/{tag_slug}", gziphandler.GzipHandler(http.HandlerFunc(tagsPageHandler)))
	r.Methods("GET").Path("/about").Handler(gziphandler.GzipHandler(http.HandlerFunc(aboutPageHandler)))
	r.Methods("GET").Path("/speaking").Handler(gziphandler.GzipHandler(http.HandlerFunc(speakingPageHandler)))
	r.Handle("/contact", gziphandler.GzipHandler(http.HandlerFunc(http.HandlerFunc(contactPageHandler))))
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

// CaselessMatcher is HTTP middleware that lowercases URL paths before routing.
func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func executeTemplate(w http.ResponseWriter, r *http.Request, templatePaths []string, data interface{}) {
	section := blog.SectionFromPath(r.URL.Path)
	err := blog.RenderPage(w, site.LayoutConfig, templatePaths, "../../web/template/layout.html", section, site.TagsList, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"../../web/template/about.html",
		"../../web/template/shared/carousel.html",
	}, blog.AboutPage{
		GeekTalksCarousel: blog.GetCarouselForTag("geek-talks", "Posts on Speaking", site.PostsByTagSlug),
		TopCarousel:       site.Carousels[0],
	})
}

func contactPageHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"../../web/template/contact.html",
		"../../web/template/shared/carousel.html",
	}, blog.ContactPage{
		ArchitectureCarousel: blog.GetCarouselForTag("architecture", "Top Picks for Architecture", site.PostsByTagSlug),
	})
}

func speakingPageHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"../../web/template/speaking.html",
		"../../web/template/shared/speaking-activity-card.html",
		"../../web/template/shared/carousel.html",
	}, blog.SpeakingPage{
		GeekTalksCarousel:  blog.GetCarouselForTag("geek-talks", "Posts on Speaking", site.PostsByTagSlug),
		SpeakingActivities: site.Speaking,
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

	for _, post := range site.Posts[:20] {
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
		return
	}
}

func tagsPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagSlug := vars["tag_slug"]
	posts, ok := site.PostsByTagSlug[tagSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	tag, ok := site.TagsBySlug[tagSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	executeTemplate(w, r, []string{
		"../../web/template/tag.html",
		"../../web/template/shared/post-item.html",
	}, blog.TagsPage{
		Posts: posts,
		Tag:   tag,
	})
}

func blogPostPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postSlug := vars["slug"]
	post, ok := site.PostsBySlug[postSlug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if !strings.EqualFold(post.Metadata.Slugs[0], postSlug) {
		http.Redirect(w, r, blog.GeneratePostURL(post), http.StatusMovedPermanently)
		return
	}

	executeTemplate(w, r, []string{
		"../../web/template/post.html",
		"../../web/template/shared/carousel.html",
	}, blog.PostPage{
		Post:                 post,
		RelatedPostsCarousel: blog.GetRelatedPostsCarousel(post, site.PostsByTagSlug),
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
		"../../web/template/home.html",
		"../../web/template/shared/carousel.html",
		"../../web/template/shared/post-item.html",
		"../../web/template/shared/speaking-activity-card.html",
	}, blog.Home{
		TopCarousel:        site.Carousels[0],
		Posts:              site.Posts[:3],
		SpeakingActivities: site.Speaking[:4],
	})
}

func blogHomeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, []string{
		"../../web/template/blog.html",
		"../../web/template/shared/post-item.html",
	}, blog.Blog{
		Posts: site.Posts,
	})
}
