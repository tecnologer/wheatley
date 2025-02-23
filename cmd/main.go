package main

import (
	"os"

	"github.com/tecnologer/wheatley/cmd/cli"
	"github.com/tecnologer/wheatley/pkg/utils/log"
)

var version string

func main() {
	newCLI := cli.NewCLI(version)

	if err := newCLI.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}
