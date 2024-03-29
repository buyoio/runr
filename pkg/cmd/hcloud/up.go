package hcloud

import (
	"fmt"
	"os/exec"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/runr/pkg/b/hcloud"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"

	"github.com/spf13/cobra"
)

type CmdUpOptions struct {
	*state.Runr
	hcloud *hcloud.Hcloud

	Name             string
	EnvironmentRef   string
	Labels           []string
	serverCreateOpts *state.HCloudServerCreateOpts
}

func NewCmdUp(runr *state.Runr) *cobra.Command {
	o := &CmdUpOptions{
		Runr: runr,
	}
	cmd := &cobra.Command{
		Use:   "up NAME --env ENVIRONMENT [flags]",
		Short: i18n.T("Provision a new hcloud runner"),
		Long:  i18n.T("Provision a new hcloud runner"),
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

func (o *CmdUpOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.EnvironmentRef, "env", "e", "", "hcloud environment to use")
	cmd.Flags().StringArrayVarP(&o.Labels, "label", "l", nil, "Runner labels to append")
}

func (o *CmdUpOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error
	if o.hcloud, err = hcloud.BinaryCmd(o.Runr, cmd); err != nil {
		return err
	}

	if len(args) > 0 {
		o.Name = args[0]
	}

	sco, ok := o.State().Environments.HCloud[o.EnvironmentRef]
	if !ok {
		return cmdutil.UsageErrorf(cmd, "hcloud environment %s not found", o.EnvironmentRef)
	}
	sco.Name = o.Name

	o.Labels = append(o.Labels, "runr="+o.EnvironmentRef)
	sco.UserData.SetLabels(o.Labels)
	sco.Labels = o.Labels
	o.serverCreateOpts = sco

	return nil
}
func (o *CmdUpOptions) Validate(cmd *cobra.Command) error {
	if err := o.hcloud.HasValidToken(); err != nil {
		return err
	}
	if o.Name == "" {
		return cmdutil.UsageErrorf(cmd, "NAME is required")
	}
	if o.EnvironmentRef == "" {
		return cmdutil.UsageErrorf(cmd, "HCloud environment is required (--env)")
	}

	return nil
}
func (o *CmdUpOptions) Run() error {
	cmd, err := o.hcloud.ServerCreate(o.serverCreateOpts)
	if err != nil {
		return fmt.Errorf("%s", err.(*exec.ExitError).Stderr)
	}
	cmd.Stdout = o.IO().Out
	cmd.Stderr = o.IO().ErrOut
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
