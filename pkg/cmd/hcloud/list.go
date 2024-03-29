package hcloud

import (
	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/runr/pkg/b/hcloud"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
)

type CmdLsOptions struct {
	*state.Runr
	hcloud *hcloud.Hcloud

	outFlags []string
}

func NewCmdLs(runr *state.Runr) *cobra.Command {
	o := &CmdLsOptions{
		Runr: runr,
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("List hcloud runners"),
		Long:  i18n.T("List hcloud runners"),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}
	o.AddFlags(cmd)
	output.AddFlag(cmd,
		output.OptionNoHeader(),
		output.OptionColumns([]string{
			"id", "name", "status", "ipv4", "ipv6", "private_net", "datacenter", "age",
		}),
		output.OptionJSON(),
		output.OptionYAML(),
	)

	return cmd
}

func (o *CmdLsOptions) AddFlags(cmd *cobra.Command) {

}

func (o *CmdLsOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error
	if o.hcloud, err = hcloud.BinaryCmd(o.Runr, cmd); err != nil {
		return err
	}

	o.outFlags = output.GetFlag(cmd)

	return nil
}
func (o *CmdLsOptions) Validate() error {
	if err := o.hcloud.HasValidToken(); err != nil {
		return err
	}

	return nil
}
func (o *CmdLsOptions) Run() error {
	cmd := o.hcloud.ServerList()
	if o.outFlags != nil {
		cmd.Args = append(cmd.Args, o.outFlags...)
	}
	cmd.Stdout = o.IO().Out
	cmd.Stderr = o.IO().ErrOut
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
