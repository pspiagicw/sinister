package update

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/config"
	"github.com/pspiagicw/sinister/pkg/database"
)

func Update(args []string) {
	config := config.ParseConfig()
	updateDatabase(config)
}
func updateDatabase(conf *config.Config) {
	for _, url := range conf.URLS {
		getLatest(url)
	}
}
func getLatest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}
	contents, err := io.ReadAll(resp.Body)
	var feed database.Feed
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}
	err = xml.Unmarshal(contents, &feed)
	fmt.Println(feed)
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}
}
