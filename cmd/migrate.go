package cmd

import (
	"hub_project/env"
	"hub_project/models"

	"bitbucket.org/jazzserve/webapi/cmd"
	"github.com/rs/zerolog/log"
)

func init() {
	rootCommand.AddCommand(cmd.New(Migrate, MigrateDescription, migrate))
}

func migrate() error {
	environment, err := env.NewBuilder(debugMode).
		ReadConfigFile(configPath).
		SetupFileManager(env.MigrationDir).
		SetupPostgresDB().
		Build()

	if err != nil {
		return err
	}

	migrations, err := environment.ReadMigrationScripts()
	if err != nil {
		return err
	}

	err = environment.Storage().Migrate(models.GetModels(), migrations)
	if err != nil {
		return err
	}

	log.Info().Msg("Migrated successfully!")
	return nil
}
