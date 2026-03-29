package blog

import (
	"fmt"
	"html/template"
	"io"
)

const (
	defaultTitle       = "Tugberk @ the Heart of Software"
	defaultDescription = "Software Engineer and Tech Product enthusiast Tugberk Ugurlu's home on the interwebs! Here, you can find out about Tugberk's conference talks, books and blog posts on software development techniques and practices."
	defaultAdTags      = "software development, asp.net, aws, azure, sql server, dynamodb, elasticsearch, mongodb, .net"
)

// RenderPage executes the named template set against the standard Layout
// context and writes the result to w.
//
// templatePaths must include all page-specific template files AND the
// layout.html file. section is the first URL path component (e.g. "archive")
// used to highlight the active nav item.
func RenderPage(w io.Writer, section string, templatePaths []string, layoutConfig LayoutConfig, tagsList TagCountPairList, data interface{}) error {
	t := template.New("")
	t = t.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})

	var err error
	t, err = t.ParseFiles(templatePaths...)
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}

	title := defaultTitle
	description := defaultDescription
	adTags := defaultAdTags

	if page, ok := data.(Page); ok {
		if page.Title() != "" {
			title = fmt.Sprintf("%s | %s", page.Title(), defaultTitle)
		}
		if page.Description() != "" {
			description = page.Description()
		}
	}

	if postPage, ok := data.(PostPage); ok {
		adTags = postPage.AdTags
	}

	ctx := &Layout{
		Title:       title,
		Description: description,
		Tags:        tagsList,
		Section:     section,
		AdTags:      adTags,
		Config:      layoutConfig,
		Data: PageData{
			Data:   data,
			Config: layoutConfig,
		},
	}
	return t.ExecuteTemplate(w, "layout", ctx)
}
