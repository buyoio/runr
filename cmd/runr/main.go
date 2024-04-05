package main

import (
	"os"

	"github.com/buyoio/runr/pkg/cli"
)

func main() {
	root := cli.NewRunrCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
