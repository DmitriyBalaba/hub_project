package cmd

import (
	"hub_project/env"

	"bitbucket.org/jazzserve/webapi/cmd"
)

func init() {
	rootCommand.AddCommand(cmd.New(Serve, ServeDescription, serve))
}

func serve() error {
	server, err := env.NewBuilder(debugMode).
		ReadConfigFile(configPath).
		SetupPostgresDB().
		SetupSessionStore().
		SetupFileManager(env.TemplatesDir).
		SetupEmail().
		ConfigureCors().
		BuildWebServer()

	if err != nil {
		return err
	}

	server.Run()

	return nil
}
