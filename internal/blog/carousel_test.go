package blog

import (
	"testing"
	"time"
)

// buildTestPost is a helper that creates a minimal Post for carousel tests.
func buildTestPost(id, slug string, hasImage bool, published time.Time, tags []string) *Post {
	images := []string{}
	if hasImage {
		images = []string{"https://example.com/" + slug + ".jpg"}
	}
	var tagPairs []TagCountPair
	for _, t := range tags {
		tagPairs = append(tagPairs, TagCountPair{Key: t, Value: &Tag{Name: t, Count: 1}})
	}
	return &Post{
		Metadata:    PostMetadata{ID: id, Slugs: []string{slug}, Tags: tags},
		Images:      images,
		PublishedOn: published,
		Tags:        tagPairs,
	}
}

var (
	t1 = time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)
	t2 = time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)
	t3 = time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)
	t4 = time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)
)

func TestCanBeACarouselPost_RequiresImage(t *testing.T) {
	withImg := buildTestPost("a", "a", true, t1, nil)
	noImg := buildTestPost("b", "b", false, t1, nil)
	if !canBeACarouselPost(withImg) {
		t.Error("post with image should be a carousel candidate")
	}
	if canBeACarouselPost(noImg) {
		t.Error("post without image should not be a carousel candidate")
	}
}

func TestCanBeCarousel_RequiresThreePosts(t *testing.T) {
	p1 := buildTestPost("1", "p1", true, t1, nil)
	p2 := buildTestPost("2", "p2", true, t2, nil)
	p3 := buildTestPost("3", "p3", true, t3, nil)

	if canBeCarousel(Carousel{Posts: []*Post{p1}}) {
		t.Error("1 post should not form a valid carousel")
	}
	if canBeCarousel(Carousel{Posts: []*Post{p1, p2}}) {
		t.Error("2 posts should not form a valid carousel")
	}
	if !canBeCarousel(Carousel{Posts: []*Post{p1, p2, p3}}) {
		t.Error("3 posts should form a valid carousel")
	}
}

func TestGetTopPicksCarousel_IncludesKnownIDs(t *testing.T) {
	postsByID := map[string]*Post{}
	for i, id := range TopPicksPostIDs {
		hasImage := i%2 == 0 // half with images
		postsByID[id] = buildTestPost(id, "slug-"+id, hasImage, t1, nil)
	}
	c := GetTopPicksCarousel(postsByID)
	if c.Title != "Top Picks for Software Engineering" {
		t.Errorf("unexpected title: %q", c.Title)
	}
	// Only posts with images should be included
	for _, post := range c.Posts {
		if len(post.Images) == 0 {
			t.Errorf("carousel post %q has no image", post.Metadata.ID)
		}
	}
}

func TestGetTopPicksCarousel_SkipsMissingPosts(t *testing.T) {
	// Empty map — no posts found
	c := GetTopPicksCarousel(map[string]*Post{})
	if len(c.Posts) != 0 {
		t.Errorf("expected 0 posts for empty map, got %d", len(c.Posts))
	}
}

func TestGetCarouselForTag_ReturnedNilWhenFewerThanThreePosts(t *testing.T) {
	p1 := buildTestPost("1", "p1", true, t1, nil)
	p2 := buildTestPost("2", "p2", true, t2, nil)
	postsByTagSlug := map[string][]*Post{"go": {p1, p2}}
	c := GetCarouselForTag("go", "Go Posts", postsByTagSlug)
	if c != nil {
		t.Error("expected nil carousel for fewer than 3 posts")
	}
}

func TestGetCarouselForTag_SkipsPostsWithoutImages(t *testing.T) {
	withImg1 := buildTestPost("1", "p1", true, t1, nil)
	withImg2 := buildTestPost("2", "p2", true, t2, nil)
	withImg3 := buildTestPost("3", "p3", true, t3, nil)
	noImg := buildTestPost("4", "p4", false, t4, nil)

	postsByTagSlug := map[string][]*Post{
		"go": {withImg1, noImg, withImg2, withImg3},
	}
	c := GetCarouselForTag("go", "Go Posts", postsByTagSlug)
	if c == nil {
		t.Fatal("expected non-nil carousel")
	}
	for _, post := range c.Posts {
		if len(post.Images) == 0 {
			t.Errorf("carousel contains post %q with no image", post.Metadata.ID)
		}
	}
}

func TestGetCarouselForTag_RespectsMaxSize(t *testing.T) {
	posts := make([]*Post, 25)
	for i := range posts {
		posts[i] = buildTestPost(string(rune('a'+i)), "slug", true, t1, nil)
	}
	c := GetCarouselForTag("go", "Go Posts", map[string][]*Post{"go": posts})
	if c != nil && len(c.Posts) > 20 {
		t.Errorf("carousel exceeds max size: got %d posts", len(c.Posts))
	}
}

func TestGetRelatedPostsCarousel_ExcludesOriginalPost(t *testing.T) {
	original := buildTestPost("orig", "original", true, t1, []string{"go"})
	rel1 := buildTestPost("r1", "r1", true, t2, []string{"go"})
	rel2 := buildTestPost("r2", "r2", true, t3, []string{"go"})
	rel3 := buildTestPost("r3", "r3", true, t4, []string{"go"})

	postsByTagSlug := map[string][]*Post{
		"go": {original, rel1, rel2, rel3},
	}
	c := GetRelatedPostsCarousel(original, postsByTagSlug)
	if c == nil {
		t.Fatal("expected non-nil carousel")
	}
	for _, post := range c.Posts {
		if post.Metadata.ID == "orig" {
			t.Error("original post should not appear in its own related posts carousel")
		}
	}
}

func TestGetRelatedPostsCarousel_DeduplicatesCandidates(t *testing.T) {
	// p1 shares both "go" and "architecture" tags with the original
	// It should appear only once in the carousel.
	original := buildTestPost("orig", "original", true, t1, []string{"go", "architecture"})
	p1 := buildTestPost("p1", "p1", true, t2, []string{"go", "architecture"})
	p2 := buildTestPost("p2", "p2", true, t3, []string{"go"})
	p3 := buildTestPost("p3", "p3", true, t4, []string{"architecture"})

	postsByTagSlug := map[string][]*Post{
		"go":           {original, p1, p2},
		"architecture": {original, p1, p3},
	}
	c := GetRelatedPostsCarousel(original, postsByTagSlug)
	if c == nil {
		t.Fatal("expected non-nil carousel")
	}
	seen := map[string]int{}
	for _, post := range c.Posts {
		seen[post.Metadata.ID]++
	}
	for id, count := range seen {
		if count > 1 {
			t.Errorf("post %q appears %d times in carousel (want 1)", id, count)
		}
	}
}

func TestGetRelatedPostsCarousel_SortedNewestFirst(t *testing.T) {
	original := buildTestPost("orig", "original", true, time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), []string{"go"})
	older := buildTestPost("older", "older", true, t4, []string{"go"})
	newer := buildTestPost("newer", "newer", true, t1, []string{"go"})
	middle := buildTestPost("middle", "middle", true, t2, []string{"go"})

	postsByTagSlug := map[string][]*Post{
		"go": {original, newer, middle, older},
	}
	c := GetRelatedPostsCarousel(original, postsByTagSlug)
	if c == nil {
		t.Fatal("expected non-nil carousel")
	}
	for i := 1; i < len(c.Posts); i++ {
		if c.Posts[i].PublishedOn.After(c.Posts[i-1].PublishedOn) {
			t.Errorf("carousel not sorted newest-first at index %d", i)
		}
	}
}

func TestGetRelatedPostsCarousel_NilWhenFewerThanThreeRelated(t *testing.T) {
	original := buildTestPost("orig", "original", true, t1, []string{"go"})
	p1 := buildTestPost("p1", "p1", true, t2, []string{"go"})

	postsByTagSlug := map[string][]*Post{
		"go": {original, p1},
	}
	c := GetRelatedPostsCarousel(original, postsByTagSlug)
	if c != nil {
		t.Error("expected nil carousel when fewer than 3 related posts")
	}
}
