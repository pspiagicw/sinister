package manage

import (
	"fmt"

	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
)

func Status(configPath string) {

	conf := config.ParseConfig(configPath)

	fmt.Println()
	fmt.Println("SINISTER")
	fmt.Println("Video: ", conf.VideoFolder)
	fmt.Println()
	fmt.Println("URLs: ")

	for _, url := range conf.URLS {
		fmt.Printf("- %s\n", url)
	}

	fmt.Println()

	totalEntries := database.TotalEntries()
	totalCreators := database.TotalCreators()
	fmt.Println("Total Videos: ", totalEntries)
	fmt.Println("Total Creators: ", totalCreators)

	unwatchedEntries := database.CountUnwatched()

	fmt.Println("Unwatched Videos: ", unwatchedEntries)
}
func printFeed(feed string, value config.Feed) {
	fmt.Println("Feed: ", feed)
	fmt.Println("URL: ", value.URL)
	fmt.Println("Tags: ")
	for _, tag := range value.Tags {
		fmt.Printf("- %s\n", tag)
	}
	fmt.Println()
}
