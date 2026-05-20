package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestParseConfig_Valid(t *testing.T) {
	path := writeTempConfig(t, `
videoFolder = "/tmp/videos"
quality = "hd720"
urls = ["https://example.com/feed"]
`)
	conf := ParseConfig(path)

	if conf.VideoFolder != "/tmp/videos" {
		t.Errorf("unexpected VideoFolder: %s", conf.VideoFolder)
	}
	if conf.Quality != "hd720" {
		t.Errorf("unexpected Quality: %s", conf.Quality)
	}
	if len(conf.URLS) != 1 || conf.URLS[0] != "https://example.com/feed" {
		t.Errorf("unexpected URLs: %v", conf.URLS)
	}
}

func TestParseConfig_HomeDirExpansion(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}

	path := writeTempConfig(t, `
videoFolder = "~/videos"
quality = "hd720"
urls = ["https://example.com/feed"]
`)
	conf := ParseConfig(path)
	expected := filepath.Join(home, "videos")
	if conf.VideoFolder != expected {
		t.Errorf("expected %s, got %s", expected, conf.VideoFolder)
	}
}

func TestResolveConfigPath_Explicit(t *testing.T) {
	path := "/some/explicit/path.toml"
	got := ResolveConfigPath(path)
	if got != path {
		t.Errorf("expected %s, got %s", path, got)
	}
}

func TestWriteAndReadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	original := &Config{
		VideoFolder: "/tmp/videos",
		Quality:     "hd1080",
		URLS:        []string{"https://a.com", "https://b.com"},
	}
	WriteConfig(path, original)

	roundtrip := ParseConfig(path)
	if roundtrip.VideoFolder != original.VideoFolder {
		t.Errorf("VideoFolder mismatch: got %s", roundtrip.VideoFolder)
	}
	if roundtrip.Quality != original.Quality {
		t.Errorf("Quality mismatch: got %s", roundtrip.Quality)
	}
	if len(roundtrip.URLS) != 2 {
		t.Errorf("expected 2 URLs, got %d", len(roundtrip.URLS))
	}
}

func TestReadConfigForUpdate_ReturnsPathAndConfig(t *testing.T) {
	path := writeTempConfig(t, `
videoFolder = "/tmp/videos"
quality = "hd720"
urls = ["https://x.com"]
`)
	conf, returnedPath := ReadConfigForUpdate(path)
	if returnedPath != path {
		t.Errorf("expected path %s, got %s", path, returnedPath)
	}
	if conf.Quality != "hd720" {
		t.Errorf("unexpected quality: %s", conf.Quality)
	}
}
