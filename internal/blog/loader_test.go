package blog

import (
	"strings"
	"testing"
	"time"
)

const testPostsDir = "testdata/posts"

func TestLoadIndex_PostCount(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	if got, want := len(idx.Posts), 4; got != want {
		t.Errorf("post count = %d, want %d", got, want)
	}
}

func TestLoadIndex_SortedNewestFirst(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	for i := 1; i < len(idx.Posts); i++ {
		prev := idx.Posts[i-1].PublishedOn
		curr := idx.Posts[i].PublishedOn
		if curr.After(prev) {
			t.Errorf("posts[%d] (%s) is newer than posts[%d] (%s): not sorted newest-first",
				i, curr.Format(time.RFC3339), i-1, prev.Format(time.RFC3339))
		}
	}
}

func TestLoadIndex_PostsBySlug(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}

	cases := []struct {
		slug   string
		wantID string
	}{
		{"test-post-alpha", "test-post-alpha-id"},
		{"test-post-alpha-old-slug", "test-post-alpha-id"}, // alternate slug
		{"test-post-beta", "test-post-beta-id"},
		{"test-post-gamma", "test-post-gamma-id"},
	}
	for _, tc := range cases {
		post, ok := idx.PostsBySlug[tc.slug]
		if !ok {
			t.Errorf("slug %q not found in PostsBySlug", tc.slug)
			continue
		}
		if post.Metadata.ID != tc.wantID {
			t.Errorf("PostsBySlug[%q].ID = %q, want %q", tc.slug, post.Metadata.ID, tc.wantID)
		}
	}
}

func TestLoadIndex_PostsByID(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	post, ok := idx.PostsByID["test-post-alpha-id"]
	if !ok {
		t.Fatal("test-post-alpha-id not found in PostsByID")
	}
	if post.Metadata.Title != "Test Post Alpha" {
		t.Errorf("title = %q, want %q", post.Metadata.Title, "Test Post Alpha")
	}
}

func TestLoadIndex_TagCounting(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	cases := []struct {
		slug      string
		wantCount int
	}{
		// "Go" appears in alpha, beta, and delta = 3
		{"go", 3},
		// "Testing" appears in alpha and gamma = 2
		{"testing", 2},
		// "Architecture" appears in beta and gamma = 2
		{"architecture", 2},
	}
	for _, tc := range cases {
		tag, ok := idx.TagsBySlug[tc.slug]
		if !ok {
			t.Errorf("tag slug %q not found", tc.slug)
			continue
		}
		if tag.Count != tc.wantCount {
			t.Errorf("tag %q count = %d, want %d", tc.slug, tag.Count, tc.wantCount)
		}
	}
}

func TestLoadIndex_PostsByTagSlug(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	goPosts := idx.PostsByTagSlug["go"]
	if len(goPosts) != 3 {
		t.Errorf("go tag post count = %d, want 3", len(goPosts))
	}
}

func TestLoadIndex_TagsListOrderedByCountDesc(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	for i := 1; i < len(idx.TagsList); i++ {
		prev := idx.TagsList[i-1].Value.Count
		curr := idx.TagsList[i].Value.Count
		if curr > prev {
			t.Errorf("TagsList not sorted descending: [%d].Count=%d > [%d].Count=%d",
				i, curr, i-1, prev)
		}
	}
}

func TestLoadIndex_MarkdownRenderedToHTML(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	post, ok := idx.PostsByID["test-post-alpha-id"]
	if !ok {
		t.Fatal("test-post-alpha-id not found")
	}
	body := string(post.Body)
	if !strings.Contains(body, "<h2") {
		t.Errorf("markdown not rendered to HTML: body does not contain <h2>, got: %s", body[:200])
	}
}

func TestLoadIndex_ImagesExtracted(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	post, ok := idx.PostsByID["test-post-alpha-id"]
	if !ok {
		t.Fatal("test-post-alpha-id not found")
	}
	if len(post.Images) == 0 {
		t.Error("expected at least one image extracted from alpha post")
	}
	if post.Images[0] != "https://example.com/alpha.jpg" {
		t.Errorf("Images[0] = %q, want %q", post.Images[0], "https://example.com/alpha.jpg")
	}
}

func TestLoadIndex_NoImagesForDelta(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	post, ok := idx.PostsByID["test-post-delta-id"]
	if !ok {
		t.Fatal("test-post-delta-id not found")
	}
	if len(post.Images) != 0 {
		t.Errorf("delta should have no images, got %v", post.Images)
	}
}

func TestLoadIndex_ReadingTimeSet(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	post, ok := idx.PostsByID["test-post-alpha-id"]
	if !ok {
		t.Fatal("test-post-alpha-id not found")
	}
	if post.ReadingTime == nil {
		t.Error("expected ReadingTime to be set")
	}
	if post.ReadingTimeDisplay == "" {
		t.Error("expected ReadingTimeDisplay to be non-empty")
	}
}

func TestLoadIndex_PublishedOnParsed(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	post, ok := idx.PostsByID["test-post-alpha-id"]
	if !ok {
		t.Fatal("test-post-alpha-id not found")
	}
	want := time.Date(2021, 6, 20, 9, 0, 0, 0, time.UTC)
	if !post.PublishedOn.Equal(want) {
		t.Errorf("PublishedOn = %v, want %v", post.PublishedOn, want)
	}
	if post.PublishedOnDisplayBrief != "20 June 2021" {
		t.Errorf("PublishedOnDisplayBrief = %q, want %q", post.PublishedOnDisplayBrief, "20 June 2021")
	}
}

func TestLoadIndex_TagsAssignedToPost(t *testing.T) {
	idx, err := LoadIndex(testPostsDir)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	post, ok := idx.PostsByID["test-post-alpha-id"]
	if !ok {
		t.Fatal("test-post-alpha-id not found")
	}
	if len(post.Tags) != 2 {
		t.Fatalf("expected 2 tags on alpha, got %d: %v", len(post.Tags), post.Tags)
	}
	tagKeys := map[string]bool{}
	for _, tp := range post.Tags {
		tagKeys[tp.Key] = true
	}
	if !tagKeys["go"] {
		t.Error("expected 'go' tag on alpha")
	}
	if !tagKeys["testing"] {
		t.Error("expected 'testing' tag on alpha")
	}
}

func TestToSlug(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Go", "go"},
		{"C#", "c-sharp"},
		{"c#", "c-sharp"},
		{"C++", "cpp"},
		{"c++", "cpp"},
		{".NET", "net"},
		{"ASP.NET MVC", "asp-net-mvc"},
		{"Distributed Systems", "distributed-systems"},
	}
	for _, tc := range cases {
		got := ToSlug(tc.input)
		if got != tc.want {
			t.Errorf("ToSlug(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestGeneratePostURL(t *testing.T) {
	post := &Post{Metadata: PostMetadata{Slugs: []string{"my-slug", "old-slug"}}}
	got := GeneratePostURL(post)
	want := "https://www.tugberkugurlu.com/archive/my-slug"
	if got != want {
		t.Errorf("GeneratePostURL = %q, want %q", got, want)
	}
}
