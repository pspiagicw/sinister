package database

import (
	"database/sql"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/feed"
)

const (
	createEntriesTableSQL = "CREATE TABLE IF NOT EXISTS entries (id INTEGER PRIMARY KEY, author TEXT, title TEXT UNIQUE, published TEXT, link TEXT, watched INTEGER, slug TEXT UNIQUE)"
)

type scanner interface {
	Scan(dest ...interface{}) error
}

type queryExecer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func UpdateWatched(entry *feed.Entry) {
	db := openDB()
	defer closeDB(db)

	_, err := db.Exec("UPDATE entries SET watched = 1 WHERE slug = ?", entry.Slug)

	if err != nil {
		goreland.LogFatal("Error while updating watched status: %v", err)
	}
}

func MarkAllUnwatched() int {
	db := openDB()
	defer closeDB(db)

	result, err := db.Exec("UPDATE entries SET watched = 0")
	if err != nil {
		goreland.LogFatal("Error while resetting watched status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		goreland.LogFatal("Error while checking updated rows: %v", err)
	}

	return int(rowsAffected)
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
	_, err := db.Exec(createEntriesTableSQL)
	if err != nil {
		goreland.LogFatal("Error while creating table: %v", err)
	}
}
func scanStrings(rows *sql.Rows) []string {
	var elements []string

	for rows.Next() {
		var element string
		err := rows.Scan(&element)
		if err != nil {
			goreland.LogFatal("Error while scanning: %v", err)
		}
		elements = append(elements, element)
	}

	if err := rows.Err(); err != nil {
		goreland.LogFatal("Error while iterating rows: %v", err)
	}

	return elements
}
func runQuery(db *sql.DB, query string, args ...interface{}) *sql.Rows {

	rows, err := db.Query(query, args...)
	if err != nil {
		goreland.LogFatal("Error while querying: %v", err)
	}

	return rows

}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		goreland.LogError("Error while closing database: %v", err)
	}
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		goreland.LogError("Error while closing rows: %v", err)
	}
}

func scanEntry(s scanner) *feed.Entry {
	entry := new(feed.Entry)
	var id int

	err := s.Scan(&id, &entry.Author.Name, &entry.Title, &entry.Published, &entry.Link.URL, &entry.Watched, &entry.Slug)
	if err != nil {
		goreland.LogFatal("Error while scanning: %v", err)
	}

	return entry
}

func countQuery(db *sql.DB, query string, args ...interface{}) int {
	var total int

	err := db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		goreland.LogFatal("Error while counting query results: %v", err)
	}

	return total
}
