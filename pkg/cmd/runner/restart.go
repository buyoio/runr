package runner

import (
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/spf13/cobra"
)

type CmdRestartOptions struct {
}

func NewCmdRestart(runr *state.Runr) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart a CI pipeline",
		Long:  "Restart a CI pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	return cmd
}

func (o *CmdRestartOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}
func (o *CmdRestartOptions) Validate() error {
	return nil
}
func (o *CmdRestartOptions) Run() error {
	return nil
}
