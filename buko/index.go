package server

import (
	"path"
)

type SiteBuilder struct {
	InputPath   string
	OutputPath  string
	BlogOptions []BlogOption
}

func NewSiteBuilder(inputPath string, outputPath string) *SiteBuilder {
	return &SiteBuilder{
		InputPath:   inputPath,
		OutputPath:  outputPath,
		BlogOptions: make([]BlogOption, 0),
	}
}

func (b *SiteBuilder) WithBlogOptions(opts []BlogOption) *SiteBuilder {
	b.BlogOptions = opts
	return b
}

func (b SiteBuilder) Build() {
	blogRenderer := NewBlogRenderer(b.InputPath, b.OutputPath, "blog.html", b.BlogOptions...)
	blogRenderer.Render()

	source := path.Join(b.OutputPath, blogRenderer.filename)
	destination := path.Join(b.OutputPath, "index.html")
	CopyFile(source, destination)
}
