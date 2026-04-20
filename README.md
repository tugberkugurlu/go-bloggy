# go-bloggy

![CI](https://github.com/tugberkugurlu/go-bloggy/actions/workflows/ci.yml/badge.svg?branch=master)

Static site generator and blog for [tugberkugurlu.com](https://www.tugberkugurlu.com), implemented in [Go](https://go.dev/). Markdown files on disk are compiled into a complete static site that is deployed to AWS S3 + CloudFront.

This is a port of [tugberkugurlu/tugberk-web](https://github.com/tugberkugurlu/tugberk-web).

## Quick start

```bash
# Build everything
go build -v ./...

# Run all tests (unit + integration)
go test -v ./...

# Generate the static site into ./output
go run ./cmd/generate \
  -posts ./web/posts \
  -templates ./web/template \
  -static ./web/static \
  -static-root ./web/static-root \
  -config . \
  -output ./output
```

The generated site is written to the `output/` directory with an `_manifest.json` listing every page, redirect, and URL pattern.

## Running the dev server

The legacy `cmd/web` server can still be used for local development. It loads posts into memory at startup and serves them on the fly:

```bash
cd cmd/web
go run .
# Visit http://localhost:8080
```

Override the port with `SERVER_PORT`:

```bash
SERVER_PORT=3000 go run .
```

## Docker (local dev)

A `docker-compose.dev.yml` is provided for container-based local development. It mounts the content directories as read-only volumes so content changes take effect on restart without rebuilding the image:

```bash
docker compose -f docker-compose.dev.yml up --build
# Visit http://localhost:9847
```

To build and run the production image directly:

```bash
docker build -t go-bloggy -f docker-web.dockerfile .
docker run -it --rm -p 9000:8080 go-bloggy
```

## Running tests

```bash
# Unit tests only
go test -v ./cmd/... ./internal/...

# All tests including comparison and E2E link integrity
go test -v ./...
```

## Project structure

```
cmd/
  generate/          Static site generator (primary build tool)
  web/               Legacy HTTP server for local development
  migrate/           One-off SQL Server to markdown migration tool
  sherlock/          AWS Lambda stub (not in active use)
internal/
  blog/              Core blog logic: post loading, templates, carousels, types
  htmltextextractor/ Extracts plain text from HTML (for reading-time calc)
  readingtime/       Computes reading time (wordCount / 200 minutes)
web/
  posts/             Blog posts as markdown files, grouped by year
  template/          Go html/template files (layout, pages, partials)
  static/            CSS, JS, images (Semantic UI, custom styles)
  static-root/       Root-level static files (robots.txt, favicon.ico)
test/
  comparison/        Side-by-side comparison tests (server vs generator)
  e2e/               E2E link integrity tests
config.yaml          Runtime configuration (assets_url, assets_prefix)
plans/               Design specs and migration plans
docs/                Operational documentation
```

## CI/CD pipeline

Two GitHub Actions workflows drive the build and deployment:

- **CI** (`.github/workflows/ci.yml`) -- runs on pull requests to `master`:
  1. Build all Go packages
  2. Run unit tests
  3. Generate the static site and verify expected output files
  4. Run side-by-side comparison and E2E link integrity tests

- **Deploy** (`.github/workflows/deploy.yml`) -- runs on push to `master`:
  1. Same build, test, and generate steps as CI
  2. Authenticate to AWS via OIDC (GitHub Actions role)
  3. Sync HTML pages, RSS feed, static assets, and root files to S3 with appropriate cache headers
  4. Invalidate the CloudFront distribution cache

The infrastructure (S3 bucket, CloudFront distribution, IAM roles, DNS) is managed in [tugberkugurlu/tugberk-infrastructure](https://github.com/tugberkugurlu/tugberk-infrastructure).

## Design spec

See [plans/2026-04-19-static-site-migration-design.md](plans/2026-04-19-static-site-migration-design.md) for the original design specification for the static site migration.

## Acknowledgements

Huge thanks to these open source projects:

- [Semantic-Org/Semantic-UI](https://github.com/Semantic-Org/Semantic-UI) -- UI framework
- [gorilla/mux](https://github.com/gorilla/mux) -- HTTP router
- [gorilla/feeds](https://github.com/gorilla/feeds) -- RSS feed generation
- [russross/blackfriday](https://github.com/russross/blackfriday) -- Markdown to HTML
- [spf13/viper](https://github.com/spf13/viper) -- Configuration loading
- [gosimple/slug](https://github.com/gosimple/slug) -- URL slug generation
- [go-yaml/yaml](https://github.com/go-yaml/yaml) -- YAML front-matter parsing
