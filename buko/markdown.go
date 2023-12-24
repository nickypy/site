package server

import (
	"fmt"
	"log"
	p "path"
	"time"

	"gopkg.in/yaml.v2"
)

type BlogMetadata struct {
	Title       string    `yaml:"title"`
	IsPublished bool      `yaml:"published"`
	Date        time.Time `yaml:"date"`
	Tags        []string  `yaml:"tags"`
}

func CreateNewMarkdownFile(title string, path string) {
	header, err := yaml.Marshal(BlogMetadata{
		Title:       title,
		IsPublished: false,
		Date:        time.Now(),
		Tags:        make([]string, 0),
	})

	if err != nil {
		log.Fatalln(err)
	}

	contents := fmt.Sprintf("---\n%s---\n", header)

	filename := makeURLSlug(title) + ".md"
	out := p.Join(path, filename)

	writeFile(out, []byte(contents))
}
