package server

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
	blogRenderer := NewBlogRenderer(b.InputPath, b.OutputPath, b.BlogOptions...)
	blogRenderer.Render()
}
