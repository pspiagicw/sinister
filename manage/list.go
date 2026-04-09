package manage

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/database"
)

func List() {
	stats := database.QueryCreatorStats()
	if len(stats) == 0 {
		goreland.LogInfo("No channels found in the database")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CHANNEL\tTOTAL\tUNWATCHED")
	for _, stat := range stats {
		fmt.Fprintf(w, "%s\t%d\t%d\n", stat.Name, stat.Total, stat.Unwatched)
	}
	w.Flush()
}
