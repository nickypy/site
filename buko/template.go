package server

import (
	"io"
	p "path"
	"text/template"
)

type TemplateRenderer struct {
	baseTmpl string
	tmpl     map[string]*template.Template
}

func NewTemplateRenderer(base string) *TemplateRenderer {
	return &TemplateRenderer{base, make(map[string]*template.Template)}
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

type BlogTemplateArgs struct {
	Title     string
	Body      string
	BlogItems []BlogPost
	Links     LinkMetadata
}

type BlogPostTemplateArgs struct {
	Title string
	Body  string
	Links LinkMetadata
}
