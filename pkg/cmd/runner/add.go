package runner

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/goodies/templates"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
)

type CmdAddOptions struct {
	*state.Runr

	state.Runner
	Name   string
	Host   string
	HCloud string

	Force bool
}

func NewCmdAddOptions(runr *state.Runr) *CmdAddOptions {
	return &CmdAddOptions{
		Runr: runr,
	}
}

func NewCmdAdd(runr *state.Runr) *cobra.Command {
	o := NewCmdAddOptions(runr)
	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: i18n.T(`Add a new runner`),
		Long:  templates.LongDesc(i18n.T(`Add a new runner`)),
		Example: templates.Examples(i18n.T(`
			# Add a new runner
			runr runner add --host 172.17.0.1:22 --user root
		`)),
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

func (o *CmdAddOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Host, "host", "", i18n.T("Host and port"))
	cmd.Flags().StringVar(&o.HCloud, "hcloud", "", i18n.T("HCloud name"))
	cmd.MarkFlagsOneRequired("host", "hcloud")
	cmd.MarkFlagsMutuallyExclusive("host", "hcloud")

	cmd.Flags().StringVarP(&o.User, "user", "u", "root", i18n.T("Username"))
	cmd.Flags().StringVarP(&o.Auth, "auth", "a", "ssh-agent", i18n.T("Type of authentication"))
	cmd.Flags().BoolVar(&o.Force, "force", false, i18n.T("Overwrite existing runner, if exists"))
}

func (o *CmdAddOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.Host != "" {
		o.Runner.Host = &o.Host
	}
	if o.HCloud != "" {
		o.Runner.HCloud = &o.HCloud
	}
	if len(args) > 0 {
		o.Name = args[0]
	}
	return nil
}
func (o *CmdAddOptions) Validate(cmd *cobra.Command) error {
	if o.Name == "" {
		return cmdutil.UsageErrorf(cmd, i18n.T("NAME is required"))
	}
	if o.Runner.Host == nil && o.Runner.HCloud == nil {
		return cmdutil.UsageErrorf(cmd, i18n.T("Either --host or --hcloud must be specified"))
	}
	return nil
}
func (o *CmdAddOptions) Run() error {
	if _, ok := o.State(true).Runners[o.Name]; !ok || o.Force {
		o.State(true).Runners[o.Name] = &o.Runner
	} else {
		return fmt.Errorf(i18n.T("Runner %v already exists, overwrite with --force"), o.Name)
	}

	if err := o.Marshal(); err != nil {
		return err
	}
	return o.IO().Print(o.State().Runners[o.Name])
}
