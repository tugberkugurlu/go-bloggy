package main

import (
	"sort"
	"strings"
)

var topPicksCarouselPostIDs = []string{
	"01ESYWJ3NHDFF5G7H4JKCY3DS0",
	"01E972TEACJE3TQW6Q59C0X4NJ",
	"4ea0c96b-06cd-4bb3-73a5-08d6999e9406",
	"01EJ1FYNB4CB3B13JVAEHB5908",
	"e48899fd-cbc7-4384-8e98-833e08a7c071",
	"ae0648d8-929f-4ddf-b2de-5564451a3754",
	"e4d7b9ac-7f3e-4bb8-8f5e-232d5b897d3f",
	"e019b369-0c8f-4415-8704-70f67b866bf5",
	"55fa0e98-4a31-48e2-9473-919b82fb9232",
	"370ecc73-3661-4462-91ac-179b9f8a3c51",
}

func GetCarousels(postsByIdMap map[string]*Post) []Carousel {
	var carousels []Carousel
	carousels = append(carousels, getTopPicksCarousel(postsByIdMap))
	return carousels
}

func GetCarouselForTag(tagSlug string, title string, postsByTagSlugMap map[string][]*Post) *Carousel {
	maxSize := 20
	tagPosts := postsByTagSlugMap[tagSlug]
	candidates := make([]*Post, 0, maxSize)
	for _, post := range tagPosts {
		if !canBeACarouselPost(post) {
			continue
		}
		candidates = append(candidates, post)
		if len(candidates) == maxSize {
			break
		}
	}
	carousel := Carousel{
		Title: title,
		Posts: candidates,
	}
	if !canBeCarousel(carousel) {
		return nil
	}
	return &carousel
}

func GetRelatedPostsCarousel(post *Post, postsByTagSlugMap map[string][]*Post) *Carousel {
	maxSize := 20
	candidatesMap := make(map[string]bool, maxSize)
	candidates := make([]*Post, 0, maxSize)
	for _, t := range post.Tags {
		if len(candidates) == maxSize {
			break
		}
		tagPosts := postsByTagSlugMap[t.Key]
		for i := 0; i < len(tagPosts); i++ {
			tagPost := tagPosts[i]
			if !canBeARelatedPostsCarouselPost(post, tagPost, candidatesMap) {
				continue
			}
			candidates = append(candidates, tagPost)
			candidatesMap[tagPost.Metadata.ID] = true
			if len(candidates) == maxSize {
				break
			}
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].PublishedOn.Unix() > candidates[j].PublishedOn.Unix()
	})
	carousel := Carousel{
		Title: "Related Posts",
		Posts: candidates,
	}
	if !canBeCarousel(carousel) {
		return nil
	}
	return &carousel
}

func getTopPicksCarousel(postsByIdMap map[string]*Post) Carousel {
	posts := make([]*Post, 0, len(topPicksCarouselPostIDs))
	for _, id := range topPicksCarouselPostIDs {
		post, ok := postsByIdMap[id]
		if !ok || !canBeACarouselPost(post) {
			continue
		}
		posts = append(posts, post)
	}
	return Carousel{
		Title: "Top Picks for Software Engineering",
		Posts: posts,
	}
}

func canBeARelatedPostsCarouselPost(originalPost *Post, post *Post, currentCandidates map[string]bool) bool {
	if !canBeACarouselPost(post) {
		return false
	}
	if currentCandidates[post.Metadata.ID] {
		return false
	}
	if strings.EqualFold(post.Metadata.ID, originalPost.Metadata.ID) {
		return false
	}
	return true
}

func canBeACarouselPost(post *Post) bool {
	return len(post.Images) > 0
}

func canBeCarousel(carousel Carousel) bool {
	return len(carousel.Posts) > 2
}