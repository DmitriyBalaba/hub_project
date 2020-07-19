package cmd

import (
	"os"
	"path"

	"bitbucket.org/jazzserve/webapi/cmd"
)

var (
	debugMode  = false
	configPath = path.Base(os.Args[0]) + ".yaml"
)

const (
	Serve            = "serve"
	ServeDescription = "Starts API server"

	Migrate            = "migrate"
	MigrateDescription = "Starts DB Migration"
)

var rootCommand = cmd.NewDefaultRoot(&debugMode, &configPath)

func Execute() {
	cmd.Execute(rootCommand)
}
