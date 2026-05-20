package manage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/pspiagicw/sinister/config"
)

func writeCleanConfig(t *testing.T, urls ...string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	urlList := ""
	for _, u := range urls {
		urlList += fmt.Sprintf(`"%s",`, u)
	}

	content := fmt.Sprintf(`videoFolder = "/tmp/videos"
quality = "hd720"
urls = [%s]
`, urlList)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	return path
}

func TestClean_Removes404URLs(t *testing.T) {
	deadSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer deadSrv.Close()

	liveSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer liveSrv.Close()

	configPath := writeCleanConfig(t, deadSrv.URL, liveSrv.URL)

	Clean(CleanOptions{ConfigPath: configPath})

	// Re-read config and verify only live URL remains
	conf, _ := readUpdatedConfig(t, configPath)
	if len(conf) != 1 || conf[0] != liveSrv.URL {
		t.Errorf("expected only live URL to remain, got %v", conf)
	}
}

func TestClean_KeepsLiveURLs(t *testing.T) {
	liveSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer liveSrv.Close()

	configPath := writeCleanConfig(t, liveSrv.URL)

	Clean(CleanOptions{ConfigPath: configPath})

	conf, _ := readUpdatedConfig(t, configPath)
	if len(conf) != 1 {
		t.Errorf("expected 1 URL kept, got %v", conf)
	}
}

func TestClean_DryRun(t *testing.T) {
	deadSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer deadSrv.Close()

	configPath := writeCleanConfig(t, deadSrv.URL)

	Clean(CleanOptions{ConfigPath: configPath, DryRun: true})

	conf, _ := readUpdatedConfig(t, configPath)
	if len(conf) != 1 || conf[0] != deadSrv.URL {
		t.Errorf("dry run should not remove URLs, got %v", conf)
	}
}

func TestClean_KeepsErrorURLs(t *testing.T) {
	// A URL that fails to connect should be kept (error ≠ 404)
	configPath := writeCleanConfig(t, "http://127.0.0.1:1")

	Clean(CleanOptions{ConfigPath: configPath})

	conf, _ := readUpdatedConfig(t, configPath)
	if len(conf) != 1 {
		t.Errorf("connection error URLs should be kept, got %v", conf)
	}
}

// readUpdatedConfig reads URLs from a saved TOML config using the config package.
func readUpdatedConfig(t *testing.T, path string) ([]string, string) {
	t.Helper()
	conf, confPath := config.ReadConfigForUpdate(path)
	return conf.URLS, confPath
}
