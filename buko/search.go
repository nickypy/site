package server

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type SearchEntry struct {
	Location string
	Token    string
}

type SearchIndex struct {
	Entries []SearchEntry
}

func NewSearchIndex() *SearchIndex {
	i := make([]SearchEntry, 1024)

	return &SearchIndex{i}
}

func (s *SearchIndex) IndexHTML(path string) {
	contents := readFile(path)

	reader := bytes.NewReader(contents)
	z := html.NewTokenizer(reader)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}

		if tt == html.TextToken {
			text := z.Token().Data
			tokens := strings.Fields(text)

			for _, t := range tokens {
				token := sanitizeToken(t)

				entry := SearchEntry{
					Location: path,
					Token: token,
				}
				s.Entries = append(s.Entries, entry)
			}
		}
	}
}

func sanitizeToken(token string) string {
	t := strings.ToLower(token)
	t = regexp.MustCompile(`[^a-z0-9 ]+`).ReplaceAllString(t, "")
	return t
}

func Search() {
	index := NewSearchIndex()
	index.IndexHTML("dist/blog/test.html")

	for _, entry := range index.Entries {
		fmt.Println(entry)
	}
}
