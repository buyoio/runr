package runner

import (
	"github.com/buyoio/goodies/templates"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

func NewCmdRunner(runr *state.Runr) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "runner",
		Short:                 i18n.T("Set up static runners"),
		Long:                  i18n.T("Set up static runners for CI/CD pipelines"),
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	templates.AddGroup(cmd, "config", "Configure static runners",
		NewCmdGet(runr),
		NewCmdAdd(runr),
		NewCmdRemove(runr),
	)

	templates.AddGroup(cmd, "remote", "Operate on static runners",
		NewCmdUpgrade(runr),
		NewCmdInspect(runr),
		NewCmdStart(runr),
		NewCmdStop(runr),
		NewCmdRestart(runr),
		NewCmdSync(runr),
	)

	return cmd
}
