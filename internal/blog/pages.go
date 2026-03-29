package blog

import "fmt"

// LayoutConfig holds per-request configuration injected into every page layout.
type LayoutConfig struct {
	GoogleAnalyticsTrackingId string
	AssetsUrl                 string
}

// PageData wraps the page-specific data alongside the layout configuration
// for use inside the shared layout template.
type PageData struct {
	Data   interface{}
	Config LayoutConfig
}

// Layout is the top-level template context passed to layout.html.
type Layout struct {
	Title       string
	Description string
	Tags        TagCountPairList
	Section     string
	Data        PageData
	AdTags      string
	Config      LayoutConfig
}

// Home is the view model for the home page (/).
type Home struct {
	TopCarousel        Carousel
	Posts              []*Post
	SpeakingActivities []*SpeakingActivity
}

// Blog is the view model for the post archive page (/archive).
type Blog struct {
	Posts []*Post
}

// TagsPage is the view model for a tag filter page (/tags/{slug}).
type TagsPage struct {
	Tag   *Tag
	Posts []*Post
}

func (t TagsPage) Title() string {
	return t.Tag.Name
}

func (t TagsPage) Description() string {
	base := fmt.Sprintf("Tugberk's thoughts on the topic of '%s'", t.Tag.Name)
	if t.Tag.Count > 1 {
		return fmt.Sprintf("%s, through %d blog posts", base, t.Tag.Count)
	}
	return base
}

// PostPage is the view model for an individual blog post (/archive/{slug}).
type PostPage struct {
	Post                 *Post
	RelatedPostsCarousel *Carousel
	AdTags               string
}

func (p PostPage) Title() string       { return p.Post.Metadata.Title }
func (p PostPage) Description() string { return p.Post.Metadata.Abstract }

// SpeakingPage is the view model for the speaking page (/speaking).
type SpeakingPage struct {
	SpeakingActivities []*SpeakingActivity
	GeekTalksCarousel  *Carousel
}

func (s SpeakingPage) Title() string { return "Tugberk Ugurlu Public Speaking Engagements" }
func (s SpeakingPage) Description() string {
	return "Tugberk speaks at conferences on technical leadership, software architecture, lean software development, Microsoft Azure and .NET."
}

// AboutPage is the view model for the about page (/about).
type AboutPage struct {
	GeekTalksCarousel *Carousel
	TopCarousel       Carousel
}

// ContactPage is the view model for the contact page (/contact).
type ContactPage struct {
	ArchitectureCarousel *Carousel
}

var _ Page = (*SpeakingPage)(nil)
var _ Page = (*PostPage)(nil)
var _ Page = (*TagsPage)(nil)
