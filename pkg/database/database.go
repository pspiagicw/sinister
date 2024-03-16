package database

import (
	"database/sql"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/feed"
)

func UpdateWatched(entry *feed.Entry) {
	db := openDB()

	_, err := db.Exec("UPDATE entries SET watched = 1 WHERE slug = ?", entry.Slug)

	if err != nil {
		goreland.LogFatal("Error while updating watched status: %v", err)
	}
}

func openDB() *sql.DB {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}

	ensureTableExists(db)

	return db
}
func getDBPath() string {
	path, err := xdg.DataFile("sinister/sinister.db")

	if err != nil {
		goreland.LogFatal("Error while getting config path: %v", err)
	}
	return path
}

func ensureTableExists(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS entries (id INTEGER PRIMARY KEY, author TEXT, title TEXT UNIQUE, published TEXT, link TEXT, watched INTEGER, SLUG TEXT UNIQUE)")
	if err != nil {
		goreland.LogFatal("Error while creating table: %v", err)
	}
}
