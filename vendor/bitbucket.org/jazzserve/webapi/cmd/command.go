package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type ExecutableCommand func() error

func New(key, description string, f ExecutableCommand) *cobra.Command {
	return &cobra.Command{
		Use:   key,
		Short: description,
		Run: func(cmd *cobra.Command, args []string) {
			if err := f(); err != nil {
				log.Fatal().Msgf("%s command failed: %s", key, err.Error())
				return
			}
		},
	}
}
