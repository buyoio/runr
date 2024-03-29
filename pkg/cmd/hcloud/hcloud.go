package hcloud

import (
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewCmdHcloud(runr *state.Runr) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "hcloud",
		Short:                 i18n.T("Manage runners on Hetzner Cloud"),
		Long:                  i18n.T("Set up runners on Hetzner Cloud for CI/CD pipelines"),
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewCmdUp(runr))
	cmd.AddCommand(NewCmdLs(runr))
	cmd.AddCommand(NewCmdRm(runr))

	return cmd
}
