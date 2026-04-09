package database

import (
	"github.com/gosimple/slug"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/feed"
)

const insertEntrySQL = "INSERT OR IGNORE INTO entries(author, title, published, link, watched, slug) values(?,?,?,?,?,?)"

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

	result, err := db.Exec(insertEntrySQL, authr, title, published, link, 0, slug)

	if err != nil {
		goreland.LogFatal("Error while executing statement: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		goreland.LogFatal("Error while checking inserted rows: %v", err)
	}

	return rowsAffected > 0
}
