# go-bloggy Static Site Migration — Design Spec

**Date:** 2026-04-19
**Status:** Draft
**Repo:** github.com/tugberkugurlu/go-bloggy
**Infra Repo:** github.com/tugberkugurlu/tugberk-infrastructure

---

## 1. Goal

Convert the go-bloggy personal blog from a Go HTTP server hosted on AWS ECS Fargate to a fully static site hosted on S3 + CloudFront. The blog is currently down (ECS deployment stopped). The aim is to bring it back online as a static site with:

- Full URL parity with the previous live site
- Custom domain (www.tugberkugurlu.com) with HTTPS
- Automated deployment on merge to master
- Side-by-side comparison testing to prove content equivalence
- Significantly reduced hosting cost (~$1-2/month vs ECS+ALB+VPC+NAT)

Blog content markdown files must not change.

---

## 2. Current State

### 2.1 go-bloggy Repository

- **Architecture:** Go HTTP server (`cmd/web/main.go`, 664 lines) that parses 219 markdown posts at startup, builds in-memory indexes, and serves everything from memory via Gorilla Mux.
- **Content:** 219 posts (2010–2021) in `web/posts/{year}/`, YAML front matter + markdown body. 20 posts have multiple slugs (canonical + aliases for backward compatibility).
- **Templates:** Go `html/template` files in `web/template/` — layout, 7 page templates, 3 shared partials.
- **Static assets:** ~7.6 MB in `web/static/` (Semantic UI 2.4.1, CSS, images). Root files (robots.txt, favicon.ico) in `web/static-root/`.
- **Tests:** Zero `_test.go` files exist. No test coverage.
- **CI:** GitHub Actions runs `go build` + `go test` on PRs. No deployment step.
- **Code organization:** Monolithic — all types, parsing logic, HTTP handlers, template rendering, and global state live in `cmd/web/main.go`.

### 2.2 Current URL Structure (must be preserved)

| Path | Behavior |
|---|---|
| `/` | Home page — top carousel, 3 latest posts, 4 speaking activities |
| `/archive` | Blog archive — all 219 posts, newest first |
| `/archive/{slug}` | Individual post — 219 canonical slugs |
| `/archive/{alias}` | 301 redirect to canonical slug (20 posts have aliases) |
| `/tags/{tag_slug}` | Posts filtered by tag |
| `/about` | About page |
| `/speaking` | Speaking engagements page |
| `/contact` | Contact page (embedded Google Form) |
| `/feeds/rss` | RSS feed (latest 20 posts, `Cache-Control: max-age=900`) |
| `/content/images/*` | 301 redirect to Azure Blob Storage |
| `/assets/*` | Static files (Semantic UI, CSS, images) |
| `/*` (root) | Static root files (robots.txt, favicon.ico) |

**Middleware behaviors:**
- `CaselessMatcher` — lowercases all incoming URL paths before routing
- `tugberkugurlu.com` (naked) → 301 to `www.tugberkugurlu.com`
- HTTP → HTTPS redirect (handled by ALB)
- GZIP compression (handled by `NYTimes/gziphandler`)

### 2.3 tugberk-infrastructure Repository

Terraform-managed AWS infrastructure in us-east-2:

- **ECS Fargate** cluster with 2 tasks (256 CPU / 512 MB each)
- **ALB** with HTTPS listener (ACM wildcard cert for `*.tugberkugurlu.com`)
- **VPC** with 3 public + 3 private subnets, single NAT Gateway
- **Route53** zone for `tugberkugurlu.com` (A record → ALB, www CNAME → ALB)
- **S3 bucket** `tugberkugurlu-blog` — stores post images (referenced in post bodies, must be preserved)
- **GitHub Actions IAM** — OIDC federation already configured

The ECS cluster is stopped but all resources still exist in Terraform state.

### 2.4 External Image References

Blog posts reference images on two external hosts:
- `tugberkugurlu-blog.s3.us-east-2.amazonaws.com/post-images/...` — S3 bucket (must be preserved)
- `tugberkugurlu.blob.core.windows.net/bloggyimages/...` — Azure Blob Storage

These are direct URLs in post markdown/HTML. They do not need to be moved or changed.

---

## 3. Approach

Build a new `cmd/generate` Go command inside the go-bloggy repo that reuses the existing Go parsing/template logic to emit a complete static site to a folder. This approach:

- Reuses all existing code (YAML parsing, markdown rendering, tag indexing, carousel logic, templates)
- Guarantees blog content markdown files don't change
- Enables side-by-side comparison testing (same data + templates → both commands should produce identical HTML)
- Keeps everything in one language/repo

---

## 4. Architecture

### 4.1 Static Site Generation

The `cmd/generate` command is a CLI that:

1. **Loads** the site data via `blog.LoadSite(postsDir)` — identical to what `cmd/web` does at startup
2. **Renders** each page to HTML bytes using the same Go templates
3. **Writes** each page as `{path}/index.html` (e.g., `archive/my-post/index.html`)
4. **Copies** static assets from `web/static/` → `output/assets/` and `web/static-root/` → `output/`
5. **Emits** `_manifest.json` with all pages and redirects

CLI interface:
```
go run ./cmd/generate \
  -posts ./web/posts \
  -templates ./web/template \
  -static ./web/static \
  -static-root ./web/static-root \
  -config ./config.yaml \
  -output ./output
```

### 4.2 Output Directory Structure

```
output/
  index.html                              # Home page
  archive/
    index.html                            # Blog archive
    redis-cluster-benefits-of.../
      index.html                          # Individual post (x219)
  tags/
    architecture/
      index.html                          # Tag page (one per tag)
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
  _manifest.json                          # Build artifact (not uploaded to S3)
```

Using `{path}/index.html` convention means CloudFront + S3 natively serve clean URLs (e.g., `/archive/my-post` resolves to `archive/my-post/index.html`).

The `feeds/rss` file has no extension and requires explicit `Content-Type: application/rss+xml` set via S3 object metadata at upload time.

### 4.3 Hosting: S3 + CloudFront

**New S3 bucket** (e.g., `tugberkugurlu-com-static`) for the static site. The existing `tugberkugurlu-blog` bucket (post images) is preserved unchanged.

- **CloudFront distribution** with OAC (Origin Access Control) pointing to the new S3 bucket
- **ACM certificate** — reuse existing cert which covers both `tugberkugurlu.com` (primary) and `*.tugberkugurlu.com` (SAN). Both names are required because CloudFront terminates TLS before the CloudFront Function can redirect naked → www.
- **CloudFront settings** — Viewer Protocol Policy: redirect-to-https (replaces ALB HTTP→HTTPS redirect). Compress objects automatically: enabled (replaces the Go server's GZIP middleware).
- **Custom error response** — configure CloudFront to return a custom `404.html` page for 403/404 errors from S3, instead of the default S3 XML error.
- **Route53** — point `tugberkugurlu.com` (A alias) and `www.tugberkugurlu.com` (A alias) to CloudFront

### 4.4 CloudFront Function (viewer-request)

A CloudFront Function handles all dynamic behaviors previously in the Go server:

```javascript
function handler(event) {
  var request = event.request;
  var uri = request.uri.toLowerCase();  // Case-insensitive matching
  var host = request.headers.host.value;

  // Naked domain → www redirect (preserving query string)
  if (host === 'tugberkugurlu.com') {
    var qs = Object.keys(request.querystring).length > 0
      ? '?' + Object.keys(request.querystring).map(function(k) {
          return k + '=' + request.querystring[k].value;
        }).join('&')
      : '';
    return {
      statusCode: 301,
      headers: { location: { value: 'https://www.tugberkugurlu.com' + uri + qs } }
    };
  }

  // Legacy image redirects
  if (uri.startsWith('/content/images/')) {
    var imagePath = uri.substring('/content/images'.length);
    return {
      statusCode: 301,
      headers: { location: { value: 'https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images' + imagePath } }
    };
  }

  // Non-canonical slug redirects — baked in at deploy time.
  // The deploy script reads _manifest.json and generates this map.
  var redirects = {
    "/archive/singalr-with-redis-running-on-a-windows-azure-virtual-machine":
      "/archive/signalr-with-redis-running-on-a-windows-azure-virtual-machine",
    // ... all 20 non-canonical slug entries from manifest
  };
  if (redirects[uri]) {
    return {
      statusCode: 301,
      headers: { location: { value: 'https://www.tugberkugurlu.com' + redirects[uri] } }
    };
  }

  // Apply lowercased URI
  request.uri = uri;

  // Add index.html for directory-style paths
  if (uri.endsWith('/')) {
    request.uri += 'index.html';
  } else if (!uri.includes('.')) {
    request.uri += '/index.html';
  }

  return request;
}
```

### 4.5 Manifest

The generator emits `_manifest.json`:

```json
{
  "pages": [
    {"path": "/", "file": "index.html", "content_hash": "abc123"},
    {"path": "/archive", "file": "archive/index.html", "content_hash": "def456"},
    {"path": "/archive/redis-cluster-...", "file": "archive/redis-cluster-.../index.html", "content_hash": "ghi789"}
  ],
  "redirects": [
    {"from": "/archive/singalr-with-redis-running-on-a-windows-azure-virtual-machine",
     "to": "/archive/signalr-with-redis-running-on-a-windows-azure-virtual-machine", "status": 301}
  ],
  "patterns": [
    {"pattern": "/content/images/*",
     "redirect_prefix": "https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images", "status": 301}
  ]
}
```

---

## 5. Codebase Refactoring

### 5.1 Extract `internal/blog/` Package

The monolithic `cmd/web/main.go` is split into a shared `internal/blog/` package:

| New File | Contents |
|---|---|
| `internal/blog/types.go` | `Post`, `PostMetadata`, `Tag`, `TagCountPair`, `Carousel`, `SpeakingActivity`, `Config`, `Site`, page view models (`Home`, `Blog`, `TagsPage`, `PostPage`, `SpeakingPage`, `AboutPage`, `ContactPage`) |
| `internal/blog/loader.go` | `LoadSite(cfg LoadSiteConfig) (*Site, error)` — accepts paths (postsDir, configPath) and options, walks markdown files, parses YAML + markdown, extracts images, calculates reading time, builds all indexes, builds carousels. Returns a `Site` struct. Config parsing (Viper) moves here. |
| `internal/blog/slugs.go` | `ToSlug(tag string) string` with predefined overrides (`c# → c-sharp`, `c++ → cpp`) |
| `internal/blog/carousels.go` | `GetCarousels()`, `GetCarouselForTag()`, `GetRelatedPostsCarousel()`, top picks post IDs |
| `internal/blog/speaking.go` | `SpeakingActivity` struct + `SpeakingActivities` data |
| `internal/blog/templates.go` | `RenderPage()` — shared template execution wrapping `ExecuteTemplate` logic |

### 5.2 The `Site` Struct

```go
type Site struct {
    Config         Config
    Posts          []*Post              // sorted newest-first
    PostsByID      map[string]*Post
    PostsBySlug    map[string]*Post     // includes non-canonical slugs
    PostsByTagSlug map[string][]*Post
    TagsBySlug     map[string]*Tag
    TagsList       TagCountPairList
    Carousels      []Carousel
    Speaking       []*SpeakingActivity
}
```

Both `cmd/web` and `cmd/generate` call `blog.LoadSite()` and get the same `Site` instance.

### 5.3 Slimmed `cmd/web/main.go`

After extraction, `cmd/web/main.go` contains only:
- HTTP server setup and route registration
- HTTP handler functions (thin — call into `blog` package for data and rendering)
- Middleware (GZIP, CaselessMatcher)

---

## 6. Caching Strategy

### 6.1 S3 Objects

Every S3 object gets an automatic ETag (MD5 hash of content). `aws s3 sync` only uploads files whose content changed.

### 6.2 CloudFront Cache-Control Headers

| Content Type | `Cache-Control` | Rationale |
|---|---|---|
| HTML pages (`.html`) | `public, max-age=300` | 5 min, then revalidate via ETag |
| Static assets (`/assets/*`) | `public, max-age=31536000, immutable` | 1 year — Semantic UI is vendored, CSS rarely changes |
| RSS feed (`/feeds/rss`) | `public, max-age=900` | 15 min, matches previous behavior |

### 6.3 Cache Invalidation

On deploy, issue a wildcard invalidation: `aws cloudfront create-invalidation --paths "/*"`. This counts as 1 path (free tier: 1,000/month). CloudFront re-fetches from S3 on next request; unchanged content gets a fast 304 via ETag.

### 6.4 Browser Revalidation

After `max-age` expires, browsers send `If-None-Match: "etag"`. CloudFront checks S3 and returns 304 Not Modified if content hasn't changed.

---

## 7. Deployment

### 7.1 Security: GitHub Actions OIDC

The go-bloggy repo is public. Credentials are handled via OIDC federation — no stored secrets:

- IAM role trust policy scoped to `repo:tugberkugurlu/go-bloggy:ref:refs/heads/master`
- Deploy workflow triggers only on `push` to `master` (not on PRs)
- PR workflows (`pull_request` event) run tests only — no AWS access needed, fork PRs cannot access OIDC
- IAM role permissions: `s3:PutObject`, `s3:DeleteObject`, `s3:ListBucket` on the static site bucket + `cloudfront:CreateInvalidation`

### 7.2 CI Workflows

**`.github/workflows/ci.yml`** (on `pull_request`):
1. `go build ./...`
2. `go test ./...`
3. `go run ./cmd/generate -output ./output`
4. Side-by-side comparison test
5. E2E link integrity test

**`.github/workflows/deploy.yml`** (on `push` to `master`):
1. All CI steps above
2. Authenticate via OIDC (`aws-actions/configure-aws-credentials@v4`)
3. Upload to S3 with per-content-type metadata
4. CloudFront wildcard invalidation

### 7.3 S3 Upload

```bash
# HTML pages (exclude manifest, delete orphans)
aws s3 sync output/ s3://BUCKET/ --delete \
  --exclude "*" --include "*.html" --exclude "_manifest.json" \
  --content-type "text/html; charset=utf-8" \
  --cache-control "public, max-age=300"

# RSS feed (no extension)
aws s3 cp output/feeds/rss s3://BUCKET/feeds/rss \
  --content-type "application/rss+xml; charset=utf-8" \
  --cache-control "public, max-age=900"

# Static assets (long cache, delete orphans)
aws s3 sync output/assets/ s3://BUCKET/assets/ --delete \
  --cache-control "public, max-age=31536000, immutable"

# Root files
aws s3 sync output/ s3://BUCKET/ \
  --exclude "*" --include "robots.txt" --include "favicon.ico"

# Invalidate CloudFront
aws cloudfront create-invalidation \
  --distribution-id DIST_ID --paths "/*"
```

---

## 8. Testing Strategy

### 8.1 Level 1: Unit Tests (`internal/blog/`)

- `loader_test.go` — parse fixture markdown files, verify Post fields, verify indexes built correctly, verify sort order, verify multi-slug indexing
- `slugs_test.go` — C# → c-sharp, C++ → cpp, normal tag slugification
- `carousels_test.go` — posts without images excluded, carousel requires ≥3 posts, related posts exclude source post

### 8.2 Level 2: Generator Tests (`cmd/generate/`)

- `generate_test.go` — run generator against fixture posts, verify:
  - Output directory structure (correct files exist)
  - `_manifest.json` is well-formed (correct pages + redirects)
  - RSS output is valid XML
  - Static assets copied correctly

### 8.3 Level 3: Side-by-Side Comparison (`test/comparison/`)

The key confidence test. Runs against the **real 219 posts** (not fixtures) for full coverage:

1. Run `cmd/generate` against `web/posts/` to produce static output
2. Start `cmd/web` on a local port (background), also using `web/posts/`
3. Read `_manifest.json`
4. For every page: HTTP GET from the Go server, read corresponding file from static output, compare HTML (after normalizing whitespace)
5. For every redirect: verify Go server returns 301 with matching Location header
6. Spider the home page, follow all internal links, verify every discovered URL has a match in the static output

### 8.4 Level 4: Link Integrity E2E (`test/e2e/`)

Validates static output standalone:
- Every `<a href="/...">` in generated HTML points to a file that exists in output
- Every `src="/assets/..."` and `href="/assets/..."` resolves to a real file
- Every tag link in sidebar/posts matches a generated tag page
- RSS `<link>` elements point to generated post pages
- `og:url` meta tags use `https://www.tugberkugurlu.com/` + canonical slug

---

## 9. Implementation Phases

### Phase 1 — Foundation (go-bloggy)

Extract `internal/blog/` package from `cmd/web/main.go`. Add unit tests. Verify `cmd/web` still works identically. This is the prerequisite for everything else and can be merged independently.

### Phase 2 — Generator (go-bloggy)

Build `cmd/generate` using the shared `internal/blog/` package. Add generator-level tests. This is purely additive — nothing in `cmd/web` changes.

### Phase 3 — Verification (go-bloggy)

Build the side-by-side comparison test and E2E link integrity test. Add both to CI. This proves the static output matches what the Go server would serve.

### Phase 4 — Infrastructure (tugberk-infrastructure)

Prerequisite: configure AWS CLI locally for the personal AWS account.

Terraform changes:
- Subscribe to CloudFront flat-rate Free plan ($0/month, includes WAF + DNS)
- New S3 bucket for static site hosting (e.g., `tugberkugurlu-com-static`) with Block Public Access enabled
- CloudFront distribution with OAC pointing to the new bucket
- CloudFront cache policy: ignore query strings, cookies, headers (prevents cache-busting)
- CloudFront Function for redirects (naked domain, legacy images, non-canonical slugs, case folding)
- Reuse existing ACM certificate (covers both `tugberkugurlu.com` and `*.tugberkugurlu.com`)
- WAF rate-based rule: limit per source IP (included in free plan)
- IAM role for GitHub Actions OIDC deploy (scoped to S3 + CloudFront)
- Route53: point `tugberkugurlu.com` and `www.tugberkugurlu.com` to CloudFront
- AWS Budget alarm at $5/month threshold with email notification
- Enable CloudFront standard logging

Do NOT touch the existing `tugberkugurlu-blog` S3 bucket (holds post images).

### Phase 5 — Deploy Pipeline (go-bloggy)

Add GitHub Actions workflows:
- `ci.yml` — build + test + generate + comparison + e2e (on PRs)
- `deploy.yml` — all CI steps + S3 upload + CloudFront invalidation (on push to master)

### Phase 6 — Go Live

Merge the deploy workflow. First deploy goes to S3 + CloudFront. Verify the site works publicly at www.tugberkugurlu.com. Run the side-by-side comparison test against the live site.

### Phase 7 — Cleanup

After the static site is confirmed working publicly:

**go-bloggy repo:**
- Remove `cmd/web/` (the HTTP server is no longer needed)
- Remove Docker-related files (`docker-web.dockerfile`, `docker-compose.yml`)
- Remove GZIP/CaselessMatcher middleware code
- Update CLAUDE.md to reflect new architecture
- Keep `cmd/generate/`, `internal/blog/`, all tests

**tugberk-infrastructure repo:**
- Remove ECS cluster, ECS service, task definitions
- Remove ALB, ALB listeners, ALB target group, ALB security group, ALB ACM resources
- Remove VPC, subnets, NAT Gateway
- Remove ECS IAM roles
- Remove ECS security groups
- Keep: Route53 zone, CloudFront, S3 buckets (both), ACM cert, GitHub Actions IAM

---

## 10. Cost Analysis

### 10.1 Traffic Assumptions

- **3,000 unique page views/day** (~90,000/month)
- Each page view loads: 1 HTML page + ~5 assets (CSS, JS, fonts, images) = ~6 requests
- Total requests: **~540,000/month**
- Per-visit data transfer (with CloudFront GZIP compression):
  - HTML: ~15 KB
  - semantic.min.css: ~100 KB (gzipped from ~700 KB)
  - semantic.min.js: ~80 KB (gzipped from ~340 KB)
  - Theme assets (fonts/icons): ~200 KB
  - Site images: ~100 KB
  - **Total per visit: ~500 KB**
- Monthly data transfer: 3,000/day × 500 KB × 30 = **~45 GB/month**

### 10.2 Recommended: CloudFront Flat-Rate Free Plan ($0/month)

AWS launched flat-rate pricing plans for CloudFront in November 2025. The **Free tier** at $0/month includes:

| Included | Allowance | Our Usage | Headroom |
|---|---|---|---|
| HTTPS Requests | 1,000,000/month | ~540,000/month | 46% spare |
| Data Transfer | 100 GB/month | ~45 GB/month | 55% spare |
| S3 Storage | 5 GB | ~8 MB | Massive spare |
| CloudFront Functions | Included in plan | ~540K invocations/month | No separate billing |
| AWS WAF (basic) | Included | — | DDoS/bot protection |
| Route 53 DNS | Included | 1 hosted zone | Saves $0.50/month |
| ACM Certificate | Free | Wildcard cert | Already exists |

**Key feature: no overage charges.** If traffic exceeds the allowance (viral post, bot traffic), CloudFront throttles performance instead of billing extra. Blocked WAF/DDoS traffic does not count against the allowance.

Up to 3 Free-tier plans per AWS account.

**Note:** Subscribing to the flat-rate Free plan is a **manual console step** — Terraform does not yet support it ([hashicorp/terraform-provider-aws#45450](https://github.com/hashicorp/terraform-provider-aws/issues/45450)). Done once after creating the CloudFront distribution. See tugberkugurlu/tugberk-infrastructure#2 for the full infrastructure plan including all manual steps.

Sources: [CloudFront Pricing](https://aws.amazon.com/cloudfront/pricing/), [Flat-Rate Plans Announcement](https://aws.amazon.com/blogs/networking-and-content-delivery/introducing-flat-rate-pricing-plans-with-no-overages/)

### 10.3 Alternative: Pay-As-You-Go Pricing (for comparison)

| Component | Calculation | Monthly Cost |
|---|---|---|
| CloudFront data transfer (US) | 45 GB × $0.085/GB | $3.83 |
| CloudFront HTTPS requests | 540K × ($0.01/10K) | $0.054 |
| CloudFront Functions | 540K invocations (free tier: 2M) | $0.00 |
| S3 Standard storage | 8 MB × $0.023/GB | ~$0.00 |
| S3 GET requests (origin fetches) | ~5K/month × ($0.0004/1K) | ~$0.00 |
| Route 53 hosted zone | 1 zone | $0.50 |
| ACM certificate | Public cert | $0.00 |
| CloudFront invalidations | ~30/month (free tier: 1000) | $0.00 |
| **Total** | | **~$4.38/month** |

Note: CloudFront always-free tier (1 TB transfer, 10M requests/month) would make this effectively **$0.50/month** (just Route 53). However, the flat-rate Free plan is strictly better because it includes Route 53, WAF, and has no-overage protection.

Sources: [S3 Pricing](https://aws.amazon.com/s3/pricing/), [Route 53 Pricing](https://aws.amazon.com/route53/pricing/)

### 10.4 Cost Comparison: Before vs After

| Component | Current (ECS) | After (S3+CF Free Plan) |
|---|---|---|
| Compute | ECS Fargate 2×(256 CPU/512MB): ~$18/month | $0 |
| Load Balancer | ALB: ~$16/month + LCU charges | $0 (CloudFront) |
| NAT Gateway | 1 gateway: ~$32/month + data | $0 |
| VPC | $0 (but enables above costs) | $0 (no VPC) |
| Data Transfer | ALB/NAT egress | $0 (included in plan) |
| DNS | Route 53: $0.50/month | $0 (included in plan) |
| SSL | ACM: $0 | $0 |
| **Total** | **~$67/month** | **$0/month** |

**Estimated annual savings: ~$800/year.**

### 10.5 Cost Protection Measures

To ensure costs stay at $0 and protect against malicious traffic:

1. **Use CloudFront flat-rate Free plan** — no overages by design; throttling instead of billing under excess load
2. **Configure CloudFront cache policy to ignore query strings, cookies, and headers** — prevents cache-busting attacks that would generate excess S3 origin fetches
3. **Enable S3 Block Public Access** on the static site bucket — only CloudFront (via OAC) can access S3, preventing direct S3 request charges
4. **Set up AWS Budget alarm at $5/month** — early warning if any unexpected charges appear
5. **Add a rate-based WAF rule** — limit requests per source IP (included in free plan WAF)
6. **Enable CloudFront standard logging** — visibility into traffic patterns for anomaly detection

---

## 11. Go Version

The current `go.mod` specifies `go 1.14` and the Dockerfile uses Go 1.15.6. As part of Phase 1 (Foundation), upgrade to a current stable Go version (1.22+) in `go.mod` and CI workflows. This is needed for modern test tooling and language features used in the new code.

---

## 11. Out of Scope

- Changing blog content markdown files
- Migrating images from Azure Blob Storage or the existing S3 bucket
- Adding new features (search, new pages, etc.)
- Changing the visual design or CSS
- Supporting incremental/partial generation (full rebuild is fast enough for 219 posts)

---

## 12. Risks and Mitigations

| Risk | Mitigation |
|---|---|
| URL breakage after migration | Side-by-side comparison test catches mismatches before merge; manifest tracks all URLs |
| Stale CloudFront cache after deploy | Wildcard `/*` invalidation on every deploy |
| RSS feed missing Content-Type | Explicit `--content-type` in S3 upload command; E2E test validates |
| CloudFront Function too large (redirects map) | 20 redirects is tiny; CF Functions support 10KB code |
| `cmd/web` refactoring breaks existing behavior | Unit tests added before refactoring; comparison test validates after |
| AWS CLI not configured locally | Phase 4 begins with AWS CLI setup for personal account |
| Bill shock from DDoS/bot traffic | CloudFront flat-rate Free plan: no overages, throttling only; WAF blocks don't count against allowance |
| S3 direct access bypass | OAC + S3 Block Public Access; bucket only accessible via CloudFront |
| Cache-busting attacks | CloudFront cache policy ignores query strings, cookies, headers |
