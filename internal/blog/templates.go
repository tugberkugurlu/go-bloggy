package blog

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

// TemplateFuncMap returns the template function map used by all blog templates.
func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"mod": func(i, j int) bool { return i%j == 0 },
	}
}

// RenderPage renders a page using the given template files and data,
// writing the result to w. The section is derived from the URL path
// (e.g., "archive" from "/archive/my-post"). layoutPath is the path
// to layout.html. tagsList is the global tag list for the sidebar.
func RenderPage(w io.Writer, layoutConfig LayoutConfig, templatePaths []string, layoutPath string, section string, tagsList TagCountPairList, data interface{}) error {
	t := template.New("")
	t = t.Funcs(TemplateFuncMap())
	allPaths := make([]string, len(templatePaths)+1)
	copy(allPaths, templatePaths)
	allPaths[len(templatePaths)] = layoutPath
	t, err := t.ParseFiles(allPaths...)
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}

	pageTitle := "Tugberk @ the Heart of Software"
	pageDescription := "Software Engineer and Tech Product enthusiast Tugberk Ugurlu's home on the interwebs! Here, you can find out about Tugberk's conference talks, books and blog posts on software development techniques and practices."
	page, ok := data.(Page)
	if ok {
		if page.Title() != "" {
			pageTitle = fmt.Sprintf("%s | %s", page.Title(), pageTitle)
		}
		if page.Description() != "" {
			pageDescription = page.Description()
		}
	}

	adTags := "software development, asp.net, aws, azure, sql server, dynamodb, elasticsearch, mongodb, .net"
	postPage, ok := data.(PostPage)
	if ok {
		adTags = postPage.AdTags
	}

	pageContext := &Layout{
		Title:       pageTitle,
		Description: pageDescription,
		Tags:        tagsList,
		Section:     section,
		AdTags:      adTags,
		Config:      layoutConfig,

		Data: PageData{
			Data:   data,
			Config: layoutConfig,
		},
	}
	return t.ExecuteTemplate(w, "layout", pageContext)
}

// SectionFromPath extracts the top-level section from a URL path.
// For example, "/archive/my-post" returns "archive", "/" returns "".
// The path must start with "/" (as r.URL.Path always does).
func SectionFromPath(urlPath string) string {
	if urlPath == "" || urlPath == "/" {
		return ""
	}
	section := urlPath[1:]
	index := strings.Index(section, "/")
	if index != -1 {
		section = urlPath[1 : index+1]
	}
	return section
}

// GeneratePostURL returns the canonical URL for a post.
func GeneratePostURL(post *Post) string {
	return fmt.Sprintf("https://www.tugberkugurlu.com/archive/%s", post.Metadata.Slugs[0])
}
