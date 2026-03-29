package blog

import "sort"

// TopPicksPostIDs is the ordered list of post IDs that appear in the
// "Top Picks for Software Engineering" carousel on the home and about pages.
// Edit this list to change which posts are featured.
var TopPicksPostIDs = []string{
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

// GetTopPicksCarousel builds the top-picks carousel from the hardcoded list of
// post IDs. Posts that are not found or have no images are silently skipped.
func GetTopPicksCarousel(postsByID map[string]*Post) Carousel {
	posts := make([]*Post, 0, len(TopPicksPostIDs))
	for _, id := range TopPicksPostIDs {
		post, ok := postsByID[id]
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

// GetCarouselForTag returns a carousel of up to 20 image-bearing posts tagged
// with tagSlug, or nil when fewer than 3 qualifying posts exist.
// Posts are expected to already be in newest-first order in postsByTagSlug.
func GetCarouselForTag(tagSlug, title string, postsByTagSlug map[string][]*Post) *Carousel {
	const maxSize = 20
	candidates := make([]*Post, 0, maxSize)
	for _, post := range postsByTagSlug[tagSlug] {
		if !canBeACarouselPost(post) {
			continue
		}
		candidates = append(candidates, post)
		if len(candidates) == maxSize {
			break
		}
	}
	c := Carousel{Title: title, Posts: candidates}
	if !canBeCarousel(c) {
		return nil
	}
	return &c
}

// GetRelatedPostsCarousel returns a carousel of posts related to post by shared
// tags, sorted newest-first. Returns nil when fewer than 3 qualifying posts are
// found. The original post itself is never included.
func GetRelatedPostsCarousel(post *Post, postsByTagSlug map[string][]*Post) *Carousel {
	const maxSize = 20
	seen := make(map[string]bool, maxSize)
	candidates := make([]*Post, 0, maxSize)

	for _, t := range post.Tags {
		if len(candidates) == maxSize {
			break
		}
		for _, candidate := range postsByTagSlug[t.Key] {
			if !canBeRelatedPost(post, candidate, seen) {
				continue
			}
			candidates = append(candidates, candidate)
			seen[candidate.Metadata.ID] = true
			if len(candidates) == maxSize {
				break
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].PublishedOn.Unix() > candidates[j].PublishedOn.Unix()
	})

	c := Carousel{Title: "Related Posts", Posts: candidates}
	if !canBeCarousel(c) {
		return nil
	}
	return &c
}

func canBeRelatedPost(original, candidate *Post, seen map[string]bool) bool {
	return canBeACarouselPost(candidate) &&
		!seen[candidate.Metadata.ID] &&
		candidate.Metadata.ID != original.Metadata.ID
}

func canBeACarouselPost(post *Post) bool {
	return len(post.Images) > 0
}

func canBeCarousel(c Carousel) bool {
	return len(c.Posts) > 2
}
