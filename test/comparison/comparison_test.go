package comparison

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// manifestPage represents a page entry in _manifest.json.
type manifestPage struct {
	Path        string `json:"path"`
	File        string `json:"file"`
	ContentHash string `json:"content_hash"`
}

// manifestRedirect represents a redirect entry in _manifest.json.
type manifestRedirect struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Status int    `json:"status"`
}

// manifest represents the full _manifest.json structure.
type manifest struct {
	Pages     []manifestPage     `json:"pages"`
	Redirects []manifestRedirect `json:"redirects"`
}

// repoRoot returns the repository root relative to this test file.
func repoRoot(t *testing.T) string {
	t.Helper()
	// test/comparison/ -> repo root is ../..
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	return root
}

// getFreePort finds a free TCP port on localhost that fits in int16 range
// (required by cmd/web which parses SERVER_PORT with ParseInt base 10 size 16).
func getFreePort(t *testing.T) int {
	t.Helper()
	// Try ports in the safe range (9000-32000)
	for port := 9847; port < 32000; port++ {
		listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			listener.Close()
			return port
		}
	}
	t.Fatal("could not find a free port in the int16 range")
	return 0
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

// startWebServer starts cmd/web on the given port and returns a cleanup function.
func startWebServer(t *testing.T, root string, port int) func() {
	t.Helper()

	// cmd/web uses relative paths (../../web/posts, ../../web/template, etc.)
	// so it must be run from within cmd/web/ directory using "go run ."
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = filepath.Join(root, "cmd", "web")
	// Build env with SERVER_PORT override. Filter out any existing SERVER_PORT.
	env := os.Environ()
	filteredEnv := make([]string, 0, len(env)+1)
	for _, e := range env {
		if !strings.HasPrefix(e, "SERVER_PORT=") {
			filteredEnv = append(filteredEnv, e)
		}
	}
	filteredEnv = append(filteredEnv, fmt.Sprintf("SERVER_PORT=%d", port))
	cmd.Env = filteredEnv

	if err := cmd.Start(); err != nil {
		t.Fatalf("starting web server: %v", err)
	}

	// Wait for the server to be ready
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL)
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Verify server is up
	resp, err := http.Get(baseURL)
	if err != nil {
		cmd.Process.Kill()
		t.Fatalf("web server did not start within timeout: %v", err)
	}
	resp.Body.Close()

	return func() {
		cmd.Process.Kill()
		cmd.Wait()
	}
}

// normalizeHTML removes leading/trailing whitespace from each line and collapses
// multiple consecutive blank lines. This handles minor whitespace differences
// between the Go server response and the static file.
func normalizeHTML(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	prevBlank := false
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \t\r")
		if trimmed == "" {
			if prevBlank {
				continue
			}
			prevBlank = true
		} else {
			prevBlank = false
		}
		result = append(result, trimmed)
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

func TestSideBySideComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping side-by-side comparison test in short mode")
	}

	root := repoRoot(t)

	// Step 1: Generate static site
	t.Log("Generating static site...")
	outputDir := generateStaticSite(t, root)

	// Step 2: Start web server
	port := getFreePort(t)
	t.Logf("Starting web server on port %d...", port)
	cleanup := startWebServer(t, root, port)
	defer cleanup()

	baseURL := fmt.Sprintf("http://localhost:%d", port)

	// Step 3: Load manifest
	manifestData, err := os.ReadFile(filepath.Join(outputDir, "_manifest.json"))
	if err != nil {
		t.Fatalf("reading manifest: %v", err)
	}

	var m manifest
	if err := json.Unmarshal(manifestData, &m); err != nil {
		t.Fatalf("parsing manifest: %v", err)
	}

	t.Logf("Comparing %d pages and %d redirects...", len(m.Pages), len(m.Redirects))

	// Step 4: Compare each HTML page
	client := &http.Client{
		// Do not follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	mismatches := 0
	for _, page := range m.Pages {
		// Skip RSS feed -- it has a dynamic timestamp from gorilla/feeds
		if page.Path == "/feeds/rss" {
			continue
		}

		// Read static file
		staticFilePath := filepath.Join(outputDir, page.File)
		staticData, err := os.ReadFile(staticFilePath)
		if err != nil {
			t.Errorf("reading static file %s: %v", page.File, err)
			continue
		}

		// Fetch from web server
		resp, err := client.Get(baseURL + page.Path)
		if err != nil {
			t.Errorf("fetching %s: %v", page.Path, err)
			continue
		}
		serverData, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("reading response for %s: %v", page.Path, err)
			continue
		}

		if resp.StatusCode != 200 {
			t.Errorf("page %s: expected status 200, got %d", page.Path, resp.StatusCode)
			continue
		}

		// Compare normalized HTML
		staticNorm := normalizeHTML(string(staticData))
		serverNorm := normalizeHTML(string(serverData))

		if staticNorm != serverNorm {
			mismatches++
			// Show first difference for debugging
			staticLines := strings.Split(staticNorm, "\n")
			serverLines := strings.Split(serverNorm, "\n")
			for i := 0; i < len(staticLines) && i < len(serverLines); i++ {
				if staticLines[i] != serverLines[i] {
					t.Errorf("page %s: first difference at line %d:\n  static: %q\n  server: %q",
						page.Path, i+1, staticLines[i], serverLines[i])
					break
				}
			}
			if len(staticLines) != len(serverLines) {
				t.Errorf("page %s: different line count: static=%d server=%d",
					page.Path, len(staticLines), len(serverLines))
			}
			if mismatches >= 5 {
				t.Fatalf("too many mismatches (%d), stopping comparison", mismatches)
			}
		}
	}

	if mismatches == 0 {
		t.Logf("All %d pages match!", len(m.Pages)-1) // -1 for skipped RSS
	}

	// Step 5: Verify redirects
	for _, redirect := range m.Redirects {
		resp, err := client.Get(baseURL + redirect.From)
		if err != nil {
			t.Errorf("fetching redirect %s: %v", redirect.From, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != redirect.Status {
			t.Errorf("redirect %s: expected status %d, got %d",
				redirect.From, redirect.Status, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		// The Go server redirects to the full URL, while the manifest has the path.
		// Check that the location ends with the expected path.
		expectedSuffix := redirect.To
		if !strings.HasSuffix(location, expectedSuffix) {
			// The server might use the full URL form
			expectedFull := fmt.Sprintf("https://www.tugberkugurlu.com%s", redirect.To)
			if location != expectedFull {
				t.Errorf("redirect %s: expected Location ending with %q or %q, got %q",
					redirect.From, expectedSuffix, expectedFull, location)
			}
		}
	}
}
