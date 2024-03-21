package database

import (
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/feed"
)

func QueryCreators() []string {
	db := openDB()

	rows, err := db.Query("SELECT DISTINCT author FROM entries WHERE watched = 0")
	if err != nil {
		goreland.LogFatal("Error while querying: %v", err)
	}

	creators := scanStrings(rows)

	rows.Close()
	defer db.Close()

	return creators
}
func QueryVideos(creator string) []string {
	db := openDB()

	rows, err := db.Query("SELECT title FROM entries WHERE author = ? AND watched = 0", creator)
	if err != nil {
		goreland.LogFatal("Error while querying: %v", err)
	}

	videos := scanStrings(rows)

	rows.Close()
	defer db.Close()

	return videos
}
func QueryEntry(creator, video string) *feed.Entry {
	db := openDB()

	entry := new(feed.Entry)

	row := db.QueryRow("SELECT * FROM entries WHERE author = ? AND title = ?", creator, video)

	var id int

	err := row.Scan(&id, &entry.Author.Name, &entry.Title, &entry.Published, &entry.Link.URL, &entry.Watched, &entry.Slug)

	if err != nil {
		goreland.LogFatal("Error while scanning: %v", err)
	}

	defer db.Close()

	return entry
}
func TotalCreators() int {
	db := openDB()

	rows := runQuery(db, "SELECT DISTINCT author FROM entries")

	total := 0

	for rows.Next() {
		total++
	}

	rows.Close()
	db.Close()

	return total
}
func TotalEntries() int {

	db := openDB()

	rows := runQuery(db, "SELECT * FROM entries")

	total := 0

	for rows.Next() {
		total++
	}

	rows.Close()
	db.Close()

	return total
}
func UnwatchedEntries() int {
	db := openDB()

	rows := runQuery(db, "SELECT * FROM entries WHERE watched = 0")

	total := 0

	for rows.Next() {
		total++
	}

	rows.Close()
	db.Close()

	return total
}
func QueryAll() []feed.Entry {
	db := openDB()

	rows := runQuery(db, "SELECT * FROM entries WHERE watched = 0")

	var entries []feed.Entry

	for rows.Next() {
		entry := new(feed.Entry)

		var id int

		err := rows.Scan(&id, &entry.Author.Name, &entry.Title, &entry.Published, &entry.Link.URL, &entry.Watched, &entry.Slug)

		if err != nil {
			goreland.LogFatal("Error while scanning: %v", err)
		}

		entries = append(entries, *entry)
	}

	rows.Close()
	db.Close()
	return entries
}
