package database

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
}
type Link struct {
	URL string `xml:"href,attr"`
}
