package server

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

const BASE_TEMPLATE = "/templates/base.tmpl"
const BLOG_TEMPLATE = "/templates/blog.tmpl"
const INDEX_TEMPLATE = "/templates/index.tmpl"
const BLOG_BASE_PATH = "blog"

type BlogOption func(*BlogRenderCache)

func RenderUnpublished() BlogOption {
	return func(b *BlogRenderCache) {
		b.ShouldListUnpublished = true
	}
}

type BlogRenderCache struct {
	Prefix                string
	OutputPath            string
	Items                 []BlogPost
	ShouldListUnpublished bool
	Links                 LinkMetadata
	markdown              goldmark.Markdown
	template              *TemplateRenderer
	mutex                 *sync.RWMutex
}

type BlogPost struct {
	Slug     string
	Path     string
	Body     string
	Metadata BlogMetadata
}

func (post *BlogPost) FormatDate() string {
	return post.Metadata.Date.Format("2006 Jan")
}

func NewBlogRenderCache(prefix string, opts ...BlogOption) *BlogRenderCache {
	template := NewTemplateRenderer(prefix + BASE_TEMPLATE)
	template.AddTemplate(
		prefix+BLOG_TEMPLATE,
		prefix+INDEX_TEMPLATE,
	)

	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
			extension.GFM,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)

	links := NewLinks()

	blog := &BlogRenderCache{
		Prefix:                prefix,
		OutputPath:            "dist",
		Items:                 nil,
		ShouldListUnpublished: false,
		Links:                 links,
		markdown:              md,
		template:              template,
		mutex:                 new(sync.RWMutex),
	}

	for _, opt := range opts {
		opt(blog)
	}

	return blog
}

func (b *BlogRenderCache) Render() {
	var wg sync.WaitGroup

	b.Items = make([]BlogPost, 0)

	CopyDirectory(
		path.Join(b.Prefix, "static"),
		path.Join(b.OutputPath, "static"),
	)

	_ = filepath.Walk(b.Prefix+"/markdown", func(path string, fi os.FileInfo, err error) (e error) {
		if !fi.IsDir() {
			if strings.HasSuffix(path, ".md") {
				wg.Add(1)
				go func(path string) {
					defer wg.Done()
					b.renderMarkdown(path)
				}(path)
			}
		}
		return nil
	})

	wg.Wait()
	b.renderIndex()
}

func (b *BlogRenderCache) renderMarkdown(filepath string) {
	md := readFile(filepath)

	var body bytes.Buffer
	ctx := parser.NewContext()
	err := b.markdown.Convert(md, &body, parser.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	var metadata BlogMetadata
	out, err := yaml.Marshal(meta.Get(ctx))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(out, &metadata)
	if err != nil {
		panic(err)
	}

	if b.ShouldListUnpublished {
		metadata.IsPublished = true
	}

	slug := makeURLSlug(metadata.Title)
	post := BlogPost{
		slug,
		path.Join(BLOG_BASE_PATH, slug) + ".html",
		body.String(),
		metadata,
	}

	b.mutex.Lock()
	b.Items = append(b.Items, post)
	b.mutex.Unlock()

	var blogPost bytes.Buffer
	b.template.Render(
		&blogPost,
		BLOG_TEMPLATE,
		PostTemplateArgs{
			post.Metadata.Title,
			post.Body,
			b.Links,
		},
	)

	writeFile(path.Join(b.OutputPath, post.Path), blogPost.Bytes())
}

func (b *BlogRenderCache) renderIndex() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	sort.Slice(b.Items, func(i, j int) bool {
		return b.Items[j].Metadata.Date.Unix() < b.Items[i].Metadata.Date.Unix()
	})

	var body bytes.Buffer
	publishableItems := make([]BlogPost, 0)
	for _, i := range b.Items {
		if i.Metadata.IsPublished {
			publishableItems = append(publishableItems, i)
		}
	}

	b.template.Render(
		&body,
		INDEX_TEMPLATE,
		IndexTemplateArgs{
			Title:    "nickypy: blog",
			Body:      body.String(),
			BlogItems: publishableItems,
			Links:     b.Links,
		},
	)

	output := body.Bytes()
	writeFile(path.Join(b.OutputPath, "index.html"), output)
}
