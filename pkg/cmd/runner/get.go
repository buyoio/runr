package runner

import (
	"fmt"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

type CmdGetOptions struct {
	*state.Runr

	All   bool
	Name  string
	Force bool
}

func NewCmdGetOptions(runr *state.Runr) *CmdGetOptions {
	return &CmdGetOptions{
		Runr: runr,
	}
}

func NewCmdGet(runr *state.Runr) *cobra.Command {
	o := NewCmdGetOptions(runr)
	cmd := &cobra.Command{
		Use:   "get NAME | -all",
		Short: i18n.T("Get a runner"),
		Long:  i18n.T("Get a runner from the configuration"),
		Example: i18n.T(`
			# Get a runner from the configuration
			runr runner get my-runner
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

func (o *CmdGetOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", false, i18n.T("Get all runners"))
	cmd.Flags().BoolVarP(&o.Force, "force", "f", false, i18n.T("Throws no error if the runner does not exist"))
}

func (o *CmdGetOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		o.Name = args[0]
	}
	return nil
}
func (o *CmdGetOptions) Validate(cmd *cobra.Command) error {
	if o.Name == "" && !o.All {
		return cmdutil.UsageErrorf(cmd, i18n.T("NAME is required"))
	} else if o.Name != "" && o.All {
		return cmdutil.UsageErrorf(cmd, i18n.T("NAME and --all are mutually exclusive"))
	}
	return nil
}

func (o *CmdGetOptions) Run() error {
	if o.All {
		return o.Runr.IO().Print(o.State().Runners)
	}

	r, ok := o.State().Runners[o.Name]
	if !ok {
		if o.Force {
			return nil
		}
		return cmdutil.UsageErrorf(nil, fmt.Sprintf(i18n.T("Runner %s not found"), o.Name))
	}

	return o.Runr.IO().Print(r)
}
