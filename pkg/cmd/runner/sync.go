package runner

import (
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

type CmdSyncOptions struct {
}

func NewCmdSync(runr *state.Runr) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: i18n.T("Sync a CI pipeline"),
		Long:  i18n.T("Sync a CI pipeline"),
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	return cmd
}

func (o *CmdSyncOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}
func (o *CmdSyncOptions) Validate() error {
	return nil
}
func (o *CmdSyncOptions) Run() error {
	return nil
}
