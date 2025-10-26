package server

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	p "path"
	"strings"
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

const THEME = "modus-vivendi"
const MD_PATH = "markdown"

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

func (mr *MarkdownRenderer) RenderBytes(contents []byte) (map[string]any, string) {
	var body bytes.Buffer
	ctx := parser.NewContext()
	err := mr.md.Convert(contents, &body, parser.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	return meta.Get(ctx), body.String()
}


func listFiles(directory string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	files := make([]string, 0)

	for _, entry := range entries {
		canon := path.Join(directory, entry.Name())

		if entry.IsDir() {
			inner := listFiles(canon)
			files = append(files, inner...)
			continue
		}

		if strings.HasSuffix(canon, ".md") {
			files = append(files, canon)
		}
	}

	return files
}

type MarkdownFileEntry struct {
	Content []byte
	Key     string
	Prefix  string
}

func GetAllMarkdownFiles() []MarkdownFileEntry {
	files := make([]MarkdownFileEntry, 0)
	for _, file := range listFiles(MD_PATH) {

		content := readFile(file)
		key := p.Base(file)

		prefix := strings.TrimPrefix(
			strings.TrimSuffix(file, key),
			MD_PATH+"/",
		)

		prefix = strings.TrimSuffix(prefix, "/")

		key = strings.TrimSuffix(key, ".md")

		entry := MarkdownFileEntry{
			content, key, prefix,
		}

		files = append(files, entry)
	}

	return files
}
