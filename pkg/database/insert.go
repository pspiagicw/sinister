package database

import (
	"database/sql"

	"github.com/gosimple/slug"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/feed"
)

func Add(entry *feed.Entry) {
	db := openDB()

	insertEntry(db, entry)

	defer db.Close()
}

func insertEntry(db *sql.DB, e *feed.Entry) {

	authr := e.Author.Name
	title := e.Title
	published := e.Published
	link := e.Link.URL
	slug := slug.Make(title)

	stmt := prepareInsertStatement(db)
	executeInsertStatement(stmt, authr, title, published, link, slug)
	stmt.Close()
}
func executeInsertStatement(stmt *sql.Stmt, authr, title, published, link, slug string) {
	_, err := stmt.Exec(authr, title, published, link, 0, slug)

	if err != nil {
		goreland.LogFatal("Error while executing statement: %v", err)
	}
}
func prepareInsertStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO entries(author, title, published, link, watched, slug) values(?,?,?,?,?,?)")
	if err != nil {
		goreland.LogFatal("Error while generating statement: %v", err)
	}
	return stmt
}
