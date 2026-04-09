package database

import (
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/feed"
)

func QueryCreators() []string {
	db := openDB()
	defer closeDB(db)

	rows, err := db.Query("SELECT DISTINCT author FROM entries WHERE watched = 0")
	if err != nil {
		goreland.LogFatal("Error while querying: %v", err)
	}
	defer closeRows(rows)

	creators := scanStrings(rows)

	return creators
}
func QueryVideos(creator string) []string {
	db := openDB()
	defer closeDB(db)

	rows, err := db.Query("SELECT title FROM entries WHERE author = ? AND watched = 0", creator)
	if err != nil {
		goreland.LogFatal("Error while querying: %v", err)
	}
	defer closeRows(rows)

	videos := scanStrings(rows)

	return videos
}
func QueryEntry(creator, video string) *feed.Entry {
	db := openDB()
	defer closeDB(db)

	row := db.QueryRow("SELECT * FROM entries WHERE author = ? AND title = ?", creator, video)
	return scanEntry(row)
}
func TotalCreators() int {
	db := openDB()
	defer closeDB(db)

	return countQuery(db, "SELECT COUNT(DISTINCT author) FROM entries")
}
func TotalEntries() int {
	db := openDB()
	defer closeDB(db)
	return countQuery(db, "SELECT COUNT(*) FROM entries")
}
func CountUnwatched() int {
	db := openDB()
	defer closeDB(db)
	return countQuery(db, "SELECT COUNT(*) FROM entries WHERE watched = 0")
}

func CountEntriesByCreator(creator string) int {
	db := openDB()
	defer closeDB(db)
	return countQuery(db, "SELECT COUNT(*) FROM entries WHERE author = ?", creator)
}

func CountUnwatchedByCreator(creator string) int {
	db := openDB()
	defer closeDB(db)
	return countQuery(db, "SELECT COUNT(*) FROM entries WHERE author = ? AND watched = 0", creator)
}

func QueryUnwatched() []feed.Entry {
	db := openDB()
	defer closeDB(db)

	rows := runQuery(db, "SELECT * FROM entries WHERE watched = 0")
	defer closeRows(rows)

	var entries []feed.Entry

	for rows.Next() {
		entry := scanEntry(rows)

		entries = append(entries, *entry)
	}

	if err := rows.Err(); err != nil {
		goreland.LogFatal("Error while iterating rows: %v", err)
	}

	return entries
}
