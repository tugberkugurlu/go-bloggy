package htmltextextractor

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"strings"
)

var textTags = []string{
	"a",
	"p", "span", "em", "string", "blockquote", "q", "cite",
	"h1", "h2", "h3", "h4", "h5", "h6",
}

// https://kananrahimov.com/post/golang-html-tokenizer-extract-text-from-a-web-page/
func Extract(htmlBody []byte) ([]string, error) {
	tag := ""
	enter := false
	tokenizer := html.NewTokenizer(bytes.NewReader(htmlBody))
	var result []string
	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			return nil, err
		case html.StartTagToken, html.SelfClosingTagToken:
			enter = false

			tag = token.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					break
				}
			}
		case html.TextToken:
			if enter {
				data := strings.TrimSpace(token.Data)

				if len(data) > 0 {
					result = append(result, data)
				}
			}
		}
	}
	return result, nil
}