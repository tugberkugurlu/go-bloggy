package blog

import (
	"strings"

	"github.com/gosimple/slug"
)

// preDefinedSlugs maps special tag names to their URL-safe slug equivalents.
var preDefinedSlugs = map[string]string{
	"c#":  "c-sharp",
	"c++": "cpp",
}

// ToSlug converts a tag name into a URL-safe slug. It handles predefined
// overrides (e.g., "C#" -> "c-sharp", "C++" -> "cpp") before falling
// through to the generic slug library.
func ToSlug(tag string) string {
	tagToSlugify := tag
	if v, ok := preDefinedSlugs[strings.ToLower(tag)]; ok {
		tagToSlugify = v
	}
	return slug.Make(tagToSlugify)
}
