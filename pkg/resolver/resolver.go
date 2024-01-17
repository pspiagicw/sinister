package resolver

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pspiagicw/goreland"
)

func HomeDir() string {
	return xdg.Home
}

func DataDir() string {
	location := filepath.Join(xdg.DataHome, "sinister")
	ensureExists(location)
	return location
}

func ensureExists(location string) {
	if !dirExists(location) {
		err := os.MkdirAll(location, 0755)
		if err != nil {
			goreland.LogFatal("Error creating directory: %s, %v", location, err)
		}
	}
}

func dirExists(dir string) bool {
	_, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	} else if err != nil {
		goreland.LogFatal("Error stating directory: %v", err)
	}
	return true
}

func DatabasePath() string {

	d := DataDir()
	d = filepath.Join(d, "db")

	goreland.LogInfo("Using %s for database", d)
	return d

}
