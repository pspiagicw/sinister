package feed

import (
	"time"

	"github.com/pspiagicw/goreland"
)

type Feed struct {
	Author  Author  `xml:"author"`
	Entries []Entry `xml:"entry"`
}

type Author struct {
	Name string `xml:"name"`
}

type Entry struct {
	Author    Author `xml:"author"`
	Title     string `xml:"title"`
	Published string `xml:"published"`
	Link      Link   `xml:"link"`
	Slug      string
	Watched   int
}
type Link struct {
	URL string `xml:"href,attr"`
}

func (e Entry) Date() time.Time {
	layout := time.RFC3339

	t, err := time.Parse(layout, e.Published)

	if err != nil {
		goreland.LogFatal("Error while parsing date: %v", err)
	}

	return t
}
