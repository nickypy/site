package server

import (
	"cmp"
	"path"
	"slices"

	"gopkg.in/yaml.v2"
)

type SiteBuilder struct {
	InputPath       string
	OutputPath      string
	ListUnpublished bool
}

func NewSiteBuilder(inputPath string, outputPath string) *SiteBuilder {
	return &SiteBuilder{
		InputPath:       inputPath,
		OutputPath:      outputPath,
		ListUnpublished: false,
	}
}

type RenderedMarkdown struct {
	Entry   MarkdownFileEntry
	Meta    map[string]any
	Content string
}

func markdownToHTML(entries []MarkdownFileEntry) []RenderedMarkdown {
	results := make([]RenderedMarkdown, len(entries))
	mdRenderer := NewMarkdownRenderer()

	for i, e := range entries {

		meta, content := mdRenderer.RenderBytes(e.Content)
		results[i] = RenderedMarkdown{
			e, meta, content,
		}

	}

	return results
}

func generateNavBar(entries []MarkdownFileEntry) NavBar {
	deduped := make(map[string]bool)

	for _, entry := range entries {
		if entry.Prefix == "" {
			continue
		}

		key := entry.Prefix
		if _, exists := deduped[key]; !exists {
			deduped[key] = true
		}
	}

	items := make([]NavBarItem, 0)

	for key := range deduped {
		items = append(items, NavBarItem{
			key, "/" + key + ".html",
		})
	}

	return NavBar{items}
}

func (b SiteBuilder) Build() {
	mdFiles := GetAllMarkdownFiles()
	navBar := generateNavBar(mdFiles)
	renderedHTML := markdownToHTML(mdFiles)
	links := NewLinks()

	tmpl := NewTemplateRenderer().WithDefault()

	mainItems := make([]RenderedMarkdown, 0)
	blogItems := make([]BlogPost, 0)

	for _, entry := range renderedHTML {
		if entry.Entry.Prefix == "posts" {
			var metadata BlogMetadata
			out, err := yaml.Marshal(entry.Meta)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(out, &metadata)
			if err != nil {
				panic(err)
			}

			post := BlogPost{
				metadata.Slug,
				metadata.Slug + ".html",
				entry.Content,
				metadata,
			}
			blogItems = append(blogItems, post)
		} else {
			mainItems = append(mainItems, entry)
		}
	}

	slices.SortFunc(blogItems, func(a, b BlogPost) int {
		return cmp.Compare(b.Metadata.Date.Unix(), a.Metadata.Date.Unix())
	})

	for _, item := range mainItems {
		title := item.Entry.Key
		body := item.Content

		res := tmpl.RenderPage(PageTemplateArgs{
			navBar,
			title,
			body,
			links,
		})

		filename := path.Join(b.OutputPath, title+".html")
		writeFile(filename, res)
	}

	for _, item := range blogItems {
		res := tmpl.RenderBlogPost(PageTemplateArgs{
			navBar,
			item.Metadata.Title,
			item.Body,
			links,
		})
		filename := path.Join(b.OutputPath, item.Path)
		writeFile(filename, res)
	}

	publishablePosts := make([]BlogPost, 0)
	for _, post := range blogItems {
		if post.Metadata.IsPublished || b.ListUnpublished {
			publishablePosts = append(publishablePosts, post)
		}
	}

	res := tmpl.RenderBlogPage(BlogPageTemplateArgs{
		navBar,
		"nickypy: blog",
		"",
		publishablePosts,
		links,
	})

	writeFile(path.Join(b.OutputPath, "posts.html"), res)

	err := CopyDirectory(
		path.Join(b.InputPath, "assets"),
		path.Join(b.OutputPath),
	)
	if err != nil {
		panic(err)
	}

	GenerateAtomFeed(blogItems, b.OutputPath)
}
