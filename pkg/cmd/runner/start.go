package runner

import (
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

type CmdStartOptions struct {
}

func NewCmdStart(runr *state.Runr) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: i18n.T("Start a CI pipeline"),
		Long:  i18n.T("Start a CI pipeline"),
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	return cmd
}

func (o *CmdStartOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}
func (o *CmdStartOptions) Validate() error {
	return nil
}
func (o *CmdStartOptions) Run() error {
	return nil
}
