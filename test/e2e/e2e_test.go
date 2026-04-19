package e2e

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// manifestPage represents a page entry in _manifest.json.
type manifestPage struct {
	Path        string `json:"path"`
	File        string `json:"file"`
	ContentHash string `json:"content_hash"`
}

// manifest represents the _manifest.json structure.
type manifest struct {
	Pages []manifestPage `json:"pages"`
}

// rss represents the top-level RSS XML structure for parsing.
type rss struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Link string `xml:"link"`
}

// repoRoot returns the repository root relative to this test file.
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	return root
}

// generateStaticSite runs cmd/generate and returns the output directory.
func generateStaticSite(t *testing.T, root string) string {
	t.Helper()
	outputDir := t.TempDir()

	cmd := exec.Command("go", "run", "./cmd/generate",
		"-posts", filepath.Join(root, "web", "posts"),
		"-templates", filepath.Join(root, "web", "template"),
		"-static", filepath.Join(root, "web", "static"),
		"-static-root", filepath.Join(root, "web", "static-root"),
		"-config", root,
		"-output", outputDir,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("running generator: %v\n%s", err, out)
	}
	return outputDir
}

// fileExistsInOutput checks if a given URL path resolves to a file in the
// output directory. It handles both directory-style paths (path/index.html)
// and direct files (path).
func fileExistsInOutput(outputDir, urlPath string) bool {
	// Remove leading /
	cleanPath := strings.TrimPrefix(urlPath, "/")
	if cleanPath == "" {
		cleanPath = "index.html"
	}

	// Try direct path
	direct := filepath.Join(outputDir, cleanPath)
	if _, err := os.Stat(direct); err == nil {
		return true
	}

	// Try as directory with index.html
	withIndex := filepath.Join(outputDir, cleanPath, "index.html")
	if _, err := os.Stat(withIndex); err == nil {
		return true
	}

	return false
}

// extractLinksFromHTML parses HTML and extracts all href and src attributes
// that point to internal paths (starting with /).
func extractLinksFromHTML(htmlContent string) (hrefs []string, srcs []string) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				val := attr.Val
				// Only consider paths starting with / but not // (protocol-relative URLs)
				isInternal := strings.HasPrefix(val, "/") && !strings.HasPrefix(val, "//")
				if attr.Key == "href" && isInternal {
					hrefs = append(hrefs, val)
				}
				if attr.Key == "src" && isInternal {
					srcs = append(srcs, val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return
}

func TestLinkIntegrity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E link integrity test in short mode")
	}

	root := repoRoot(t)

	t.Log("Generating static site...")
	outputDir := generateStaticSite(t, root)

	// Load manifest to know all page paths
	manifestData, err := os.ReadFile(filepath.Join(outputDir, "_manifest.json"))
	if err != nil {
		t.Fatalf("reading manifest: %v", err)
	}

	var m manifest
	if err := json.Unmarshal(manifestData, &m); err != nil {
		t.Fatalf("parsing manifest: %v", err)
	}

	// Build set of known paths from manifest
	knownPaths := make(map[string]bool)
	for _, page := range m.Pages {
		knownPaths[page.Path] = true
	}

	t.Logf("Checking links across %d pages...", len(m.Pages))

	brokenLinks := 0
	brokenAssets := 0

	for _, page := range m.Pages {
		// Skip non-HTML files (RSS feed)
		if !strings.HasSuffix(page.File, ".html") && page.Path != "/feeds/rss" {
			continue
		}
		if page.Path == "/feeds/rss" {
			continue // RSS links checked separately
		}

		filePath := filepath.Join(outputDir, page.File)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("reading %s: %v", page.File, err)
			continue
		}

		hrefs, srcs := extractLinksFromHTML(string(content))

		// Check internal links (hrefs)
		for _, href := range hrefs {
			// Strip fragment
			if idx := strings.Index(href, "#"); idx != -1 {
				href = href[:idx]
			}
			// Strip query string
			if idx := strings.Index(href, "?"); idx != -1 {
				href = href[:idx]
			}
			if href == "" || href == "/" {
				continue
			}

			// Skip /Feeds/Rss (legacy casing in template, will be lowercased by CloudFront Function)
			if strings.EqualFold(href, "/feeds/rss") {
				continue
			}

			// Skip /content/images/* (legacy redirect handled by CloudFront Function)
			if strings.HasPrefix(strings.ToLower(href), "/content/images/") {
				continue
			}

			if !fileExistsInOutput(outputDir, href) {
				brokenLinks++
				if brokenLinks <= 20 {
					t.Errorf("page %s: broken link %q - no file found in output", page.Path, href)
				}
			}
		}

		// Check internal sources (srcs)
		for _, src := range srcs {
			if !fileExistsInOutput(outputDir, src) {
				brokenAssets++
				if brokenAssets <= 20 {
					t.Errorf("page %s: broken asset src %q - no file found in output", page.Path, src)
				}
			}
		}
	}

	if brokenLinks > 20 || brokenAssets > 20 {
		t.Errorf("too many broken links (%d) or assets (%d) to show all", brokenLinks, brokenAssets)
	}

	if brokenLinks == 0 && brokenAssets == 0 {
		t.Log("All internal links and asset references are valid!")
	}
}

func TestRSSLinksPointToGeneratedPages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RSS link check in short mode")
	}

	root := repoRoot(t)
	outputDir := generateStaticSite(t, root)

	rssPath := filepath.Join(outputDir, "feeds", "rss")
	rssData, err := os.ReadFile(rssPath)
	if err != nil {
		t.Fatalf("reading RSS: %v", err)
	}

	var feed rss
	if err := xml.Unmarshal(rssData, &feed); err != nil {
		t.Fatalf("parsing RSS XML: %v", err)
	}

	if len(feed.Channel.Items) == 0 {
		t.Fatal("RSS feed has no items")
	}

	t.Logf("Checking %d RSS item links...", len(feed.Channel.Items))

	for _, item := range feed.Channel.Items {
		// RSS links are full URLs like https://www.tugberkugurlu.com/archive/slug
		link := item.Link
		if !strings.HasPrefix(link, "https://www.tugberkugurlu.com/") {
			t.Errorf("RSS link does not start with expected prefix: %s", link)
			continue
		}

		// Extract the path portion
		urlPath := strings.TrimPrefix(link, "https://www.tugberkugurlu.com")

		if !fileExistsInOutput(outputDir, urlPath) {
			t.Errorf("RSS link %s does not correspond to a generated page", link)
		}
	}
}

func TestAllTagLinksResolve(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping tag link resolution test in short mode")
	}

	root := repoRoot(t)
	outputDir := generateStaticSite(t, root)

	// Read the home page to find tag links in the sidebar
	homePath := filepath.Join(outputDir, "index.html")
	homeContent, err := os.ReadFile(homePath)
	if err != nil {
		t.Fatalf("reading home page: %v", err)
	}

	hrefs, _ := extractLinksFromHTML(string(homeContent))

	tagLinks := 0
	brokenTags := 0
	for _, href := range hrefs {
		if strings.HasPrefix(href, "/tags/") {
			tagLinks++
			if !fileExistsInOutput(outputDir, href) {
				brokenTags++
				t.Errorf("tag link %s does not resolve to a generated page", href)
			}
		}
	}

	if tagLinks == 0 {
		t.Error("no tag links found in home page sidebar")
	}

	t.Logf("Checked %d tag links, %d broken", tagLinks, brokenTags)
}

func TestManifestPageFilesExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping manifest file existence test in short mode")
	}

	root := repoRoot(t)
	outputDir := generateStaticSite(t, root)

	manifestData, err := os.ReadFile(filepath.Join(outputDir, "_manifest.json"))
	if err != nil {
		t.Fatalf("reading manifest: %v", err)
	}

	var m manifest
	if err := json.Unmarshal(manifestData, &m); err != nil {
		t.Fatalf("parsing manifest: %v", err)
	}

	for _, page := range m.Pages {
		filePath := filepath.Join(outputDir, page.File)
		info, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			t.Errorf("manifest page %s references non-existent file %s", page.Path, page.File)
			continue
		}
		if err != nil {
			t.Errorf("stat %s: %v", page.File, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("manifest page %s: file %s is empty", page.Path, page.File)
		}
		if page.ContentHash == "" {
			t.Errorf("manifest page %s has empty content hash", page.Path)
		}
	}

	t.Logf("All %d manifest pages have valid files", len(m.Pages))
}

func TestOgUrlMetaTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping og:url meta tag test in short mode")
	}

	root := repoRoot(t)
	outputDir := generateStaticSite(t, root)

	// Check a sample of post pages for correct og:url
	postDirs, err := os.ReadDir(filepath.Join(outputDir, "archive"))
	if err != nil {
		t.Fatalf("reading archive directory: %v", err)
	}

	checked := 0
	for _, dir := range postDirs {
		if !dir.IsDir() {
			continue
		}
		indexPath := filepath.Join(outputDir, "archive", dir.Name(), "index.html")
		content, err := os.ReadFile(indexPath)
		if err != nil {
			continue
		}

		htmlStr := string(content)
		expectedOgUrl := fmt.Sprintf("https://www.tugberkugurlu.com/archive/%s", dir.Name())

		if !strings.Contains(htmlStr, expectedOgUrl) {
			// The post template sets og:url, so this should match
			t.Errorf("post %s: og:url does not contain expected URL %s", dir.Name(), expectedOgUrl)
		}

		checked++
		if checked >= 10 {
			break // Spot-check 10 posts
		}
	}

	t.Logf("Verified og:url for %d post pages", checked)
}
