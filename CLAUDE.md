# CLAUDE.md -- go-bloggy

This file provides guidance for AI agents (Claude Code and others) working on the `go-bloggy` repository.

---

## Project Overview

**go-bloggy** is a Go-based personal blogging platform for [tugberkugurlu.com](https://www.tugberkugurlu.com). It is a port of [tugberkugurlu/tugberk-web](https://github.com/tugberkugurlu/tugberk-web) to Go. The site serves blog posts, a speaking engagements page, an about page, a contact page, and an RSS feed.

- **Architecture:** Static site generator + HTTP development server. Markdown files on disk are parsed and rendered to HTML at build time. The site is deployed as static files to S3 + CloudFront.
- **Two commands:** `cmd/generate` (static site generator) and `cmd/web` (development HTTP server). Both use the shared `internal/blog` package.
- **UI:** Semantic UI 2.4.1 + jQuery 3.1.1, server-side rendered via Go `html/template`.

---

## Repository Layout

```
cmd/
  generate/                  # Static site generator CLI
    main.go                  # Loads site, renders pages, copies assets, emits manifest
    generate_test.go         # Tests: output structure, manifest, RSS, assets, 404
  web/                       # HTTP development server
    main.go                  # Router setup, HTTP handlers (uses internal/blog)
  migrate/
    main.go                  # One-off migration tool: SQL Server -> markdown files
  sherlock/
    main.go                  # AWS Lambda stub (not in active use)
internal/
  blog/                      # Shared blog engine (used by both cmd/generate and cmd/web)
    types.go                 # Post, Tag, Carousel, Site, page view models, Config
    loader.go                # LoadSite() -- parses posts, builds indexes, carousels
    templates.go             # RenderPage() -- shared template rendering
    carousels.go             # Carousel construction (Top Picks, tag-based, related posts)
    slugs.go                 # ToSlug() with C#/C++ overrides
    speaking.go              # Hardcoded speaking engagement data
    testdata/                # 4 fixture posts + config.yaml for unit tests
    *_test.go                # Unit tests for loader, carousels, slugs
  htmltextextractor/         # Extracts plain text from HTML (for reading-time calc)
  readingtime/               # Computes reading time (wordCount / 200 minutes)
test/
  comparison/                # Side-by-side comparison: cmd/generate vs cmd/web output
    comparison_test.go       # Starts both, compares all 356 pages + 22 redirects
  e2e/                       # E2E link integrity tests on generated output
    e2e_test.go              # Validates all internal links, assets, RSS links, og:url
web/
  posts/                     # 219 blog posts as markdown files, grouped by year
  template/                  # Go html/template files
    layout.html              # Master layout (nav, sidebar, footer)
    home.html, blog.html, post.html, tag.html, about.html, speaking.html, contact.html
    shared/                  # Reusable partials (carousel, post-item, speaking-card)
  static/                    # CSS, JS, images (Semantic UI, custom styles)
  static-root/               # Root-level static files (robots.txt, favicon.ico)
config.yaml                  # Runtime configuration (assets_url, assets_prefix)
docker-web.dockerfile        # Docker build for the web server
docker-compose.yml           # Compose config for local container runs
go.mod / go.sum              # Go modules (Go 1.22)
.github/workflows/
  ci.yml                     # PR checks: build, test, generate, comparison, e2e
  deploy.yml                 # Deploy to S3 + CloudFront on push to master
  stale.yml                  # Stale issue/PR management
plans/                       # Design specs and implementation plans
```

---

## Build and Run

### Static site generation

```bash
# From the repository root
go run ./cmd/generate \
  -posts ./web/posts \
  -templates ./web/template \
  -static ./web/static \
  -static-root ./web/static-root \
  -config . \
  -output ./output
```

This generates the complete static site in `./output/`:
- `index.html` (home), `archive/index.html`, `archive/{slug}/index.html` (219 posts)
- `tags/{tag}/index.html` (132 tag pages)
- `about/`, `speaking/`, `contact/` index pages
- `feeds/rss` (RSS XML, no extension)
- `404.html` (custom error page)
- `assets/` (copied from `web/static/`)
- `robots.txt`, `favicon.ico` (copied from `web/static-root/`)
- `_manifest.json` (all pages, redirects, patterns)

### Development HTTP server

```bash
# From cmd/web directory (uses relative paths)
cd cmd/web
go run .
```

The server defaults to `:8080`. Override with `SERVER_PORT`:

```bash
SERVER_PORT=3000 go run .
```

### Build and test

```bash
go build -v ./...           # Compile everything
go test -v ./...            # Run all tests (unit + comparison + e2e)
go vet ./...                # Static analysis
```

### Docker

```bash
docker build -t go-bloggy -f docker-web.dockerfile .
docker-compose up           # Exposes port 9000 -> container :8080
```

### Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `SERVER_PORT` | `8080` | HTTP listen port (cmd/web only) |
| `TUGBERKWEB_GoogleAnalytics__TrackingId` | _(empty)_ | Google Analytics tracking ID injected into templates |
| `BLOGGY_SQL_SERVER_CONN_STR` | -- | SQL Server connection string (migrate tool only) |

### Configuration File

`config.yaml` is loaded by Viper. `cmd/generate` accepts `-config` flag pointing to its directory. `cmd/web` searches `../../ci/` then `../../` relative to its working directory.

```yaml
assets_url: "/assets"    # Base URL for static assets; set to CDN URL in production
assets_prefix: ""        # Optional sub-path appended to assets_url
```

`assets_url` is **required** -- both commands exit fatally if missing.

---

## cmd/generate CLI Flags

| Flag | Default | Purpose |
|---|---|---|
| `-posts` | `./web/posts` | Path to blog posts directory |
| `-templates` | `./web/template` | Path to templates directory |
| `-static` | `./web/static` | Path to static assets directory |
| `-static-root` | `./web/static-root` | Path to static root files |
| `-config` | `.` | Path to directory containing config.yaml |
| `-output` | `./output` | Path to output directory (cleaned on each run) |

---

## Output Directory Structure

```
output/
  index.html                              # Home page
  archive/
    index.html                            # Blog archive
    redis-cluster-benefits-of.../
      index.html                          # Individual post (x219)
  tags/
    architecture/
      index.html                          # Tag page (one per tag, 132 total)
  about/
    index.html
  speaking/
    index.html
  contact/
    index.html
  feeds/
    rss                                   # RSS XML (no extension)
  assets/
    semantic-ui-css-2.4.1/                # Copied from web/static/
    stylesheets/
    images/
  404.html                                # Custom error page
  robots.txt                              # Copied from web/static-root/
  favicon.ico
  _manifest.json                          # Build manifest (not uploaded to S3)
```

---

## Manifest (_manifest.json)

The generator emits `_manifest.json` with three sections:

- **pages**: All generated pages with `path`, `file`, and `content_hash` (SHA-256 prefix)
- **redirects**: Non-canonical slug redirects (from alias to canonical, status 301)
- **patterns**: URL pattern redirects (e.g., `/content/images/*` to Azure Blob Storage)

The deploy pipeline reads this manifest to configure CloudFront Function redirects.

---

## Internal Packages

| Package | Purpose |
|---|---|
| `internal/blog` | Shared blog engine: `LoadSite()`, `RenderPage()`, `GetCarousels()`, `ToSlug()`, types, speaking data |
| `internal/readingtime` | `Calculate(text string) time.Duration` -- words/200 rounded up |
| `internal/htmltextextractor` | `Extract(html string) []string` -- pulls text from semantic HTML tags |

### internal/blog key exports

| Function/Type | Purpose |
|---|---|
| `LoadSite(cfg LoadSiteConfig) (*Site, error)` | Parse all posts, build indexes and carousels |
| `RenderPage(w, layoutConfig, templatePaths, layoutPath, section, tagsList, data)` | Render any page type through the layout template |
| `SectionFromPath(urlPath string) string` | Extract section name from URL path |
| `GeneratePostURL(post *Post) string` | Generate canonical full URL for a post |
| `GetCarousels(postsByID)` | Build site-wide carousels |
| `GetCarouselForTag(tagSlug, title, postsByTagSlug)` | Build tag-specific carousel |
| `GetRelatedPostsCarousel(post, postsByTagSlug)` | Build related posts carousel |
| `ToSlug(tag string) string` | Convert tag name to URL-safe slug |
| `Site` struct | All loaded site data: posts, indexes, carousels, config |

The `internal/blog` package has **no dependency on `net/http`** -- it is shared between both commands.

---

## Testing

### Test levels

| Level | Location | What it tests | Speed |
|---|---|---|---|
| Unit tests | `internal/blog/*_test.go`, `cmd/generate/generate_test.go` | Post parsing, indexing, carousels, slugs, generator output | Fast (~1s) |
| Side-by-side comparison | `test/comparison/comparison_test.go` | Every page matches between cmd/generate and cmd/web | Medium (~3s) |
| E2E link integrity | `test/e2e/e2e_test.go` | All internal links/assets resolve, RSS links valid | Medium (~7s) |

### Running specific test levels

```bash
go test ./internal/blog/           # Unit tests only
go test ./cmd/generate/            # Generator unit tests
go test ./test/comparison/         # Side-by-side comparison
go test ./test/e2e/                # E2E link integrity
go test ./...                      # All tests
```

The comparison and e2e tests skip in `-short` mode.

---

## Data Flow

### Site loading (internal/blog)

`LoadSite()` walks `web/posts/**/*.md`, then for each file:

1. Parse YAML front matter -> `PostMetadata`
2. Compile markdown body -> `template.HTML` (via `blackfriday/v2`)
3. Extract images from HTML (for carousel eligibility)
4. Calculate reading time (`internal/readingtime`)
5. Build in-memory indexes: `PostsByID`, `PostsBySlug`, `PostsByTagSlug`, `TagsBySlug`, `TagsList`
6. Build carousels (Top Picks from hardcoded post IDs in `carousels.go`)
7. Return `*Site` struct

### Static generation (cmd/generate)

1. Call `blog.LoadSite()` to load all site data
2. Render each page type using `blog.RenderPage()` with the appropriate templates
3. Write each page as `{path}/index.html`
4. Generate RSS feed using `gorilla/feeds`
5. Generate `404.html` using a temporary template through the layout
6. Copy `web/static/` -> `output/assets/` and `web/static-root/` -> `output/`
7. Write `_manifest.json` with all pages, redirects, and patterns

---

## CI/CD

### PR checks (.github/workflows/ci.yml)

1. `go build -v ./...`
2. Unit tests (`go test ./cmd/... ./internal/...`)
3. Run `cmd/generate` to produce static output
4. Verify expected output files exist
5. Side-by-side comparison test
6. E2E link integrity test

### Deploy (.github/workflows/deploy.yml)

Triggers on push to `master`:

1. All CI steps
2. OIDC authentication via `aws-actions/configure-aws-credentials@v4`
3. S3 upload with per-content-type metadata:
   - HTML: `Cache-Control: public, max-age=300`
   - RSS: `Cache-Control: public, max-age=900`
   - Assets: `Cache-Control: public, max-age=31536000, immutable`
4. CloudFront wildcard invalidation (`/*`)

**Required GitHub secrets:** `AWS_DEPLOY_ROLE_ARN`, `S3_BUCKET_NAME`, `CLOUDFRONT_DISTRIBUTION_ID`

---

## Key Data Models

```go
type Site struct {
    Config         Config
    LayoutConfig   LayoutConfig
    Posts          []*Post              // sorted newest-first
    PostsByID      map[string]*Post
    PostsBySlug    map[string]*Post     // includes non-canonical slugs
    PostsByTagSlug map[string][]*Post
    TagsBySlug     map[string]*Tag
    TagsList       TagCountPairList     // sorted by count desc, then alphabetical
    Carousels      []Carousel
    Speaking       []*SpeakingActivity
}

type Post struct {
    Body                    template.HTML
    Images                 []string
    Tags                   []TagCountPair
    Metadata               PostMetadata
    ReadingTimeDisplay     string
    ReadingTime            *time.Duration
    PublishedOn            time.Time
    PublishedOnDisplay     string
    PublishedOnDisplayBrief string
}

type PostMetadata struct {
    ID        string   `yaml:"id"`
    Title     string   `yaml:"title"`
    Abstract  string   `yaml:"abstract"`
    Format    string   `yaml:"format"`     // always "md"
    CreatedOn string   `yaml:"created_at"`
    Tags      []string `yaml:"tags"`
    Slugs     []string `yaml:"slugs"`      // first slug is canonical
}
```

Page view models: `Home`, `Blog`, `TagsPage`, `PostPage`, `SpeakingPage`, `AboutPage`, `ContactPage` -- passed to `RenderPage()` which wraps them in a `Layout` struct.

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
created_at: "2026-03-29 10:00:00 +0000 UTC"
tags:
  - Go
  - Software Engineering
slugs:
  - your-post-slug
---

Post body in markdown here.
```

3. Run the generator or restart the dev server -- posts are indexed at startup only.

---

## Conventions

- **Go version:** 1.22 (specified in `go.mod` and CI workflows).
- **No ORM, no database at runtime** -- data lives in markdown files.
- **No authentication** -- purely public read-only blog.
- **Shared code in `internal/blog`** -- both commands use the same parsing, rendering, and indexing logic. `internal/blog` must have NO dependency on `net/http`.
- **Hardcoded data** for carousels (`internal/blog/carousels.go`) and speaking engagements (`internal/blog/speaking.go`) -- change these files to update that content.
- **Tag slugs** follow `gosimple/slug` rules with explicit overrides for `C# -> c-sharp` and `C++ -> cpp`.
- **Deterministic builds** -- tag lists are sorted by count (desc) then alphabetically (asc) for reproducible output.
- **Clean URL convention** -- pages are written as `{path}/index.html` so CloudFront/S3 serves them at `/{path}`.
- Error handling: startup errors are `log.Fatal`; handler errors return HTTP 404/500.

---

## Dependencies (key ones)

| Package | Version | Use |
|---|---|---|
| `gorilla/mux` | v1.7.4 | HTTP router (cmd/web) |
| `gorilla/feeds` | v1.1.1 | RSS feed generation |
| `russross/blackfriday/v2` | v2.1.0 | Markdown -> HTML |
| `spf13/viper` | v1.7.1 | Config loading (YAML + env) |
| `gosimple/slug` | v1.9.0 | URL slug generation |
| `go-yaml/yaml` | v2.2.8 | YAML front-matter parsing |
| `golang.org/x/net/html` | latest | HTML tokenizing |
| `NYTimes/gziphandler` | v1.1.1 | GZIP middleware (cmd/web) |

---

## Migration Status

The site is being migrated from ECS Fargate (dynamic Go server) to S3 + CloudFront (static site). Current status:

- [x] Phase 1: Extract `internal/blog` package + tests (PR #22, merged)
- [x] Phase 2: Build `cmd/generate` static site generator (PR #23, merged)
- [x] Phase 3: Side-by-side comparison + E2E tests (PR #24, merged)
- [ ] Phase 4: AWS infrastructure (tugberkugurlu/tugberk-infrastructure#2)
- [ ] Phase 5: Deploy pipeline (PR #25, pending infrastructure)
- [ ] Phase 6: Go live
- [ ] Phase 7: Remove `cmd/web`, Docker files, ECS infrastructure

Design spec: `plans/2026-04-19-static-site-migration-design.md`
