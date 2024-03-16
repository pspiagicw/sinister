package config

import (
	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/help"
)

type Config struct {
	VideoFolder string   `toml:"video_folder"`
	URLS        []string `toml:"urls"`
}

func ParseConfig() *Config {
	path := getConfigPath()

	conf := readConfig(path)

	checkConfig(conf)

	return conf
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
}

func readConfig(path string) *Config {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		goreland.LogFatal("Error while reading config file: %v", err)
	}
	return &conf
}
func getConfigPath() string {
	path, err := xdg.SearchConfigFile("sinister/config.toml")
	if err != nil {
		help.HelpConfig()
		goreland.LogFatal("Error while searching for config file: %v", err)
	}
	return path
}
