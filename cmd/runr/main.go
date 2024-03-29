package main

import (
	"os"

	"github.com/buyoio/runr/pkg/cli"
	v "github.com/buyoio/runr/pkg/cmd/version"
)

// Magic variables set by goreleaser
var (
	version   string
	date      string
	commit    string
	treestate string
)

func main() {
	v.SetBuildInfo(v.Info{
		Version:      version,
		BuildDate:    date,
		GitCommit:    commit,
		GitTreeState: treestate,
	})
	root := cli.NewRunrCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
