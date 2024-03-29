package runner

import (
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

type CmdStopOptions struct {
}

func NewCmdStop(runr *state.Runr) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: i18n.T("Stop a CI pipeline"),
		Long:  i18n.T("Stop a CI pipeline"),
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	return cmd
}

func (o *CmdStopOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}
func (o *CmdStopOptions) Validate() error {
	return nil
}
func (o *CmdStopOptions) Run() error {
	return nil
}
