// cmd/staticgen generates a complete static website from the blog markdown
// sources. The output directory is suitable for direct upload to an S3 bucket
// configured for static-website hosting behind a CloudFront distribution.
//
// Usage:
//
//	staticgen [flags]
//
// Flags:
//
//	-posts-dir      Path to the directory containing blog post .md files.
//	                Default: ../../web/posts
//	-template-dir   Path to the directory containing Go HTML templates.
//	                Default: ../../web/template
//	-assets-src-dir Path to the local assets directory to copy into the output.
//	                Default: ../../web/static
//	-static-root-dir Path to root-level static files (robots.txt, favicon.ico).
//	                Default: ../../web/static-root
//	-output-dir     Destination directory for generated files. Created if absent.
//	                Default: dist
//	-assets-url     Value injected as AssetsUrl in templates (e.g. /assets or
//	                a CDN URL). Default: /assets
//
// URL structure of the generated output mirrors the live server exactly:
//
//	/                       → dist/index.html
//	/archive                → dist/archive/index.html
//	/archive/{slug}         → dist/archive/{slug}/index.html
//	/tags/{slug}            → dist/tags/{slug}/index.html
//	/about                  → dist/about/index.html
//	/speaking               → dist/speaking/index.html
//	/contact                → dist/contact/index.html
//	/feeds/rss              → dist/feeds/rss  (no extension; Content-Type set via S3 metadata)
//	/assets/*               → dist/assets/*   (copied verbatim)
//	/robots.txt, /favicon   → dist/ (copied from static-root)
//
// Redirects that the live server handles dynamically are written to
// dist/_redirects.json so they can be applied via a CloudFront Function or
// an S3 routing rule during deployment.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/tugberkugurlu/go-bloggy/internal/blog"
)

func main() {
	postsDir := flag.String("posts-dir", "../../web/posts", "path to posts directory")
	templateDir := flag.String("template-dir", "../../web/template", "path to templates directory")
	assetsSrcDir := flag.String("assets-src-dir", "../../web/static", "path to static assets directory")
	staticRootDir := flag.String("static-root-dir", "../../web/static-root", "path to static-root directory")
	outputDir := flag.String("output-dir", "dist", "output directory")
	assetsURL := flag.String("assets-url", "/assets", "base URL for static assets in templates")
	flag.Parse()

	if err := run(*postsDir, *templateDir, *assetsSrcDir, *staticRootDir, *outputDir, *assetsURL); err != nil {
		log.Fatalf("staticgen: %v", err)
	}
	log.Printf("staticgen: done — output written to %s", *outputDir)
}

func run(postsDir, templateDir, assetsSrcDir, staticRootDir, outputDir, assetsURL string) error {
	index, err := blog.LoadIndex(postsDir)
	if err != nil {
		return fmt.Errorf("loading posts: %w", err)
	}
	log.Printf("loaded %d posts", len(index.Posts))

	layoutConfig := blog.LayoutConfig{AssetsUrl: assetsURL}
	carousels := []blog.Carousel{blog.GetTopPicksCarousel(index.PostsByID)}

	gen := &generator{
		templateDir:  templateDir,
		outputDir:    outputDir,
		layoutConfig: layoutConfig,
		index:        index,
		carousels:    carousels,
	}

	if err := gen.generateAll(); err != nil {
		return err
	}

	if err := copyDir(assetsSrcDir, filepath.Join(outputDir, "assets")); err != nil {
		return fmt.Errorf("copying assets: %w", err)
	}
	if err := copyDir(staticRootDir, outputDir); err != nil {
		return fmt.Errorf("copying static-root: %w", err)
	}
	return nil
}

type generator struct {
	templateDir  string
	outputDir    string
	layoutConfig blog.LayoutConfig
	index        *blog.BlogIndex
	carousels    []blog.Carousel
}

func (g *generator) generateAll() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"home", g.generateHome},
		{"archive", g.generateArchive},
		{"posts", g.generatePosts},
		{"tags", g.generateTags},
		{"about", g.generateAbout},
		{"speaking", g.generateSpeaking},
		{"contact", g.generateContact},
		{"rss", g.generateRSS},
		{"redirects", g.generateRedirects},
	}
	for _, s := range steps {
		log.Printf("generating: %s", s.name)
		if err := s.fn(); err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
	}
	return nil
}

// renderPage renders a page and writes it to outPath, creating parent
// directories as needed. templateFiles are names relative to templateDir.
func (g *generator) renderPage(outPath, section string, templateFiles []string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	paths := make([]string, len(templateFiles)+1)
	for i, tf := range templateFiles {
		paths[i] = filepath.Join(g.templateDir, tf)
	}
	paths[len(templateFiles)] = filepath.Join(g.templateDir, "layout.html")

	return blog.RenderPage(f, section, paths, g.layoutConfig, g.index.TagsList, data)
}

func (g *generator) generateHome() error {
	return g.renderPage(
		filepath.Join(g.outputDir, "index.html"),
		"",
		[]string{
			"home.html",
			"shared/carousel.html",
			"shared/post-item.html",
			"shared/speaking-activity-card.html",
		},
		blog.Home{
			TopCarousel:        g.carousels[0],
			Posts:              g.index.Posts[:min(3, len(g.index.Posts))],
			SpeakingActivities: blog.SpeakingActivities[:min(4, len(blog.SpeakingActivities))],
		},
	)
}

func (g *generator) generateArchive() error {
	return g.renderPage(
		filepath.Join(g.outputDir, "archive", "index.html"),
		"archive",
		[]string{
			"blog.html",
			"shared/post-item.html",
		},
		blog.Blog{Posts: g.index.Posts},
	)
}

func (g *generator) generatePosts() error {
	for _, post := range g.index.Posts {
		slug := post.Metadata.Slugs[0]
		err := g.renderPage(
			filepath.Join(g.outputDir, "archive", slug, "index.html"),
			"archive",
			[]string{
				"post.html",
				"shared/carousel.html",
			},
			blog.PostPage{
				Post:                 post,
				RelatedPostsCarousel: blog.GetRelatedPostsCarousel(post, g.index.PostsByTagSlug),
				AdTags:               strings.Join(post.Metadata.Tags, ","),
			},
		)
		if err != nil {
			return fmt.Errorf("post %q: %w", slug, err)
		}
	}
	return nil
}

func (g *generator) generateTags() error {
	for tagSlug, tagPosts := range g.index.PostsByTagSlug {
		tag := g.index.TagsBySlug[tagSlug]
		err := g.renderPage(
			filepath.Join(g.outputDir, "tags", tagSlug, "index.html"),
			"tags",
			[]string{
				"tag.html",
				"shared/post-item.html",
			},
			blog.TagsPage{Tag: tag, Posts: tagPosts},
		)
		if err != nil {
			return fmt.Errorf("tag %q: %w", tagSlug, err)
		}
	}
	return nil
}

func (g *generator) generateAbout() error {
	return g.renderPage(
		filepath.Join(g.outputDir, "about", "index.html"),
		"about",
		[]string{
			"about.html",
			"shared/carousel.html",
		},
		blog.AboutPage{
			GeekTalksCarousel: blog.GetCarouselForTag("geek-talks", "Posts on Speaking", g.index.PostsByTagSlug),
			TopCarousel:       g.carousels[0],
		},
	)
}

func (g *generator) generateSpeaking() error {
	return g.renderPage(
		filepath.Join(g.outputDir, "speaking", "index.html"),
		"speaking",
		[]string{
			"speaking.html",
			"shared/speaking-activity-card.html",
			"shared/carousel.html",
		},
		blog.SpeakingPage{
			GeekTalksCarousel:  blog.GetCarouselForTag("geek-talks", "Posts on Speaking", g.index.PostsByTagSlug),
			SpeakingActivities: blog.SpeakingActivities,
		},
	)
}

func (g *generator) generateContact() error {
	return g.renderPage(
		filepath.Join(g.outputDir, "contact", "index.html"),
		"contact",
		[]string{
			"contact.html",
			"shared/carousel.html",
		},
		blog.ContactPage{
			ArchitectureCarousel: blog.GetCarouselForTag("architecture", "Top Picks for Architecture", g.index.PostsByTagSlug),
		},
	)
}

func (g *generator) generateRSS() error {
	author := &feeds.Author{Name: "Tugberk Ugurlu"}
	feed := &feeds.Feed{
		Title:       "Tugberk Ugurlu @ the Heart of Software",
		Link:        &feeds.Link{Href: "https://www.tugberkugurlu.com"},
		Description: "Software Engineer and Tech Product enthusiast Tugberk Ugurlu's home on the interwebs! Here, you can find out about Tugberk's conference talks, books and blog posts on software development techniques and practices",
		Author:      author,
	}
	for _, post := range g.index.Posts[:min(20, len(g.index.Posts))] {
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
		return err
	}

	outPath := filepath.Join(g.outputDir, "feeds", "rss")
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(rss), 0644)
}

// Redirect describes a single HTTP redirect rule for inclusion in the
// CloudFront Function or S3 routing configuration.
type Redirect struct {
	// From is the request path that should be redirected.
	From string `json:"from"`
	// To is the target URL (absolute) or path.
	To string `json:"to"`
	// Status is the HTTP status code (301 or 302).
	Status int `json:"status"`
}

// generateRedirects writes _redirects.json containing:
//   - Secondary-slug → canonical-slug redirects for all posts with aliases.
//   - The /content/images/* → Azure Blob Storage redirect rule.
func (g *generator) generateRedirects() error {
	var redirects []Redirect

	// Secondary slugs → canonical slug (301)
	for _, post := range g.index.Posts {
		if len(post.Metadata.Slugs) < 2 {
			continue
		}
		canonical := blog.GeneratePostURL(post)
		for _, altSlug := range post.Metadata.Slugs[1:] {
			redirects = append(redirects, Redirect{
				From:   "/archive/" + altSlug,
				To:     canonical,
				Status: 301,
			})
		}
	}

	// Legacy image URL redirect (wildcard; handled in CloudFront Function)
	redirects = append(redirects, Redirect{
		From:   "/content/images/*",
		To:     "https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images/*",
		Status: 301,
	})

	outPath := filepath.Join(g.outputDir, "_redirects.json")
	data, err := json.MarshalIndent(redirects, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, data, 0644)
}

// copyDir recursively copies src into dst, creating dst if it doesn't exist.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
