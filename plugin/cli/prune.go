package cli

import (
	"fmt"

	// external
	"github.com/spf13/cobra"

	// internal
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/config"
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/model"
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/service"

	_ "github.com/sniperkit/snk.golang.vcs-starred/plugin/service/vcs/remote"
)

var delete = false

// PruneCmd prunes local stars that are no longer starred on a remote service
var PruneCmd = &cobra.Command{
	Use:     "prune",
	Short:   "Prune unstarred repositories",
	Long:    "Prune from your local database any repositories you no longer have starred on [--service] (default: github).",
	Example: fmt.Sprintf("  %s prune", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		db, err := getDatabase()
		fatalOnError(err)

		svc, err := getService()
		fatalOnError(err)

		serviceName := service.Name(svc)
		dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
		fatalOnError(err)

		prunable, err := model.FindPrunableStars(db, dbSvc)
		fatalOnError(err)

		for _, star := range prunable {
			output.StarLine(&star)
			if delete {
				fatalOnError(star.Delete(db))
			}
		}
	},
}

func init() {
	PruneCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Actually delete from your local database (default: display what would be deleted)")
	RootCmd.AddCommand(PruneCmd)
}
