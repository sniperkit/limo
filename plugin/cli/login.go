package cli

import (
	"context"
	"fmt"

	// external
	"github.com/spf13/cobra"

	// internal
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/config"
	"github.com/sniperkit/snk.golang.vcs-starred/pkg/service"

	_ "github.com/sniperkit/snk.golang.vcs-starred/plugin/service/vcs/remote"
)

// LoginCmd lets you log in
var LoginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Log in to a service",
	Long:    "Log in to the service specified by [--service] (default: github).",
	Example: fmt.Sprintf("  %s login", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Get the specified service and log in
		svc, err := getService()
		fatalOnError(err)

		token, err := svc.Login(ctx)
		fatalOnError(err)

		// Update configuration with token
		config, err := getConfiguration()
		fatalOnError(err)

		config.GetService(service.Name(svc)).Token = token
		fatalOnError(config.WriteConfig())
	},
}

func init() {
	RootCmd.AddCommand(LoginCmd)
}
