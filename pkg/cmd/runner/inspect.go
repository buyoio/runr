package runner

import (
	"context"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/goodies/ssh"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

type CmdInspectOptions struct {
	*state.Runr

	ctx  context.Context
	Name string
}

type CmdInspectOutput struct {
	System *ssh.RemoteInfo `json:"system" yaml:"system"`
}

func NewCmdInspectOptions(runr *state.Runr) *CmdInspectOptions {
	return &CmdInspectOptions{
		Runr: runr,
	}
}

func NewCmdInspect(runr *state.Runr) *cobra.Command {
	o := NewCmdInspectOptions(runr)
	cmd := &cobra.Command{
		Use:   "inspect NAME",
		Short: i18n.T("Inspect a runner"),
		Long:  i18n.T("Inspect a runner remotly"),
		Example: i18n.T(`
			# Inspect a runner remotly
			runr runner inspect my-runner
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

func (o *CmdInspectOptions) AddFlags(cmd *cobra.Command) {
}

func (o *CmdInspectOptions) Complete(cmd *cobra.Command, args []string) error {
	o.ctx = cmd.Context()
	if len(args) > 0 {
		o.Name = args[0]
	}
	return nil
}
func (o *CmdInspectOptions) Validate(cmd *cobra.Command) error {
	if o.Name == "" {
		return cmdutil.UsageErrorf(cmd, i18n.T("NAME is required"))
	}
	return nil
}
func (o *CmdInspectOptions) Run() error {
	client, err := getRunnerSSHClient(o.Runr, o.Name)
	if err != nil {
		return err
	}

	var out CmdInspectOutput
	out.System, err = client.RemoteInfo()
	if err != nil {
		return err
	}

	return o.Runr.IO().Print(out)
}
