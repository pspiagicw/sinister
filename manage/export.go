package manage

import (
	"os"
	"strings"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/database"
)

const exportFile = "urls.txt"

func Export() {
	entries := database.QueryUnwatched()

	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, entry.Link.URL)
	}

	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}

	if err := os.WriteFile(exportFile, []byte(content), 0644); err != nil {
		goreland.LogFatal("Error while writing export file: %v", err)
	}

	for _, entry := range entries {
		entryCopy := entry
		database.UpdateWatched(&entryCopy)
	}

	goreland.LogSuccess("Exported %d URLs to %s", len(entries), exportFile)
}
