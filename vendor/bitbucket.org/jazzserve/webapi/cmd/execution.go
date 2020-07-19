package cmd

import (
	"bitbucket.org/jazzserve/webapi/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

// Execute runs command configurations and exits on its finish
func Execute(command *cobra.Command) {
	if err := command.Execute(); err != nil {
		log.Fatal().Msg("Execution failed: " + err.Error())
		os.Exit(utils.ErrorExitCode)
	}
	os.Exit(utils.SuccessExitCode)
}
