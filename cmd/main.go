package main

import (
	"os"

	"github.com/tecnologer/wheatley/cmd/cli"
	"github.com/tecnologer/wheatley/pkg/constants/envvarname"
	"github.com/tecnologer/wheatley/pkg/utils/log"
)

var version string

func main() {
	newCLI := cli.NewCLI(version)

	err := os.Setenv(envvarname.BotVersion, version)
	if err != nil {
		log.Errorf("failed to set %s environment variable: %v", envvarname.BotVersion, err)
	}

	if err := newCLI.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}
