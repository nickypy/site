package server

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

var LINKS_PATH = "config/links.yml"

type LinkMetadata struct {
	Profiles []ProfileMetadata       `yaml:"profiles"`
	External []ExternalPostsMetadata `yaml:"external_posts"`
}

type ProfileMetadata struct {
	Title string `yaml:"title"`
	Link  string `yaml:"link"`
}

func (pm ProfileMetadata) SVGSource() string {
	path := fmt.Sprintf("svg/%s.svg", pm.Title)
	contents := readFile(path)
	return string(contents)
}

type ExternalPostsMetadata struct {
	Title string    `yaml:"title"`
	Link  string    `yaml:"link"`
	Date  time.Time `yaml:"date"`
}

func (epm ExternalPostsMetadata) FormatDate() string {
	return epm.Date.Format("2006 Jan")
}

func NewLinks() LinkMetadata {
	contents := readFile(LINKS_PATH)
	var l LinkMetadata
	err := yaml.Unmarshal(contents, &l)
	if err != nil {
		log.Printf("%s: may be invalid", LINKS_PATH)
		log.Fatalln(err)
	}

	return l
}
