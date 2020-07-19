package cmd

import (
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Default command configurations flags, shorts and descriptions
const (
	Debug      = "debug"
	DebugShort = "d"
	DebugUse   = "server debug mode"

	Config      = "config"
	ConfigShort = "c"
	ConfigUse   = "config file path"
)

// Creates default root command with debug and config flags
func NewDefaultRoot(debugMode *bool, configPath *string) *cobra.Command {
	binaryName := path.Base(os.Args[0])

	root := &cobra.Command{
		Use:   binaryName,
		Short: "By JazzServe",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// showing debug messages in debug mode only
			if *debugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
			log.Debug().Msgf("%s is in debug mode: you are able to see debug messages", binaryName)
		},
	}

	root.PersistentFlags().StringVarP(configPath, Config, ConfigShort, binaryName, ConfigUse)
	root.PersistentFlags().BoolVarP(debugMode, Debug, DebugShort, false, DebugUse)

	return root
}
