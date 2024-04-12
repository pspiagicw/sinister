package resolver

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pspiagicw/fs"
	"github.com/pspiagicw/goreland"
)

func HomeDir() string {
	return xdg.Home
}

func DataDir() string {
	location := filepath.Join(xdg.DataHome, "sinister")
	fs.EnsurePathExists(location)
	return location
}

func DatabasePath() string {

	d := DataDir()
	d = filepath.Join(d, "db")

	goreland.LogInfo("Using %s for database", d)
	return d

}
