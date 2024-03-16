package database

import (
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/feed"
)

func GetCreators() []string {
	db := openDB()

	rows, err := db.Query("SELECT DISTINCT author FROM entries WHERE watched = 0")
	if err != nil {
		goreland.LogFatal("Error while querying: %v", err)
	}

	var creators []string

	for rows.Next() {
		var creator string
		err = rows.Scan(&creator)
		if err != nil {
			goreland.LogFatal("Error while scanning: %v", err)
		}
		creators = append(creators, creator)
	}

	rows.Close()
	defer db.Close()

	return creators
}
func GetVideos(creator string) []string {
	db := openDB()

	rows, err := db.Query("SELECT title FROM entries WHERE author = ? AND watched = 0", creator)
	if err != nil {
		goreland.LogFatal("Error while querying: %v", err)
	}

	var videos []string

	for rows.Next() {
		var video string
		err = rows.Scan(&video)
		if err != nil {
			goreland.LogFatal("Error while scanning: %v", err)
		}
		videos = append(videos, video)
	}

	rows.Close()
	defer db.Close()

	return videos
}
func GetEntry(creator, video string) *feed.Entry {
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
