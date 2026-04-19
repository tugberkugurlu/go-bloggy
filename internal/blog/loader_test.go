package blog

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

func loadTestSite(t *testing.T) *Site {
	t.Helper()
	dir := testdataDir()
	site, err := LoadSite(LoadSiteConfig{
		PostsDir:   dir,
		ConfigPath: dir,
	})
	if err != nil {
		t.Fatalf("LoadSite failed: %v", err)
	}
	return site
}

func TestLoadSite_ParsesPosts(t *testing.T) {
	site := loadTestSite(t)

	if len(site.Posts) != 4 {
		t.Fatalf("expected 4 posts, got %d", len(site.Posts))
	}
}

func TestLoadSite_SortOrder_NewestFirst(t *testing.T) {
	site := loadTestSite(t)

	for i := 1; i < len(site.Posts); i++ {
		if site.Posts[i].PublishedOn.After(site.Posts[i-1].PublishedOn) {
			t.Errorf("posts not sorted newest first: post[%d] (%s) is after post[%d] (%s)",
				i, site.Posts[i].PublishedOn, i-1, site.Posts[i-1].PublishedOn)
		}
	}
}

func TestLoadSite_PostsBySlug(t *testing.T) {
	site := loadTestSite(t)

	// post1 has two slugs
	post1ByCanonical, ok := site.PostsBySlug["first-test-post"]
	if !ok {
		t.Fatal("expected to find post by canonical slug 'first-test-post'")
	}
	post1ByAlias, ok := site.PostsBySlug["first-test-post-alias"]
	if !ok {
		t.Fatal("expected to find post by alias slug 'first-test-post-alias'")
	}
	if post1ByCanonical != post1ByAlias {
		t.Error("expected both slugs to resolve to the same Post pointer")
	}
	if post1ByCanonical.Metadata.ID != "test-post-1" {
		t.Errorf("expected post ID 'test-post-1', got %q", post1ByCanonical.Metadata.ID)
	}
}

func TestLoadSite_PostsByID(t *testing.T) {
	site := loadTestSite(t)

	post, ok := site.PostsByID["test-post-2"]
	if !ok {
		t.Fatal("expected to find post by ID 'test-post-2'")
	}
	if post.Metadata.Title != "Second Test Post" {
		t.Errorf("expected title 'Second Test Post', got %q", post.Metadata.Title)
	}
}

func TestLoadSite_TagsIndexed(t *testing.T) {
	site := loadTestSite(t)

	// "Go" tag should be present (used by post1 and post2)
	goTag, ok := site.TagsBySlug["go"]
	if !ok {
		t.Fatal("expected 'go' tag to be indexed")
	}
	if goTag.Count != 2 {
		t.Errorf("expected 'go' tag count 2, got %d", goTag.Count)
	}

	// "C#" -> "c-sharp" tag should be present (used by post3)
	csharpTag, ok := site.TagsBySlug["c-sharp"]
	if !ok {
		t.Fatal("expected 'c-sharp' tag to be indexed")
	}
	if csharpTag.Count != 1 {
		t.Errorf("expected 'c-sharp' tag count 1, got %d", csharpTag.Count)
	}

	// "C++" -> "cpp" tag should be present (used by post4)
	cppTag, ok := site.TagsBySlug["cpp"]
	if !ok {
		t.Fatal("expected 'cpp' tag to be indexed")
	}
	if cppTag.Count != 1 {
		t.Errorf("expected 'cpp' tag count 1, got %d", cppTag.Count)
	}
}

func TestLoadSite_PostsByTagSlug(t *testing.T) {
	site := loadTestSite(t)

	goPosts, ok := site.PostsByTagSlug["go"]
	if !ok {
		t.Fatal("expected 'go' tag in PostsByTagSlug")
	}
	if len(goPosts) != 2 {
		t.Errorf("expected 2 posts with 'go' tag, got %d", len(goPosts))
	}
}

func TestLoadSite_TagsList(t *testing.T) {
	site := loadTestSite(t)

	if len(site.TagsList) == 0 {
		t.Fatal("expected non-empty TagsList")
	}

	// TagsList should be sorted by count descending
	for i := 1; i < len(site.TagsList); i++ {
		if site.TagsList[i].Value.Count > site.TagsList[i-1].Value.Count {
			t.Errorf("TagsList not sorted descending: tag[%d] count %d > tag[%d] count %d",
				i, site.TagsList[i].Value.Count, i-1, site.TagsList[i-1].Value.Count)
		}
	}
}

func TestLoadSite_PostFields(t *testing.T) {
	site := loadTestSite(t)

	// Find post1
	post, ok := site.PostsByID["test-post-1"]
	if !ok {
		t.Fatal("expected to find test-post-1")
	}

	if post.Metadata.Title != "First Test Post" {
		t.Errorf("expected title 'First Test Post', got %q", post.Metadata.Title)
	}
	if post.Metadata.Abstract != "This is the abstract for the first test post." {
		t.Errorf("unexpected abstract: %q", post.Metadata.Abstract)
	}
	if post.Metadata.Format != "md" {
		t.Errorf("expected format 'md', got %q", post.Metadata.Format)
	}
	if len(post.Metadata.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(post.Metadata.Tags))
	}
	if len(post.Metadata.Slugs) != 2 {
		t.Errorf("expected 2 slugs, got %d", len(post.Metadata.Slugs))
	}
	if len(post.Images) != 2 {
		t.Errorf("expected 2 images, got %d", len(post.Images))
	}
	if post.PublishedOnDisplay != "2021-06-15 10:00:00" {
		t.Errorf("unexpected PublishedOnDisplay: %q", post.PublishedOnDisplay)
	}
	if post.PublishedOnDisplayBrief != "15 June 2021" {
		t.Errorf("unexpected PublishedOnDisplayBrief: %q", post.PublishedOnDisplayBrief)
	}
	if post.Body == "" {
		t.Error("expected non-empty Body")
	}
	if post.ReadingTimeDisplay == "" {
		t.Error("expected non-empty ReadingTimeDisplay")
	}
}

func TestLoadSite_PostWithNoImages(t *testing.T) {
	site := loadTestSite(t)

	post, ok := site.PostsByID["test-post-4"]
	if !ok {
		t.Fatal("expected to find test-post-4")
	}
	if len(post.Images) != 0 {
		t.Errorf("expected 0 images, got %d", len(post.Images))
	}
}

func TestLoadSite_Config(t *testing.T) {
	site := loadTestSite(t)

	if site.Config.AssetsUrl != "/assets" {
		t.Errorf("expected AssetsUrl '/assets', got %q", site.Config.AssetsUrl)
	}
}
