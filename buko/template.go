package server

import (
	"bytes"
	"io"
	p "path"
	"text/template"
)

const BASE_TEMPLATE = "./templates/base.tmpl"
const BLOG_POST_TEMPLATE = "./templates/blog_post.tmpl"
const BLOG_PAGE_TEMPLATE = "./templates/blog.tmpl"
const PAGE_TEMPLATE = "./templates/page.tmpl"

type TemplateRenderer struct {
	baseTmpl string
	tmpl     map[string]*template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{BASE_TEMPLATE, make(map[string]*template.Template)}
}

func (t *TemplateRenderer) WithDefault() *TemplateRenderer {
	t.AddTemplate(
		BLOG_POST_TEMPLATE,
		BLOG_PAGE_TEMPLATE,
		PAGE_TEMPLATE,
	)

	return t
}

func (t *TemplateRenderer) AddTemplate(paths ...string) {
	for _, path := range paths {
		name := p.Base(path)
		t.tmpl[name] = template.Must(template.ParseFiles(path, t.baseTmpl))
	}
}

func (t *TemplateRenderer) Render(out io.Writer, path string, data any) {
	name := p.Base(path)
	err := t.tmpl[name].ExecuteTemplate(out, "base", data)
	if err != nil {
		panic(err)
	}
}

func (t *TemplateRenderer) RenderPage(p PageTemplateArgs) []byte {
	var out bytes.Buffer
	t.Render(&out, PAGE_TEMPLATE, p)
	return out.Bytes()
}

func (t *TemplateRenderer) RenderBlogPost(p PageTemplateArgs) []byte {
	var out bytes.Buffer
	t.Render(&out, BLOG_POST_TEMPLATE, p)
	return out.Bytes()
}

func (t *TemplateRenderer) RenderBlogPage(p BlogPageTemplateArgs) []byte {
	var out bytes.Buffer
	t.Render(&out, BLOG_PAGE_TEMPLATE, p)
	return out.Bytes()
}


type BlogPageTemplateArgs struct {
	Nav       NavBar
	Title     string
	Body      string
	BlogItems []BlogPost
	Links     LinkMetadata
}

type PageTemplateArgs struct {
	Nav   NavBar
	Title string
	Body  string
	Links LinkMetadata
}
