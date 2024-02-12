package server

import (
	"crypto/sha256"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/gorilla/feeds"
)

func (b *BlogRenderCache) generateAtomFeed() {
	now := time.Now()
	author := &feeds.Author{Name: "nickypy", Email: ""}
	baseLink := "https://nickypy.com"

	feed := &feeds.Feed{
		Title:       "nickypy",
		Link:        &feeds.Link{Href: baseLink},
		Description: "nickypy's blog",
		Author:      author,
		Created:     now,
	}

	var items []*feeds.Item

	for _, item := range b.Items {
		title := item.Metadata.Title
		hash := sha256.Sum256([]byte(item.Metadata.Title))

		if !item.Metadata.IsPublished {
			continue
		}

		items = append(items, &feeds.Item{
			Content:     item.Body,
			Created:     item.Metadata.Date,
			Description: title,
			Id:          fmt.Sprintf("%x", hash),
			Link:        &feeds.Link{Href: baseLink + item.Slug},
			Title:       title,
		})
	}
	feed.Items = items

	output, err := feed.ToAtom()
	if err != nil {
		log.Fatalln("Failed to generate atom feed %w", err)
	}

	writeFile(path.Join(b.OutputPath, "feed/atom.xml"), []byte(output))
}
