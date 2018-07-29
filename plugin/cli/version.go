package cli

import (
	"fmt"

	// external
	"github.com/spf13/cobra"

	// internal
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/config"
)

// VersionCmd shows the version
var VersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Display version information",
	Long:    fmt.Sprintf("Display version information for %s.", config.ProgramName),
	Example: fmt.Sprintf("  %s version", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		getOutput().Info(config.Version)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
