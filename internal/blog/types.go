package blog

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// Config holds site configuration loaded from config.yaml.
type Config struct {
	AssetsUrl    string
	AssetsPrefix string
}

// IsLocal returns true if the assets URL is a relative path.
func (c Config) IsLocal() bool {
	return strings.Index(c.AssetsUrl, "/") == 0
}

// FullAssetsUrl returns the complete assets URL including prefix.
func (c Config) FullAssetsUrl() string {
	if c.AssetsPrefix == "" {
		return c.AssetsUrl
	}
	return fmt.Sprintf("%s/%s",
		strings.TrimSuffix(c.AssetsUrl, "/"),
		strings.TrimSuffix(c.AssetsPrefix, "/"),
	)
}

// LayoutConfig holds configuration used in template rendering.
type LayoutConfig struct {
	GoogleAnalyticsTrackingId string
	AssetsUrl                 string
}

// Post is the runtime representation of a parsed markdown blog post.
type Post struct {
	Body      template.HTML
	Highlight string
	Images    []string
	Tags      []TagCountPair
	Metadata  PostMetadata

	ReadingTimeDisplay string
	ReadingTime        *time.Duration

	PublishedOn             time.Time
	PublishedOnDisplay      string
	PublishedOnDisplayBrief string
}

// PostMetadata holds the YAML front matter from a markdown file.
type PostMetadata struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Abstract  string   `yaml:"abstract"`
	Format    string   `yaml:"format"`
	CreatedOn string   `yaml:"created_at"`
	Tags      []string `yaml:"tags"`
	Slugs     []string `yaml:"slugs"`
}

// Tag represents a blog tag with its usage count.
type Tag struct {
	Name  string
	Count int
}

// TagCountPair pairs a tag slug with its Tag.
type TagCountPair struct {
	Key   string
	Value *Tag
}

// TagCountPairList is a sortable list of TagCountPair, sorted by tag count.
type TagCountPairList []TagCountPair

func (p TagCountPairList) Len() int           { return len(p) }
func (p TagCountPairList) Less(i, j int) bool {
	if p[i].Value.Count != p[j].Value.Count {
		return p[i].Value.Count < p[j].Value.Count
	}
	// Tiebreaker: alphabetical by slug key for deterministic ordering
	return p[i].Key < p[j].Key
}
func (p TagCountPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Carousel holds a titled collection of posts for carousel display.
type Carousel struct {
	Title string
	Posts []*Post
}

// SpeakingActivity represents a public speaking engagement.
type SpeakingActivity struct {
	Title       string
	Activity    string
	ImageURL    string
	City        string
	Country     string
	DisplayDate string
	Extras      []SpeakingActivityExtra

	EmbededHTMLData template.HTML
}

// SpeakingActivityExtra holds supplemental links for a speaking activity.
type SpeakingActivityExtra struct {
	Name           string
	Link           string
	IconCSSClasses string
}

// Page is the interface for page view models that provide title and description.
type Page interface {
	Title() string
	Description() string
}

// PageData wraps page-specific data with layout config for template rendering.
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

// Home is the view model for the home page.
type Home struct {
	TopCarousel        Carousel
	Posts              []*Post
	SpeakingActivities []*SpeakingActivity
}

// Blog is the view model for the archive page.
type Blog struct {
	Posts []*Post
}

// TagsPage is the view model for a tag-filtered listing page.
type TagsPage struct {
	Tag   *Tag
	Posts []*Post
}

// Title returns the tag name as page title.
func (t TagsPage) Title() string {
	return t.Tag.Name
}

// Description returns a descriptive string for the tag page.
func (t TagsPage) Description() string {
	baseTitle := fmt.Sprintf("Tugberk's thoughts on the topic of '%s'", t.Tag.Name)
	if t.Tag.Count > 1 {
		return fmt.Sprintf("%s, through %d blog posts", baseTitle, t.Tag.Count)
	}
	return baseTitle
}

// PostPage is the view model for an individual blog post page.
type PostPage struct {
	Post                 *Post
	RelatedPostsCarousel *Carousel
	AdTags               string
}

// Title returns the post title.
func (p PostPage) Title() string {
	return p.Post.Metadata.Title
}

// Description returns the post abstract.
func (p PostPage) Description() string {
	return p.Post.Metadata.Abstract
}

// SpeakingPage is the view model for the speaking engagements page.
type SpeakingPage struct {
	SpeakingActivities []*SpeakingActivity
	GeekTalksCarousel  *Carousel
}

// Title returns the speaking page title.
func (s SpeakingPage) Title() string {
	return "Tugberk Ugurlu Public Speaking Engagements"
}

// Description returns the speaking page description.
func (s SpeakingPage) Description() string {
	return "Tugberk speaks at conferences on technical leadership, software architecture, lean software development, Microsoft Azure and .NET."
}

// AboutPage is the view model for the about page.
type AboutPage struct {
	GeekTalksCarousel *Carousel
	TopCarousel       Carousel
}

// ContactPage is the view model for the contact page.
type ContactPage struct {
	ArchitectureCarousel *Carousel
}

// Compile-time checks that page view models implement Page.
var _ Page = (*SpeakingPage)(nil)
var _ Page = (*PostPage)(nil)
var _ Page = (*TagsPage)(nil)

// Site holds all loaded site data: posts, indexes, carousels, config, etc.
type Site struct {
	Config         Config
	LayoutConfig   LayoutConfig
	Posts          []*Post
	PostsByID      map[string]*Post
	PostsBySlug    map[string]*Post
	PostsByTagSlug map[string][]*Post
	TagsBySlug     map[string]*Tag
	TagsList       TagCountPairList
	Carousels      []Carousel
	Speaking       []*SpeakingActivity
}
