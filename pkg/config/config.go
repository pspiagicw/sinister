package config

type Config struct {
	URLS []string
}

func ParseConfig() *Config {
	return &Config{
		URLS: []string{
			"https://www.youtube.com/feeds/videos.xml?channel_id=UCeeFfhMcJa1kjtfZAGskOCA",
		},
	}
}
