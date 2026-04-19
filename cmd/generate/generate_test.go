package main

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
)

// testFixtureDir returns the path to the shared testdata directory.
func testFixtureDir() string {
	return filepath.Join("..", "..", "internal", "blog", "testdata")
}

// setupTestTemplates creates minimal template files for testing.
// These templates define the required "layout", "twitter_card", and "content"
// blocks to match the real template structure.
func setupTestTemplates(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// layout.html
	layoutContent := `{{define "layout"}}<!doctype html><html><head><title>{{ .Title }}</title></head><body>{{template "twitter_card" .Data}}{{template "content" .Data}}</body></html>{{end}}`
	writeTestFile(t, filepath.Join(dir, "layout.html"), layoutContent)

	// home.html
	homeContent := `{{define "twitter_card"}}{{end}}{{define "content"}}<div id="home">{{ range $post := .Data.Posts }}<div>{{ $post.Metadata.Title }}</div>{{ end }}</div>{{end}}`
	writeTestFile(t, filepath.Join(dir, "home.html"), homeContent)

	// blog.html (archive)
	blogContent := `{{define "twitter_card"}}{{end}}{{define "content"}}<div id="archive">{{ range $post := .Data.Posts }}<div>{{ $post.Metadata.Title }}</div>{{ end }}</div>{{end}}`
	writeTestFile(t, filepath.Join(dir, "blog.html"), blogContent)

	// post.html
	postContent := `{{define "twitter_card"}}{{end}}{{define "content"}}<div id="post"><h1>{{ .Data.Post.Metadata.Title }}</h1><div>{{ .Data.Post.Body }}</div></div>{{end}}`
	writeTestFile(t, filepath.Join(dir, "post.html"), postContent)

	// tag.html
	tagContent := `{{define "twitter_card"}}{{end}}{{define "content"}}<div id="tag"><h1>{{ .Data.Tag.Name }}</h1>{{ range $post := .Data.Posts }}<div>{{ $post.Metadata.Title }}</div>{{ end }}</div>{{end}}`
	writeTestFile(t, filepath.Join(dir, "tag.html"), tagContent)

	// about.html
	aboutContent := `{{define "twitter_card"}}{{end}}{{define "content"}}<div id="about">About page</div>{{end}}`
	writeTestFile(t, filepath.Join(dir, "about.html"), aboutContent)

	// speaking.html
	speakingContent := `{{define "twitter_card"}}{{end}}{{define "content"}}<div id="speaking">Speaking page</div>{{end}}`
	writeTestFile(t, filepath.Join(dir, "speaking.html"), speakingContent)

	// contact.html
	contactContent := `{{define "twitter_card"}}{{end}}{{define "content"}}<div id="contact">Contact page</div>{{end}}`
	writeTestFile(t, filepath.Join(dir, "contact.html"), contactContent)

	// shared/carousel.html (empty define, referenced but not needed in minimal templates)
	os.MkdirAll(filepath.Join(dir, "shared"), 0755)
	writeTestFile(t, filepath.Join(dir, "shared", "carousel.html"), `{{define "carousel"}}{{end}}`)
	writeTestFile(t, filepath.Join(dir, "shared", "post-item.html"), `{{define "post-item"}}{{end}}`)
	writeTestFile(t, filepath.Join(dir, "shared", "speaking-activity-card.html"), `{{define "speaking_activity_card"}}{{end}}`)

	return dir
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// setupTestStaticDirs creates minimal static asset and static-root directories.
func setupTestStaticDirs(t *testing.T) (string, string) {
	t.Helper()
	staticDir := t.TempDir()
	staticRootDir := t.TempDir()

	// Create a sample static asset
	os.MkdirAll(filepath.Join(staticDir, "stylesheets"), 0755)
	writeTestFile(t, filepath.Join(staticDir, "stylesheets", "main.css"), "body { margin: 0; }")

	// Create static root files
	writeTestFile(t, filepath.Join(staticRootDir, "robots.txt"), "User-agent: *\nAllow: /")
	writeTestFile(t, filepath.Join(staticRootDir, "favicon.ico"), "fake-ico-data")

	return staticDir, staticRootDir
}

func TestGenerateOutputStructure(t *testing.T) {
	templatesDir := setupTestTemplates(t)
	staticDir, staticRootDir := setupTestStaticDirs(t)
	outputDir := t.TempDir()

	err := run(testFixtureDir(), templatesDir, staticDir, staticRootDir, testFixtureDir(), outputDir)
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	// Verify expected files exist
	expectedFiles := []string{
		"index.html",
		"archive/index.html",
		"archive/first-test-post/index.html",
		"archive/second-test-post/index.html",
		"archive/third-test-post/index.html",
		"archive/fourth-test-post-no-images/index.html",
		"tags/go/index.html",
		"tags/software-engineering/index.html",
		"tags/architecture/index.html",
		"tags/c-sharp/index.html",
		"tags/cpp/index.html",
		"about/index.html",
		"speaking/index.html",
		"contact/index.html",
		"feeds/rss",
		"404.html",
		"_manifest.json",
		"assets/stylesheets/main.css",
		"robots.txt",
		"favicon.ico",
	}

	for _, expected := range expectedFiles {
		path := filepath.Join(outputDir, expected)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", expected)
		}
	}
}

func TestManifestIsComplete(t *testing.T) {
	templatesDir := setupTestTemplates(t)
	staticDir, staticRootDir := setupTestStaticDirs(t)
	outputDir := t.TempDir()

	err := run(testFixtureDir(), templatesDir, staticDir, staticRootDir, testFixtureDir(), outputDir)
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	manifestPath := filepath.Join(outputDir, "_manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("reading manifest: %v", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parsing manifest: %v", err)
	}

	// 4 posts + archive + home + about + speaking + contact + 5 tags + RSS = 15 pages
	// Tags: go, software-engineering, architecture, c-sharp, cpp
	expectedPageCount := 15
	if len(manifest.Pages) != expectedPageCount {
		t.Errorf("expected %d pages in manifest, got %d", expectedPageCount, len(manifest.Pages))
		for _, p := range manifest.Pages {
			t.Logf("  page: %s -> %s", p.Path, p.File)
		}
	}

	// post1 has "first-test-post-alias" as a second slug
	expectedRedirectCount := 1
	if len(manifest.Redirects) != expectedRedirectCount {
		t.Errorf("expected %d redirects in manifest, got %d", expectedRedirectCount, len(manifest.Redirects))
		for _, r := range manifest.Redirects {
			t.Logf("  redirect: %s -> %s", r.From, r.To)
		}
	}

	// Verify the redirect
	if len(manifest.Redirects) > 0 {
		r := manifest.Redirects[0]
		if r.From != "/archive/first-test-post-alias" {
			t.Errorf("expected redirect from /archive/first-test-post-alias, got %s", r.From)
		}
		if r.To != "/archive/first-test-post" {
			t.Errorf("expected redirect to /archive/first-test-post, got %s", r.To)
		}
		if r.Status != 301 {
			t.Errorf("expected redirect status 301, got %d", r.Status)
		}
	}

	// Verify pattern
	if len(manifest.Patterns) != 1 {
		t.Errorf("expected 1 pattern in manifest, got %d", len(manifest.Patterns))
	}
	if len(manifest.Patterns) > 0 {
		p := manifest.Patterns[0]
		if p.Pattern != "/content/images/*" {
			t.Errorf("expected pattern /content/images/*, got %s", p.Pattern)
		}
	}

	// Verify all pages have non-empty content hashes
	for _, page := range manifest.Pages {
		if page.ContentHash == "" {
			t.Errorf("page %s has empty content hash", page.Path)
		}
		if page.File == "" {
			t.Errorf("page %s has empty file path", page.Path)
		}
	}
}

func TestRSSIsValidXML(t *testing.T) {
	templatesDir := setupTestTemplates(t)
	staticDir, staticRootDir := setupTestStaticDirs(t)
	outputDir := t.TempDir()

	err := run(testFixtureDir(), templatesDir, staticDir, staticRootDir, testFixtureDir(), outputDir)
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	rssPath := filepath.Join(outputDir, "feeds", "rss")
	data, err := os.ReadFile(rssPath)
	if err != nil {
		t.Fatalf("reading RSS: %v", err)
	}

	// Verify it is valid XML
	var xmlDoc interface{}
	if err := xml.Unmarshal(data, &xmlDoc); err != nil {
		t.Fatalf("RSS is not valid XML: %v", err)
	}

	// Verify it contains expected RSS elements
	rssContent := string(data)
	if !contains(rssContent, "<rss") {
		t.Error("RSS missing <rss> element")
	}
	if !contains(rssContent, "<channel>") {
		t.Error("RSS missing <channel> element")
	}
	if !contains(rssContent, "<item>") {
		t.Error("RSS missing <item> elements")
	}
	if !contains(rssContent, "Tugberk Ugurlu") {
		t.Error("RSS missing author name")
	}
}

func TestStaticAssetsCopied(t *testing.T) {
	templatesDir := setupTestTemplates(t)
	staticDir, staticRootDir := setupTestStaticDirs(t)
	outputDir := t.TempDir()

	err := run(testFixtureDir(), templatesDir, staticDir, staticRootDir, testFixtureDir(), outputDir)
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	// Check static assets
	cssPath := filepath.Join(outputDir, "assets", "stylesheets", "main.css")
	cssData, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("reading copied CSS: %v", err)
	}
	if string(cssData) != "body { margin: 0; }" {
		t.Error("copied CSS content does not match source")
	}

	// Check static root files
	robotsPath := filepath.Join(outputDir, "robots.txt")
	robotsData, err := os.ReadFile(robotsPath)
	if err != nil {
		t.Fatalf("reading copied robots.txt: %v", err)
	}
	if string(robotsData) != "User-agent: *\nAllow: /" {
		t.Error("copied robots.txt content does not match source")
	}
}

func Test404PageGenerated(t *testing.T) {
	templatesDir := setupTestTemplates(t)
	staticDir, staticRootDir := setupTestStaticDirs(t)
	outputDir := t.TempDir()

	err := run(testFixtureDir(), templatesDir, staticDir, staticRootDir, testFixtureDir(), outputDir)
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	path404 := filepath.Join(outputDir, "404.html")
	data, err := os.ReadFile(path404)
	if err != nil {
		t.Fatalf("reading 404.html: %v", err)
	}

	content := string(data)
	if !contains(content, "Page Not Found") {
		t.Error("404.html missing 'Page Not Found' text")
	}
	if !contains(content, "<html") {
		t.Error("404.html missing HTML structure")
	}
}

func TestGeneratedHTMLNotEmpty(t *testing.T) {
	templatesDir := setupTestTemplates(t)
	staticDir, staticRootDir := setupTestStaticDirs(t)
	outputDir := t.TempDir()

	err := run(testFixtureDir(), templatesDir, staticDir, staticRootDir, testFixtureDir(), outputDir)
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	filesToCheck := []string{
		"index.html",
		"archive/index.html",
		"archive/first-test-post/index.html",
		"about/index.html",
		"speaking/index.html",
		"contact/index.html",
	}

	for _, file := range filesToCheck {
		path := filepath.Join(outputDir, file)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("%s: stat error: %v", file, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("%s: file is empty", file)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
