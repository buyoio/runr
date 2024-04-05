package runner

import (
	"fmt"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

type CmdRemoveOptions struct {
	*state.Runr

	name    string
	force   bool
	onlySCM bool
}

func NewCmdRemoveOptions(runr *state.Runr) *CmdRemoveOptions {
	return &CmdRemoveOptions{
		Runr: runr,
	}
}

func NewCmdRemove(runr *state.Runr) *cobra.Command {
	o := NewCmdRemoveOptions(runr)
	cmd := &cobra.Command{
		Use:     "remove NAME",
		Aliases: []string{"rm", "delete"},
		Short:   i18n.T("Remove a runner"),
		Long:    i18n.T("Remove a runner from the configuration and scm"),
		Example: i18n.T(`
			# Remove a runner from the configuration
			runr runner remove my-runner
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate(cmd))
			cmdutil.CheckErr(o.Run())
		},
	}
	o.AddFlags(cmd)
	output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML(), output.OptionFormat())

	return cmd
}

func (o *CmdRemoveOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&o.force, "force", "f", false, i18n.T("Force remove"))
	cmd.Flags().BoolVarP(&o.onlySCM, "only-scm", "r", false, i18n.T("Only remove from scm"))
}

func (o *CmdRemoveOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		o.name = args[0]
	}
	return nil
}
func (o *CmdRemoveOptions) Validate(cmd *cobra.Command) error {
	if o.name == "" {
		return cmdutil.UsageErrorf(cmd, i18n.T("NAME is required"))
	}
	return nil
}
func (o *CmdRemoveOptions) Run() error {
	var runner *state.Runner
	var ok bool
	if runner, ok = o.State().Runners[o.name]; !ok {
		if o.force {
			return nil
		}
		return fmt.Errorf(i18n.T("No runner found with name %s"), o.name)
	}

	provider, err := runner.Setup.SCMPlatform.NewProvider(o.GetContext(), o.Logger())
	if err != nil {
		return err
	}
	if err := provider.RemoveRunner(o.name); err != nil {
		return err
	}

	if !o.onlySCM {
		delete(o.State(true).Runners, o.name)
		if err := o.Runr.Marshal(); err != nil {
			return err
		}
	}

	return o.Runr.IO().Print(runner)
}
