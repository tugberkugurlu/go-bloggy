package main

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

func getTopPicksCarousel(postsByIdMap map[string]*Post) Carousel {
	posts := make([]*Post, 0, len(topPicksCarouselPostIDs))
	for _, id := range topPicksCarouselPostIDs {
		post, ok := postsByIdMap[id]
		if !ok || !canBeCarousel(post) {
			continue
		}
		posts = append(posts, post)
	}
	return Carousel{
		Title: "Top Picks for Software Engineering",
		Posts: posts,
	}
}

func canBeCarousel(post *Post) bool {
	return len(post.Images) > 0
}