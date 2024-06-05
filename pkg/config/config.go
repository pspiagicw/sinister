package config

import (
	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
	"github.com/mitchellh/go-homedir"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/help"
)

type Config struct {
	VideoFolder string          `toml:"videoFolder"`
	URLS        []string        `toml:"urls"`
	Feeds       map[string]Feed `toml:"feed"`
	Quality     string          `toml:"quality"`
}
type Feed struct {
	URL  string   `toml:"url"`
	Tags []string `toml:"tags"`
}

func (c Config) GetURLs() []string {
	urls := make([]string, 0)
	for _, url := range c.URLS {
		urls = append(urls, url)
	}

	for _, feed := range c.Feeds {
		urls = append(urls, feed.URL)
	}

	return urls
}

func ParseConfig(opts *argparse.Opts) *Config {
	path := getConfigPath(opts)

	conf := readConfig(path)

	sanitizeConfig(conf)

	checkConfig(conf)

	return conf
}
func sanitizeConfig(conf *Config) {
	path, err := homedir.Expand(conf.VideoFolder)
	if err != nil {
		goreland.LogFatal("Error while expanding video folder path: %v", err)
	}
	conf.VideoFolder = path
}
func checkConfig(conf *Config) {
	if conf.VideoFolder == "" {
		help.HelpConfig()
		goreland.LogFatal("Video folder not set in config file")
	}
	if len(conf.URLS) == 0 {
		help.HelpConfig()
		goreland.LogFatal("No URLs set in config file")
	}
	if conf.Quality == "" {
		help.HelpConfig()
		goreland.LogFatal("Quality not set in config file")
	}
}

func readConfig(path string) *Config {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		goreland.LogFatal("Error while reading config file: %v", err)
	}
	return &conf
}
func getConfigPath(opts *argparse.Opts) string {
	if opts.Config != "" {
		return opts.Config
	}
	path, err := xdg.SearchConfigFile("sinister/config.toml")
	if err != nil {
		help.HelpConfig()
		goreland.LogFatal("Error while searching for config file: %v", err)
	}
	return path
}
