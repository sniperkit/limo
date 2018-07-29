package cli

import (
	"fmt"

	// external
	"github.com/spf13/cobra"

	// internal
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/config"
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/model"
)

// DeleteCmd renames a tag
var DeleteCmd = &cobra.Command{
	Use:     "delete <tag>",
	Aliases: []string{"rm"},
	Short:   "Delete a tag",
	Long:    "Delete the tag named <tag>.",
	Example: fmt.Sprintf("  %s delete frameworks", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a tag")
		}

		db, err := getDatabase()
		fatalOnError(err)

		tag, err := model.FindTagByName(db, args[0])
		fatalOnError(err)

		if tag == nil {
			output.Fatal(fmt.Sprintf("Tag '%s' not found", args[0]))
		}

		fatalOnError(tag.Delete(db))

		output.Info(fmt.Sprintf("Deleted tag '%s'", tag.Name))
	},
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
}
