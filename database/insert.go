package database

import (
	neturl "net/url"
	"strings"

	"github.com/gosimple/slug"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/feed"
)

const insertEntrySQL = "INSERT OR IGNORE INTO entries(author, title, published, link, watched, slug, video_id) values(?,?,?,?,?,?,?)"

func Add(entry *feed.Entry) bool {
	db := openDB()
	defer closeDB(db)
	return insertEntry(db, entry)
}

func insertEntry(db queryExecer, e *feed.Entry) bool {
	authr := e.Author.Name
	title := e.Title
	published := e.Published
	link := e.Link.URL
	slug := slug.Make(title)
	videoID := ExtractVideoID(link)

	result, err := db.Exec(insertEntrySQL, authr, title, published, link, 0, slug, nilIfEmpty(videoID))
	if err != nil {
		goreland.LogFatal("Error while executing statement: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		goreland.LogFatal("Error while checking inserted rows: %v", err)
	}

	return rowsAffected > 0
}

func ExtractVideoID(rawURL string) string {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return ""
	}
	if id := u.Query().Get("v"); id != "" {
		return id
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) == 2 && parts[0] == "shorts" {
		return parts[1]
	}
	if len(parts) == 1 && (u.Host == "youtu.be" || u.Host == "www.youtu.be") {
		return parts[0]
	}
	return ""
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
