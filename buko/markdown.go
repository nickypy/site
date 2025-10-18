package server

import (
	"bytes"
	"fmt"
	"log"
	p "path"
	"time"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"gopkg.in/yaml.v2"
)

const THEME = "github-dark"

type BlogMetadata struct {
	Title       string    `yaml:"title"`
	IsPublished bool      `yaml:"published"`
	Date        time.Time `yaml:"date"`
	Slug        string    `yaml:"slug"`
	Tags        []string  `yaml:"tags"`
}

func CreateNewMarkdownFile(title string, path string) {
	header, err := yaml.Marshal(BlogMetadata{
		Title:       title,
		IsPublished: false,
		Date:        time.Now(),
		Slug:        "",
		Tags:        make([]string, 0),
	})

	if err != nil {
		log.Fatalln(err)
	}

	contents := fmt.Sprintf("---\n%s---\n", header)

	filename := makeFilename(title) + ".md"
	out := p.Join(path, filename)

	writeFile(out, []byte(contents))
}

type MarkdownRenderer struct {
	md goldmark.Markdown
}

func NewMarkdownRenderer() MarkdownRenderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle(THEME),
			),
			extension.GFM,
			&anchor.Extender{
				Position: anchor.After,
				Texter:   anchor.Text("#"),
			},
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	return MarkdownRenderer{
		md: md,
	}
}

func (mr *MarkdownRenderer) Render(filepath string) (map[string]interface{}, string) {
	contents := readFile(filepath)

	var body bytes.Buffer
	ctx := parser.NewContext()
	err := mr.md.Convert(contents, &body, parser.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	return meta.Get(ctx), body.String()
}
