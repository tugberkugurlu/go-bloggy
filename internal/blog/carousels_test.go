package blog

import (
	"testing"
	"time"
)

func TestCanBeACarouselPost_WithImages(t *testing.T) {
	post := &Post{
		Images: []string{"https://example.com/img.png"},
	}
	if !CanBeACarouselPost(post) {
		t.Error("expected post with images to be carousel-eligible")
	}
}

func TestCanBeACarouselPost_WithoutImages(t *testing.T) {
	post := &Post{
		Images: nil,
	}
	if CanBeACarouselPost(post) {
		t.Error("expected post without images to NOT be carousel-eligible")
	}
}

func TestCanBeCarousel_ThreeOrMorePosts(t *testing.T) {
	c := Carousel{
		Posts: []*Post{
			{Images: []string{"a"}},
			{Images: []string{"b"}},
			{Images: []string{"c"}},
		},
	}
	if !CanBeCarousel(c) {
		t.Error("expected carousel with 3 posts to be valid")
	}
}

func TestCanBeCarousel_FewerThanThreePosts(t *testing.T) {
	c := Carousel{
		Posts: []*Post{
			{Images: []string{"a"}},
			{Images: []string{"b"}},
		},
	}
	if CanBeCarousel(c) {
		t.Error("expected carousel with 2 posts to NOT be valid")
	}
}

func TestCanBeCarousel_EmptyPosts(t *testing.T) {
	c := Carousel{
		Posts: nil,
	}
	if CanBeCarousel(c) {
		t.Error("expected carousel with 0 posts to NOT be valid")
	}
}

func TestGetCarouselForTag_ExcludesPostsWithoutImages(t *testing.T) {
	postsByTagSlug := map[string][]*Post{
		"go": {
			{Metadata: PostMetadata{ID: "1"}, Images: []string{"img1.png"}},
			{Metadata: PostMetadata{ID: "2"}, Images: nil},
			{Metadata: PostMetadata{ID: "3"}, Images: []string{"img3.png"}},
			{Metadata: PostMetadata{ID: "4"}, Images: []string{"img4.png"}},
		},
	}

	carousel := GetCarouselForTag("go", "Go Posts", postsByTagSlug)
	if carousel == nil {
		t.Fatal("expected non-nil carousel")
	}
	// Should have 3 posts (1, 3, 4) — post 2 has no images
	if len(carousel.Posts) != 3 {
		t.Errorf("expected 3 posts in carousel, got %d", len(carousel.Posts))
	}
	for _, p := range carousel.Posts {
		if len(p.Images) == 0 {
			t.Errorf("carousel contains post %q with no images", p.Metadata.ID)
		}
	}
}

func TestGetCarouselForTag_ReturnsNilWhenFewerThanThreeEligible(t *testing.T) {
	postsByTagSlug := map[string][]*Post{
		"go": {
			{Metadata: PostMetadata{ID: "1"}, Images: []string{"img1.png"}},
			{Metadata: PostMetadata{ID: "2"}, Images: []string{"img2.png"}},
		},
	}

	carousel := GetCarouselForTag("go", "Go Posts", postsByTagSlug)
	if carousel != nil {
		t.Error("expected nil carousel when fewer than 3 eligible posts")
	}
}

func TestGetRelatedPostsCarousel_ExcludesSourcePost(t *testing.T) {
	sourcePost := &Post{
		Metadata: PostMetadata{ID: "source"},
		Images:   []string{"img.png"},
		Tags: []TagCountPair{
			{Key: "go"},
		},
	}

	relatedPost1 := &Post{
		Metadata:    PostMetadata{ID: "related-1"},
		Images:      []string{"img1.png"},
		PublishedOn: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	relatedPost2 := &Post{
		Metadata:    PostMetadata{ID: "related-2"},
		Images:      []string{"img2.png"},
		PublishedOn: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	}
	relatedPost3 := &Post{
		Metadata:    PostMetadata{ID: "related-3"},
		Images:      []string{"img3.png"},
		PublishedOn: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
	}

	postsByTagSlug := map[string][]*Post{
		"go": {sourcePost, relatedPost1, relatedPost2, relatedPost3},
	}

	carousel := GetRelatedPostsCarousel(sourcePost, postsByTagSlug)
	if carousel == nil {
		t.Fatal("expected non-nil carousel")
	}

	for _, p := range carousel.Posts {
		if p.Metadata.ID == "source" {
			t.Error("related posts carousel should NOT contain the source post")
		}
	}
	if len(carousel.Posts) != 3 {
		t.Errorf("expected 3 related posts, got %d", len(carousel.Posts))
	}
}

func TestGetRelatedPostsCarousel_NoDuplicates(t *testing.T) {
	sharedPost := &Post{
		Metadata:    PostMetadata{ID: "shared"},
		Images:      []string{"img.png"},
		PublishedOn: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	other1 := &Post{
		Metadata:    PostMetadata{ID: "other1"},
		Images:      []string{"img.png"},
		PublishedOn: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	}
	other2 := &Post{
		Metadata:    PostMetadata{ID: "other2"},
		Images:      []string{"img.png"},
		PublishedOn: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
	}

	sourcePost := &Post{
		Metadata: PostMetadata{ID: "source"},
		Images:   []string{"img.png"},
		Tags: []TagCountPair{
			{Key: "tag1"},
			{Key: "tag2"},
		},
	}

	// sharedPost appears in both tags
	postsByTagSlug := map[string][]*Post{
		"tag1": {sharedPost, other1, other2},
		"tag2": {sharedPost, other1, other2},
	}

	carousel := GetRelatedPostsCarousel(sourcePost, postsByTagSlug)
	if carousel == nil {
		t.Fatal("expected non-nil carousel")
	}

	seen := make(map[string]bool)
	for _, p := range carousel.Posts {
		if seen[p.Metadata.ID] {
			t.Errorf("duplicate post %q in related posts carousel", p.Metadata.ID)
		}
		seen[p.Metadata.ID] = true
	}
}
