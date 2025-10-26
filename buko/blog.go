package server

import (
	"sync"
)

const BLOG_BASE_PATH = "blog"

type BlogRenderCache struct {
	Prefix                string
	OutputPath            string
	Items                 []BlogPost
	ShouldListUnpublished bool
	Links                 LinkMetadata
	filename              string
	markdown              MarkdownRenderer
	template              *TemplateRenderer
	mutex                 *sync.RWMutex
	navBar                NavBar
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
