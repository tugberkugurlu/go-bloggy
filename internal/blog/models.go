package blog

import (
	"html/template"
	"time"
)

// Post is the runtime representation of a single markdown blog post.
type Post struct {
	Body                    template.HTML
	Highlight               string
	Images                  []string
	Tags                    []TagCountPair
	Metadata                PostMetadata
	ReadingTimeDisplay      string
	ReadingTime             *time.Duration
	PublishedOn             time.Time
	PublishedOnDisplay      string
	PublishedOnDisplayBrief string
}

// PostMetadata holds the YAML front matter of a post file.
type PostMetadata struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Abstract  string   `yaml:"abstract"`
	Format    string   `yaml:"format"`
	CreatedOn string   `yaml:"created_at"`
	Tags      []string `yaml:"tags"`
	Slugs     []string `yaml:"slugs"`
}

// Tag represents a blog tag with its usage count across all posts.
type Tag struct {
	Name  string
	Count int
}

// TagCountPair associates a URL slug with its Tag.
type TagCountPair struct {
	Key   string // URL-safe slug derived from Tag.Name
	Value *Tag
}

// TagCountPairList is a sortable slice of TagCountPair (by tag count).
type TagCountPairList []TagCountPair

func (p TagCountPairList) Len() int           { return len(p) }
func (p TagCountPairList) Less(i, j int) bool { return p[i].Value.Count < p[j].Value.Count }
func (p TagCountPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Carousel is a horizontally-scrollable collection of featured posts.
// Requires at least 3 image-bearing posts to be considered valid.
type Carousel struct {
	Title string
	Posts []*Post
}

// Page is implemented by page view models that supply their own SEO metadata.
type Page interface {
	Title() string
	Description() string
}

// SpeakingActivity represents a single public speaking engagement.
type SpeakingActivity struct {
	Title           string
	Activity        string
	ImageURL        string
	City            string
	Country         string
	DisplayDate     string
	Extras          []SpeakingActivityExtra
	EmbededHTMLData template.HTML
}

// SpeakingActivityExtra is a labelled link attached to a speaking activity
// (e.g. slides, video, conference profile).
type SpeakingActivityExtra struct {
	Name           string
	Link           string
	IconCSSClasses string
}
