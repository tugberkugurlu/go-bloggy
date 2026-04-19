package main

import (
	"crypto/sha256"
	"encoding/hex"
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

// ManifestPage represents a generated page in the manifest.
type ManifestPage struct {
	Path        string `json:"path"`
	File        string `json:"file"`
	ContentHash string `json:"content_hash"`
}

// ManifestRedirect represents a slug redirect in the manifest.
type ManifestRedirect struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Status int    `json:"status"`
}

// ManifestPattern represents a URL pattern redirect in the manifest.
type ManifestPattern struct {
	Pattern        string `json:"pattern"`
	RedirectPrefix string `json:"redirect_prefix"`
	Status         int    `json:"status"`
}

// Manifest is the build artifact listing all pages and redirects.
type Manifest struct {
	Pages     []ManifestPage     `json:"pages"`
	Redirects []ManifestRedirect `json:"redirects"`
	Patterns  []ManifestPattern  `json:"patterns"`
}

func main() {
	postsDir := flag.String("posts", "./web/posts", "Path to blog posts directory")
	templatesDir := flag.String("templates", "./web/template", "Path to templates directory")
	staticDir := flag.String("static", "./web/static", "Path to static assets directory")
	staticRootDir := flag.String("static-root", "./web/static-root", "Path to static root files directory")
	configPath := flag.String("config", ".", "Path to directory containing config.yaml")
	outputDir := flag.String("output", "./output", "Path to output directory")
	flag.Parse()

	if err := run(*postsDir, *templatesDir, *staticDir, *staticRootDir, *configPath, *outputDir); err != nil {
		log.Fatal(err)
	}
}

func run(postsDir, templatesDir, staticDir, staticRootDir, configPath, outputDir string) error {
	site, err := blog.LoadSite(blog.LoadSiteConfig{
		PostsDir:   postsDir,
		ConfigPath: configPath,
	})
	if err != nil {
		return fmt.Errorf("loading site: %w", err)
	}

	site.LayoutConfig = blog.LayoutConfig{
		AssetsUrl: site.Config.FullAssetsUrl(),
	}

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("cleaning output directory: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	g := &generator{
		site:         site,
		templatesDir: templatesDir,
		outputDir:    outputDir,
	}

	// Generate all pages
	log.Println("Generating home page...")
	if err := g.generateHome(); err != nil {
		return fmt.Errorf("generating home: %w", err)
	}

	log.Println("Generating archive page...")
	if err := g.generateArchive(); err != nil {
		return fmt.Errorf("generating archive: %w", err)
	}

	log.Printf("Generating %d post pages...\n", len(site.Posts))
	if err := g.generatePosts(); err != nil {
		return fmt.Errorf("generating posts: %w", err)
	}

	log.Printf("Generating %d tag pages...\n", len(site.TagsBySlug))
	if err := g.generateTags(); err != nil {
		return fmt.Errorf("generating tags: %w", err)
	}

	log.Println("Generating about page...")
	if err := g.generateAbout(); err != nil {
		return fmt.Errorf("generating about: %w", err)
	}

	log.Println("Generating speaking page...")
	if err := g.generateSpeaking(); err != nil {
		return fmt.Errorf("generating speaking: %w", err)
	}

	log.Println("Generating contact page...")
	if err := g.generateContact(); err != nil {
		return fmt.Errorf("generating contact: %w", err)
	}

	log.Println("Generating RSS feed...")
	if err := g.generateRSS(); err != nil {
		return fmt.Errorf("generating RSS: %w", err)
	}

	log.Println("Generating 404 page...")
	if err := g.generate404(); err != nil {
		return fmt.Errorf("generating 404: %w", err)
	}

	// Copy static assets
	log.Println("Copying static assets...")
	if err := copyDir(staticDir, filepath.Join(outputDir, "assets")); err != nil {
		return fmt.Errorf("copying static assets: %w", err)
	}

	// Copy static root files
	log.Println("Copying static root files...")
	if err := copyDir(staticRootDir, outputDir); err != nil {
		return fmt.Errorf("copying static root files: %w", err)
	}

	// Build and write manifest
	log.Println("Writing manifest...")
	if err := g.writeManifest(); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	log.Printf("Generation complete: %d pages, %d redirects\n", len(g.manifest.Pages), len(g.manifest.Redirects))
	return nil
}

type generator struct {
	site         *blog.Site
	templatesDir string
	outputDir    string
	manifest     Manifest
}

func (g *generator) templatePath(name string) string {
	return filepath.Join(g.templatesDir, name)
}

func (g *generator) renderAndWrite(urlPath, filePath string, templatePaths []string, section string, data interface{}) error {
	absFilePath := filepath.Join(g.outputDir, filePath)
	if err := os.MkdirAll(filepath.Dir(absFilePath), 0755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", filePath, err)
	}

	f, err := os.Create(absFilePath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", filePath, err)
	}
	defer f.Close()

	layoutPath := g.templatePath("layout.html")
	if err := blog.RenderPage(f, g.site.LayoutConfig, templatePaths, layoutPath, section, g.site.TagsList, data); err != nil {
		return fmt.Errorf("rendering %s: %w", urlPath, err)
	}

	// Compute content hash
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seeking %s: %w", filePath, err)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return fmt.Errorf("hashing %s: %w", filePath, err)
	}
	contentHash := hex.EncodeToString(hash.Sum(nil))[:16]

	g.manifest.Pages = append(g.manifest.Pages, ManifestPage{
		Path:        urlPath,
		File:        filePath,
		ContentHash: contentHash,
	})

	return nil
}

func (g *generator) generateHome() error {
	data := blog.Home{
		TopCarousel:        g.site.Carousels[0],
		Posts:              g.site.Posts[:min(3, len(g.site.Posts))],
		SpeakingActivities: g.site.Speaking[:min(4, len(g.site.Speaking))],
	}
	return g.renderAndWrite("/", "index.html", []string{
		g.templatePath("home.html"),
		g.templatePath("shared/carousel.html"),
		g.templatePath("shared/post-item.html"),
		g.templatePath("shared/speaking-activity-card.html"),
	}, "", data)
}

func (g *generator) generateArchive() error {
	data := blog.Blog{
		Posts: g.site.Posts,
	}
	return g.renderAndWrite("/archive", "archive/index.html", []string{
		g.templatePath("blog.html"),
		g.templatePath("shared/post-item.html"),
	}, "archive", data)
}

func (g *generator) generatePosts() error {
	for _, post := range g.site.Posts {
		canonicalSlug := post.Metadata.Slugs[0]
		urlPath := fmt.Sprintf("/archive/%s", canonicalSlug)
		filePath := fmt.Sprintf("archive/%s/index.html", canonicalSlug)

		data := blog.PostPage{
			Post:                 post,
			RelatedPostsCarousel: blog.GetRelatedPostsCarousel(post, g.site.PostsByTagSlug),
			AdTags:               strings.Join(post.Metadata.Tags, ","),
		}
		if err := g.renderAndWrite(urlPath, filePath, []string{
			g.templatePath("post.html"),
			g.templatePath("shared/carousel.html"),
		}, "archive", data); err != nil {
			return err
		}

		// Register non-canonical slug redirects
		for _, slug := range post.Metadata.Slugs[1:] {
			g.manifest.Redirects = append(g.manifest.Redirects, ManifestRedirect{
				From:   fmt.Sprintf("/archive/%s", slug),
				To:     urlPath,
				Status: 301,
			})
		}
	}
	return nil
}

func (g *generator) generateTags() error {
	for tagSlug, tag := range g.site.TagsBySlug {
		posts := g.site.PostsByTagSlug[tagSlug]
		urlPath := fmt.Sprintf("/tags/%s", tagSlug)
		filePath := fmt.Sprintf("tags/%s/index.html", tagSlug)

		data := blog.TagsPage{
			Posts: posts,
			Tag:   tag,
		}
		if err := g.renderAndWrite(urlPath, filePath, []string{
			g.templatePath("tag.html"),
			g.templatePath("shared/post-item.html"),
		}, "tags", data); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) generateAbout() error {
	data := blog.AboutPage{
		GeekTalksCarousel: blog.GetCarouselForTag("geek-talks", "Posts on Speaking", g.site.PostsByTagSlug),
		TopCarousel:       g.site.Carousels[0],
	}
	return g.renderAndWrite("/about", "about/index.html", []string{
		g.templatePath("about.html"),
		g.templatePath("shared/carousel.html"),
	}, "about", data)
}

func (g *generator) generateSpeaking() error {
	data := blog.SpeakingPage{
		GeekTalksCarousel:  blog.GetCarouselForTag("geek-talks", "Posts on Speaking", g.site.PostsByTagSlug),
		SpeakingActivities: g.site.Speaking,
	}
	return g.renderAndWrite("/speaking", "speaking/index.html", []string{
		g.templatePath("speaking.html"),
		g.templatePath("shared/speaking-activity-card.html"),
		g.templatePath("shared/carousel.html"),
	}, "speaking", data)
}

func (g *generator) generateContact() error {
	data := blog.ContactPage{
		ArchitectureCarousel: blog.GetCarouselForTag("architecture", "Top Picks for Architecture", g.site.PostsByTagSlug),
	}
	return g.renderAndWrite("/contact", "contact/index.html", []string{
		g.templatePath("contact.html"),
		g.templatePath("shared/carousel.html"),
	}, "contact", data)
}

func (g *generator) generateRSS() error {
	author := &feeds.Author{Name: "Tugberk Ugurlu"}
	feed := &feeds.Feed{
		Title:       "Tugberk Ugurlu @ the Heart of Software",
		Link:        &feeds.Link{Href: "https://www.tugberkugurlu.com"},
		Description: "Software Engineer and Tech Product enthusiast Tugberk Ugurlu's home on the interwebs! Here, you can find out about Tugberk's conference talks, books and blog posts on software development techniques and practices",
		Author:      author,
	}

	count := min(20, len(g.site.Posts))
	for _, post := range g.site.Posts[:count] {
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
		return fmt.Errorf("generating RSS XML: %w", err)
	}

	rssPath := filepath.Join(g.outputDir, "feeds", "rss")
	if err := os.MkdirAll(filepath.Dir(rssPath), 0755); err != nil {
		return fmt.Errorf("creating feeds directory: %w", err)
	}
	if err := os.WriteFile(rssPath, []byte(rss), 0644); err != nil {
		return fmt.Errorf("writing RSS file: %w", err)
	}

	// Add RSS to manifest
	hash := sha256.Sum256([]byte(rss))
	g.manifest.Pages = append(g.manifest.Pages, ManifestPage{
		Path:        "/feeds/rss",
		File:        "feeds/rss",
		ContentHash: hex.EncodeToString(hash[:])[:16],
	})

	return nil
}

func (g *generator) generate404() error {
	// The 404 page uses the layout template with a simple error message.
	// We create a minimal "not found" content template inline is not practical,
	// so we render it directly by writing the file.
	absFilePath := filepath.Join(g.outputDir, "404.html")

	f, err := os.Create(absFilePath)
	if err != nil {
		return fmt.Errorf("creating 404.html: %w", err)
	}
	defer f.Close()

	// Use the home template structure but render a 404 page.
	// We need a template that defines "twitter_card" and "content" blocks.
	// Since there is no dedicated 404 template, we create the HTML directly
	// using RenderPage with a custom template.
	err = g.renderNotFoundPage(f)
	if err != nil {
		return fmt.Errorf("rendering 404 page: %w", err)
	}

	return nil
}

func (g *generator) renderNotFoundPage(w io.Writer) error {
	// Write a 404 template file temporarily
	notFoundTemplate := `{{define "twitter_card"}}{{end}}
{{define "content"}}
<div class="ui container" style="padding: 40px 0; text-align: center;">
    <h1 class="ui header">
        <i class="frown outline icon"></i>
        Page Not Found
    </h1>
    <p>Sorry, the page you are looking for does not exist.</p>
    <a class="ui purple button" href="/">Go to Home Page</a>
</div>
{{end}}`

	tmpFile, err := os.CreateTemp("", "404-template-*.html")
	if err != nil {
		return fmt.Errorf("creating temp 404 template: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(notFoundTemplate); err != nil {
		return fmt.Errorf("writing temp 404 template: %w", err)
	}
	tmpFile.Close()

	layoutPath := g.templatePath("layout.html")
	return blog.RenderPage(w, g.site.LayoutConfig, []string{tmpFile.Name()}, layoutPath, "", g.site.TagsList, nil)
}

func (g *generator) writeManifest() error {
	// Add legacy image redirect pattern
	g.manifest.Patterns = append(g.manifest.Patterns, ManifestPattern{
		Pattern:        "/content/images/*",
		RedirectPrefix: "https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images",
		Status:         301,
	})

	data, err := json.MarshalIndent(g.manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	manifestPath := filepath.Join(g.outputDir, "_manifest.json")
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	return nil
}

// copyDir recursively copies a directory tree from src to dst.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
