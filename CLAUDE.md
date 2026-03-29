# CLAUDE.md — go-bloggy

This file provides guidance for AI agents (Claude Code and others) working on the `go-bloggy` repository.

---

## Project Overview

**go-bloggy** is a Go-based personal blogging platform for [tugberkugurlu.com](https://www.tugberkugurlu.com). It is a port of [tugberkugurlu/tugberk-web](https://github.com/tugberkugurlu/tugberk-web) to Go. The site serves blog posts, a speaking engagements page, an about page, a contact page, and an RSS feed.

- **Architecture:** Static-content blog — markdown files on disk, parsed and indexed in-memory at startup. No database at runtime.
- **Router:** Gorilla Mux with host-based and path-based routing.
- **UI:** Semantic UI 2.4.1 + jQuery 3.1.1, server-side rendered via Go `html/template`.

---

## Repository Layout

```
cmd/
  web/                     # Main HTTP server (primary application)
    main.go                # Router setup, HTTP handlers, startup data loading
    carousels.go           # Carousel construction logic
    speaking_activity_data.go  # Hardcoded speaking engagement data
  migrate/
    main.go                # One-off migration tool: SQL Server → markdown files
  sherlock/
    main.go                # AWS Lambda stub (not in active use)
internal/
  htmltextextractor/       # Extracts plain text from HTML (for reading-time calc)
  readingtime/             # Computes reading time (wordCount / 200 minutes)
web/
  posts/                   # Blog posts as markdown files, grouped by year
  template/                # Go html/template files
    layout.html            # Master layout (nav, sidebar, footer)
    home.html, blog.html, post.html, tag.html, about.html, speaking.html, contact.html
    shared/                # Reusable partials (carousel, post-item, speaking-card)
  static/                  # CSS, JS, images (Semantic UI, custom styles)
  static-root/             # Root-level static files (robots.txt, favicon.ico)
config.yaml                # Runtime configuration (assets_url, assets_prefix)
docker-web.dockerfile      # Docker build for the web server
docker-compose.yml         # Compose config for local container runs
go.mod / go.sum            # Go modules
.github/workflows/go.yml   # CI: go build + go test
```

---

## Build and Run

### Local development

```bash
# From the repository root
go build -v ./...          # Compile everything
go test -v ./...           # Run tests (currently minimal)

# Run the web server (from cmd/web)
cd cmd/web
go run .                   # OR: go build -o web . && ./web
```

The server defaults to `:8080`. Override with `SERVER_PORT`:

```bash
SERVER_PORT=3000 go run .
```

Visit `http://localhost:8080`.

### Docker

```bash
docker build -t go-bloggy -f docker-web.dockerfile .
docker-compose up          # Exposes port 9000 → container :8080
```

### Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `SERVER_PORT` | `8080` | HTTP listen port |
| `TUGBERKWEB_GoogleAnalytics__TrackingId` | _(empty)_ | Google Analytics tracking ID injected into templates |
| `BLOGGY_SQL_SERVER_CONN_STR` | — | SQL Server connection string (migrate tool only) |

### Configuration File

`config.yaml` is loaded by Viper; the server searches for it in `../../ci/` then `../../` relative to the binary location.

```yaml
assets_url: "/assets"    # Base URL for static assets; set to CDN URL in production
assets_prefix: ""        # Optional sub-path appended to assets_url
```

`assets_url` is **required** — the server exits fatally if missing.

---

## HTTP Endpoints

All routes are registered in `cmd/web/main.go` (around line 284).

| Method | Path | Handler | Notes |
|---|---|---|---|
| GET | `/` | `homeHandler` | Home page: top carousel + 3 latest posts + 4 speaking activities |
| GET | `/archive` | `blogHomeHandler` | Full post archive (all posts, newest first) |
| GET | `/archive/{slug}` | `blogPostPageHandler` | Individual post page with related-posts carousel |
| GET | `/tags/{tag_slug}` | `tagsPageHandler` | Posts filtered by a tag slug |
| GET | `/about` | `aboutPageHandler` | About page with GeekTalks + Top carousel |
| GET | `/speaking` | `speakingPageHandler` | All 12 speaking engagements + GeekTalks carousel |
| GET | `/contact` | `contactPageHandler` | Contact page with Architecture carousel |
| GET | `/feeds/rss` | `rssHandler` | RSS feed — latest 20 posts, Cache-Control: max-age=900 |
| GET | `/content/images/*` | `legacyBlogImagesRedirector` | 301 redirect to Azure Blob Storage for legacy image URLs |
| GET | `/assets/*` | `http.FileServer` | Semantic UI, stylesheets, custom assets |
| GET | `/*` | `http.FileServer` | Static root files (robots.txt, favicon.ico) |

**Host routing:** requests to `tugberkugurlu.com` (without `www`) are redirected 301 to `www.tugberkugurlu.com`.

**Middleware applied to all HTML handlers:** GZIP compression (`NYTimes/gziphandler`), case-insensitive URL matching (`CaselessMatcher` — paths lowercased before routing).

---

## Data Flow at Startup

On startup, `cmd/web/main.go` walks `web/posts/**/*.md`, then for each file:

1. Parse YAML front matter → `PostMetadata`
2. Compile markdown body → `template.HTML` (via `blackfriday/v2`)
3. Extract images from HTML (for carousel eligibility)
4. Calculate reading time (`internal/readingtime`)
5. Build in-memory indexes:
   - `postsByID` — `map[string]*Post`
   - `postsBySlug` — `map[string]*Post` (multiple slugs per post supported)
   - `postsByTagSlug` — `map[string][]*Post`
   - `tagsBySlug` — `map[string]*Tag`
   - `tagsList` — `TagCountPairList` sorted by post count
6. Build carousels (Top Picks, GeekTalks, Architecture) from hardcoded post ID lists in `carousels.go`
7. Start HTTP server

**All content is served entirely from memory after startup. There is no database at runtime.**

---

## Key Data Models

```go
// Post (runtime representation of a markdown file)
type Post struct {
    Body                    template.HTML
    Images                 []string
    Tags                   []TagCountPair
    Metadata               PostMetadata
    ReadingTimeDisplay     string        // e.g., "5 minutes read"
    ReadingTime            *time.Duration
    PublishedOn            time.Time
    PublishedOnDisplay     string        // "2006-01-02 15:04:05"
    PublishedOnDisplayBrief string       // "2 January 2006"
}

// PostMetadata (YAML front matter inside each .md file)
type PostMetadata struct {
    ID        string   `yaml:"id"`
    Title     string   `yaml:"title"`
    Abstract  string   `yaml:"abstract"`
    Format    string   `yaml:"format"`     // always "md"
    CreatedOn string   `yaml:"created_at"`
    Tags      []string `yaml:"tags"`
    Slugs     []string `yaml:"slugs"`      // first slug is the canonical URL
}

// Carousel — at least 3 image-bearing posts required to be valid
type Carousel struct {
    Title string
    Posts []*Post
}
```

Page view models: `Home`, `Blog`, `TagsPage`, `PostPage`, `SpeakingPage`, `AboutPage`, `ContactPage` — passed to `ExecuteTemplate` which wraps them in a `Layout` struct.

---

## Adding a New Blog Post

1. Create a markdown file under `web/posts/{YEAR}/{YYYY-MM-DDTHH:MM:SS}_{slug}.md`.
2. Add YAML front matter at the top:

```yaml
---
id: "unique-post-id"
title: "Your Post Title"
abstract: "Short description for meta tags and archive listing."
format: "md"
created_at: "2026-03-29T10:00:00.0000000Z"
tags:
  - Go
  - Software Engineering
slugs:
  - your-post-slug
---

Post body in markdown here.
```

3. Restart the server — posts are indexed at startup only.

---

## Adding a New HTTP Handler

1. Define a handler function in `cmd/web/main.go`.
2. Create a page view model struct implementing the `Page` interface (`Title() string`, `Description() string`).
3. Register the route in the router block (~line 284).
4. Add a corresponding template file under `web/template/`.
5. Update `ExecuteTemplate` to handle the new page type if needed.

---

## Internal Packages

| Package | Purpose |
|---|---|
| `internal/readingtime` | `Calculate(text string) time.Duration` — words/200 rounded up |
| `internal/htmltextextractor` | `Extract(html string) []string` — pulls text from semantic HTML tags |

These are small and self-contained; tests belong alongside the package files.

---

## CI/CD

GitHub Actions workflow: `.github/workflows/go.yml`

```yaml
steps:
  - go build -v ./...
  - go test -v ./...
```

No deployment step is defined in the workflow (deployment is manual or via Docker).

---

## Conventions

- **Go version:** 1.14+ (Docker image uses 1.15.6; use latest stable for development).
- **No ORM, no database at runtime** — data lives in markdown files.
- **No authentication** — purely public read-only blog.
- **Global vars** for in-memory post indexes — intentional for a single-binary, single-process server.
- **Hardcoded data** for carousels (`carousels.go`) and speaking engagements (`speaking_activity_data.go`) — change these files to update that content.
- **Tag slugs** follow `gosimple/slug` rules with explicit overrides for `C# → c-sharp` and `C++ → cpp`.
- **Template function** `mod` (modulo) is registered for grid layout in templates.
- Error handling: startup errors are `log.Fatal`; handler errors return HTTP 404/500.

---

## Dependencies (key ones)

| Package | Version | Use |
|---|---|---|
| `gorilla/mux` | v1.7.4 | HTTP router |
| `gorilla/feeds` | v1.1.1 | RSS feed generation |
| `russross/blackfriday/v2` | v2.1.0 | Markdown → HTML |
| `spf13/viper` | v1.7.1 | Config loading (YAML + env) |
| `gosimple/slug` | v1.9.0 | URL slug generation |
| `go-yaml/yaml` | v2.2.8 | YAML front-matter parsing |
| `golang.org/x/net/html` | latest | HTML tokenizing |
| `NYTimes/gziphandler` | v1.1.1 | GZIP middleware |
| `pkg/errors` | v0.8.1 | Error wrapping |
