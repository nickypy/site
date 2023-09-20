package server

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func Search() {
	contents := readFile("./dist/blog/test.html")

	reader := bytes.NewReader(contents)
	z := html.NewTokenizer(reader)

	index := make(map[string]int)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}

		if tt == html.TextToken {
			text := z.Token().Data
			tokens := strings.Fields(text)

			for _, token := range tokens {
				if _, ok := index[token]; ok {
					index[token]++
				} else {
					index[token] = 1
				}
			}
		}
	}

	for k, v := range index {
		fmt.Println(k, v)
	}
}
