# S3 Static Hosting Migration Plan

## Goal

Replace the always-on Go HTTP server with a static site generator whose
output can be hosted on Amazon S3 + CloudFront, giving:

- Zero-idle cost (S3 + CloudFront are billed per request/byte, not uptime)
- HTTPS with custom domain (`www.tugberkugurlu.com`)
- Full URL parity with the existing server — no broken links, no SEO regression
- Reproducible, CI-verified builds
- Smart CloudFront invalidation (only changed pages re-invalidated)

The existing Go server and all markdown source files are **preserved**; the
static generator is an additive new command (`cmd/staticgen`) that runs at
CI time and produces a `dist/` directory.

---

## Architecture: Before → After

### Current (Go HTTP server on AWS)

```
Browser → Route53 → Load Balancer → EC2/ECS (Go server)
                                          ↓
                                    In-memory post index
                                    (built at startup from web/posts/*.md)
```

### Target (S3 + CloudFront)

```
Browser → Route53 (ALIAS) → CloudFront Distribution
                                  ↓               ↓
                          CloudFront Function   S3 Bucket (origin)
                          (URL rewrites /       (dist/ directory
                           redirects /           uploaded by CI)
                           www enforcement /
                           case-normalise)
```

CI pipeline generates `dist/` from markdown files and uploads to S3 on every
push to `master`. CloudFront serves the files globally with HTTPS.

---

## URL Parity Analysis

Every URL the current server handles must be preserved:

| Current URL | Static file | Method |
|---|---|---|
| `/` | `dist/index.html` | S3 index document |
| `/archive` | `dist/archive/index.html` | S3 sub-directory index |
| `/archive/{primary-slug}` | `dist/archive/{slug}/index.html` | S3 sub-directory index |
| `/archive/{alt-slug}` → 301 → canonical | `_redirects.json` entry | CloudFront Function |
| `/tags/{tag-slug}` | `dist/tags/{slug}/index.html` | S3 sub-directory index |
| `/about` | `dist/about/index.html` | S3 sub-directory index |
| `/speaking` | `dist/speaking/index.html` | S3 sub-directory index |
| `/contact` | `dist/contact/index.html` | S3 sub-directory index |
| `/feeds/rss` | `dist/feeds/rss` (no extension) | Content-Type: application/rss+xml via S3 metadata |
| `/content/images/*` → 301 → Azure Blob | CloudFront Function rule | CloudFront Function |
| `/assets/*` | `dist/assets/*` | S3 files |
| `/robots.txt`, `/favicon.ico` | `dist/robots.txt`, `dist/favicon.ico` | S3 files |
| `tugberkugurlu.com` → `www.tugberkugurlu.com` | CloudFront Function | CloudFront Function |

**S3 Static Website Hosting** handles directory-to-index-document resolution
(`/archive` → `/archive/index.html`) when the index document is set to
`index.html`.

**Note:** Use S3 website endpoint (not REST API endpoint) as the CloudFront
origin so that index document resolution works correctly.

---

## Redirect Strategy

Three categories of redirects need to be preserved:

### 1. Secondary slug → canonical slug (per-post 301s)

Many posts have multiple slugs (`slugs` YAML array). Only `slugs[0]` is
canonical. The generator writes these into `_redirects.json`.

**Implementation:** CloudFront Function reads `_redirects.json` at deploy time
(embedded as a constant) and returns 301 for matching paths.

Alternative: Generate an S3 object for each secondary slug with the
`x-amz-website-redirect-location` metadata header. Simpler but adds objects
to S3.

**Recommended:** S3 redirect metadata — one object per secondary slug,
Content-Type `text/html`, S3 metadata:
```
x-amz-website-redirect-location: https://www.tugberkugurlu.com/archive/{primary-slug}
```
The deploy script generates these during `aws s3 sync`.

### 2. Legacy image URLs (`/content/images/*`)

Redirect to Azure Blob Storage. Handled by a CloudFront Function (a wildcard
rule cannot be expressed as an S3 routing rule).

### 3. `tugberkugurlu.com` → `www.tugberkugurlu.com`

Handled by a separate S3 bucket configured as a redirect-only website, or by a
CloudFront Function that intercepts requests to the bare domain.

**Recommended:** Route53 + a second CloudFront distribution (or a redirect S3
bucket) for `tugberkugurlu.com`.

---

## CloudFront Function (Viewer Request)

A single JavaScript CloudFront Function handles all runtime URL logic:

```javascript
function handler(event) {
    var request = event.request;
    var uri = request.uri;

    // 1. Enforce www
    if (request.headers.host.value === 'tugberkugurlu.com') {
        return {
            statusCode: 301,
            headers: { location: { value: 'https://www.tugberkugurlu.com' + uri } }
        };
    }

    // 2. Normalise to lowercase
    var lower = uri.toLowerCase();
    if (lower !== uri) {
        return { statusCode: 301, headers: { location: { value: lower } } };
    }

    // 3. Legacy image redirect
    if (uri.startsWith('/content/images/')) {
        var blobBase = 'https://tugberkugurlu.blob.core.windows.net/bloggyimages/legacy-blog-images/images';
        return {
            statusCode: 301,
            headers: { location: { value: blobBase + uri.substring('/content/images'.length) } }
        };
    }

    // 4. Append index.html for extensionless directory paths
    if (!uri.includes('.') || uri.endsWith('/')) {
        if (!uri.endsWith('/')) uri += '/';
        request.uri = uri + 'index.html';
    }

    return request;
}
```

**Note:** The secondary-slug redirects are better handled as S3 redirect
metadata (approach in §Redirect Strategy §1) rather than embedding a large
lookup table in the CloudFront Function (which has a 10 KB size limit on
Key-Value Store-backed functions).

---

## S3 + CloudFront Infrastructure

### AWS Resources

| Resource | Configuration |
|---|---|
| **S3 Bucket** `www.tugberkugurlu.com` | Static website hosting enabled; index document: `index.html`; error document: `404.html`; NOT publicly accessible directly (blocked by bucket policy); accessible only via CloudFront OAC |
| **ACM Certificate** | Region: `us-east-1` (required for CloudFront); covers `www.tugberkugurlu.com` and `tugberkugurlu.com`; validated via Route53 DNS records |
| **CloudFront Distribution** | Origin: S3 website endpoint; Alternate domain names: both www and bare; ACM cert attached; Default root object: `index.html`; Viewer protocol: redirect HTTP → HTTPS; Custom error responses: 403→`/404.html`, 404→`/404.html` |
| **CloudFront Function** | Viewer Request event; handles www redirect, lowercase normalisation, legacy image redirects, extensionless URL rewriting |
| **Route53** | `www.tugberkugurlu.com` → ALIAS to CloudFront; `tugberkugurlu.com` → ALIAS to CloudFront (function redirects to www) |

### Infrastructure as Code

Terraform (or AWS CDK) definitions will live in a new `infrastructure/`
directory. This is a separate concern from the Go codebase and will be added in
a dedicated phase.

---

## ETag-based Deployment and Invalidation

S3 ETags are the MD5 hash of an object's content. The deploy script:

1. Dry-run `aws s3 sync dist/ s3://BUCKET/ --delete --dryrun` to identify
   which files would change.
2. Perform the real sync.
3. Collect the list of changed paths from the sync output.
4. Create a CloudFront invalidation **only for changed paths** rather than a
   blanket `/*`.

This keeps invalidation costs low (first 1,000 paths/month free) and avoids
unnecessary cache evictions for unchanged content.

```bash
#!/usr/bin/env bash
set -euo pipefail

BUCKET="s3://www.tugberkugurlu.com"
CF_DIST_ID="${CLOUDFRONT_DISTRIBUTION_ID}"

# Sync and capture which files actually changed
CHANGED=$(aws s3 sync dist/ "$BUCKET" --delete --output text 2>&1 \
  | grep "^upload:" \
  | awk '{print $4}' \
  | sed "s|$BUCKET||")

if [ -z "$CHANGED" ]; then
  echo "No files changed."
  exit 0
fi

# Create invalidation for changed paths only
PATHS=$(echo "$CHANGED" | tr '\n' ' ')
aws cloudfront create-invalidation \
  --distribution-id "$CF_DIST_ID" \
  --paths $PATHS
```

---

## Side-by-Side Comparison Testing

The existing Go server and the static generator must produce equivalent
content. A comparison test runs both servers and validates they match:

### Approach

1. Start the Go server on `localhost:8080`.
2. Start an HTTP file server (or `aws s3 website` equivalent) serving `dist/`
   on `localhost:8081`.
3. For each URL in the known URL set, fetch from both servers and compare the
   HTML (normalised — whitespace collapsed, dynamic elements like timestamps
   ignored).
4. Report any differences as test failures.

### Implementation (future phase)

A new `cmd/compare` command (or `internal/e2e` test package) will:

- Accept `-live-url` and `-static-url` flags.
- Enumerate all URLs from the blog index.
- Issue parallel HTTP GETs and diff the results.
- Run as a step in CI after static site generation.

### What to compare

| Check | How |
|---|---|
| HTTP 200 for all primary URLs | Status code assertion |
| 301 for secondary slugs | Follow redirect, compare destination |
| RSS feed is valid XML with ≥20 items | XML parse + item count |
| Post HTML contains post title and abstract | Substring match |
| Tag pages list the correct posts | Post title presence check |
| `_redirects.json` references only real slugs | Cross-reference with index |

---

## Implementation Phases

### Phase 0 — Refactoring & Unit Tests ✅ (this PR)

Extract shared blog logic from `cmd/web/main.go` into `internal/blog/` so that
both the web server and the static generator share the same data loading,
carousel building, and template rendering code.

**Deliverables:**
- `internal/blog/` package: `models.go`, `pages.go`, `loader.go`,
  `carousel.go`, `render.go`, `speaking.go`
- `internal/blog/*_test.go` with unit tests for loading and carousel logic
- `cmd/staticgen/main.go` skeleton that generates the full `dist/` directory
- Updated `cmd/web/main.go` using `internal/blog`
- Updated `.github/workflows/go.yml` with static site generation + validation

### Phase 1 — Static Generator Hardening

- Generate a `404.html` page.
- Write secondary-slug S3 redirect metadata upload logic into the deploy script.
- Write a `cmd/staticgen` integration test that generates into a temp directory
  and checks the output structure, RSS validity, and a sample post's HTML.
- Bump Go module version to `1.21` and update Docker image.

### Phase 2 — Side-by-Side Comparison Tests

- Implement `cmd/compare` (or `internal/e2e`).
- Add a CI job that: starts the Go server, generates static output, runs
  `cmd/compare`, and fails the build on any mismatch.
- Establish a baseline tolerance list for known acceptable differences (e.g.
  Disqus counter in dynamic server vs. static page).

### Phase 3 — AWS Infrastructure

- Create `infrastructure/` with Terraform (or CDK) modules for:
  - S3 buckets (www + redirect-only bare domain)
  - ACM certificate
  - CloudFront distributions
  - Route53 records
  - CloudFront Function
- Document IAM permissions required for the CI deploy role.

### Phase 4 — CI/CD Deployment Pipeline

- Add a GitHub Actions `deploy` job that runs on push to `master`:
  1. Build static site.
  2. Run comparison tests against production Go server.
  3. Sync to S3 with ETag-based invalidation.
  4. Create targeted CloudFront invalidation.
- Keep the Go server running in parallel during transition (shadow mode).

### Phase 5 — DNS Cutover

- Validate static site in production behind a staging subdomain.
- Switch Route53 records to CloudFront.
- Monitor error rates for 24 hours.
- Decommission Go server after confidence period.

---

## Open Questions / Decisions

1. **Terraform vs CDK**: CDK allows type-safe infrastructure in Go (matching
   the repo's language), but Terraform has wider operator familiarity. Decision
   needed before Phase 3.

2. **Google Analytics migration**: The current site uses Universal Analytics
   (UA-) which is sunset. Migrate to GA4 when rebuilding templates.

3. **Disqus comments**: Disqus works identically on static pages (pure JS).
   No change needed.

4. **Contact form**: The current `/contact` page embeds a Google Form. This
   works identically on a static page.

5. **CDN for assets**: Currently `/assets` is served locally. For production,
   point `assets_url` at a CDN (e.g. CloudFront with a separate S3 origin for
   assets) to benefit from edge caching. The `assets_url` config already
   supports this.

6. **`feeds/rss` Content-Type**: S3 serves files without extensions with
   `application/octet-stream`. The deploy script must set
   `--content-type "application/rss+xml"` for the `feeds/rss` object via
   `aws s3 cp` (after the sync).

7. **Secondary slug redirect approach**: S3 object metadata redirect vs.
   CloudFront Function with embedded lookup table. S3 metadata approach
   scales better (no CloudFront Function size limit concern) but requires the
   deploy script to create/maintain these objects explicitly.
