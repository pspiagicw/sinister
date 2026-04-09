package config

import (
	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
	"github.com/mitchellh/go-homedir"
	"github.com/pspiagicw/goreland"
)

type Config struct {
	VideoFolder string   `toml:"videoFolder"`
	URLS        []string `toml:"urls"`
	Quality     string   `toml:"quality"`
}
type Feed struct {
	URL  string   `toml:"url"`
	Tags []string `toml:"tags"`
}

func ParseConfig(configPath string) *Config {
	path := getConfigPath(configPath)

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
		goreland.LogFatal("Video folder not set in config file")
	}
	if len(conf.URLS) == 0 {
		goreland.LogFatal("No URLs set in config file")
	}
	if conf.Quality == "" {
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
func getConfigPath(configPath string) string {

	if configPath != "" {
		return configPath
	}

	path, err := xdg.SearchConfigFile("sinister/config.toml")
	if err != nil {
		goreland.LogFatal("Error while searching for config file: %v", err)
	}
	return path
}
