package server

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type SearchEntry struct {
	Location string
	Token    string
}

type SearchIndex struct {
	BasePath string
	Entries  []SearchEntry
}

func NewSearchIndex() *SearchIndex {
	i := make([]SearchEntry, 1024)

	return &SearchIndex{"", i}
}

func (s *SearchIndex) WithPath(path string) *SearchIndex {
	s.BasePath = path
	return s
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
					Token:    token,
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

func (s *SearchIndex) IndexAll() *SearchIndex {
	_ = filepath.Walk(s.BasePath, func(p string, fi os.FileInfo, err error) (e error) {
		if !fi.IsDir() {
			if strings.HasSuffix(p, ".html") {
				s.IndexHTML(p)
			}
		}

		return nil
	})

	return s
}

func (s *SearchIndex) Invert() *InvertedIndex {
	// TODO: implement ranking by occurence

	i := NewInvertedIndex()

	for _, entry := range s.Entries {
		i.Index(entry.Token, entry.Location)
	}

	return i
}

type InvertedIndexNode map[string]*InvertedIndex

type InvertedIndex struct {
	Nodes InvertedIndexNode `json:"n"`
	Leaf  *string           `json:"l"`
}

func NewInvertedIndex() *InvertedIndex {
	nodes := make(map[string]*InvertedIndex)
	return &InvertedIndex{nodes, nil}
}

func separate(s string) (string, string) {
	r := []rune(s)
	if len(r) == 0 {
		return "", ""
	}
	return string(r[0]), string(r[1:])
}

func (i *InvertedIndex) Index(token string, location string) {
	first, rest := separate(token)
	if first == "" {
		i.Leaf = &location
		return
	}

	if node, ok := i.Nodes[first]; ok {
		node.Index(rest, location)
	} else {
		i.Nodes[first] = NewInvertedIndex()
		i.Nodes[first].Index(rest, location)
	}
}

func (i *InvertedIndex) Search(term string) []string {
	results := i.search(term, make([]string, 0))
	return results
}

func (i *InvertedIndex) search(term string, acc []string) []string {
	first, rest := separate(term)

	if leaf := i.Leaf; leaf != nil && first == "" {
		acc = append(acc, *leaf)
	}

	if node, ok := i.Nodes[first]; ok && first != "" {
		results := node.search(rest, acc)
		acc = append(acc, results...)
	}

	return acc
}

func Search(term string) {
	index := NewSearchIndex().WithPath("dist/blog").IndexAll().Invert()
	results := index.Search(term)

	fmt.Printf("%v\n", results)
}
