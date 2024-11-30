package server

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

const BASE_TEMPLATE = "/templates/base.tmpl"
const BLOG_POST_TEMPLATE = "/templates/blog_post.tmpl"
const BLOG_PAGE_TEMPLATE = "/templates/blog.tmpl"
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
	markdown              MarkdownRenderer
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

func NewBlogRenderer(prefix string, outputPath string, opts ...BlogOption) *BlogRenderCache {
	template := NewTemplateRenderer(prefix + BASE_TEMPLATE)
	template.AddTemplate(
		prefix+BLOG_POST_TEMPLATE,
		prefix+BLOG_PAGE_TEMPLATE,
	)

	md := NewMarkdownRenderer()
	links := NewLinks()

	blog := &BlogRenderCache{
		Prefix:                prefix,
		OutputPath:            outputPath,
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
		path.Join(b.Prefix, "assets"),
		path.Join(b.OutputPath),
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
	b.renderPage()
	b.generateAtomFeed()
}

func (b *BlogRenderCache) renderMarkdown(filepath string) {
	frontmatter, body := b.markdown.Render(filepath)

	var metadata BlogMetadata
	out, err := yaml.Marshal(frontmatter)
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

	post := BlogPost{
		metadata.Slug,
		metadata.Slug + ".html",
		body,
		metadata,
	}

	b.mutex.Lock()
	b.Items = append(b.Items, post)
	b.mutex.Unlock()

	var blogPost bytes.Buffer
	b.template.Render(
		&blogPost,
		BLOG_POST_TEMPLATE,
		PostTemplateArgs{
			post.Metadata.Title,
			post.Body,
			b.Links,
		},
	)

	writeFile(path.Join(b.OutputPath, post.Path), blogPost.Bytes())
}

func (b *BlogRenderCache) renderPage() {
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
		BLOG_PAGE_TEMPLATE,
		IndexTemplateArgs{
			Title:     "nickypy: blog",
			Body:      body.String(),
			BlogItems: publishableItems,
			Links:     b.Links,
		},
	)

	output := body.Bytes()
	writeFile(path.Join(b.OutputPath, "index.html"), output)
}
