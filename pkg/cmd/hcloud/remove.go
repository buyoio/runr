package hcloud

import (
	"bytes"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/runr/pkg/b/hcloud"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/buyoio/runr/pkg/scm"
	"github.com/spf13/cobra"

	"github.com/tidwall/gjson"
)

type CmdRmOptions struct {
	*state.Runr
	hcloud *hcloud.Hcloud

	server         string
	environmentRef string
	runnerRef      string
	serverRef      string
	scm            scm.Provider
}

func NewCmdRm(runr *state.Runr) *cobra.Command {
	o := &CmdRmOptions{
		Runr: runr,
	}
	cmd := &cobra.Command{
		Use:   "remove NAME",
		Short: i18n.T("Remove a hcloud runner"),
		Long:  i18n.T("Remove a hcloud runner"),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate(cmd))
			cmdutil.CheckErr(o.Run())
		},
	}

	return cmd
}

func (o *CmdRmOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error
	if o.hcloud, err = hcloud.BinaryCmd(o.Runr, cmd); err != nil {
		return err
	}

	if len(args) > 0 {
		o.server = args[0]
	}

	return nil
}
func (o *CmdRmOptions) Validate(cmd *cobra.Command) error {
	if err := o.hcloud.HasValidToken(); err != nil {
		return err
	}

	var out bytes.Buffer
	c := o.hcloud.ServerList()
	c.Args = append(c.Args, "-o", "json")
	c.Stdout = &out
	c.Stderr = o.IO().ErrOut
	if err := c.Run(); err != nil {
		return err
	}
	result := gjson.Get(out.String(), `#(labels|@keys.#(="runr"))#.{id,name,labels.runr}`)
	if !result.Exists() || len(result.Array()) == 0 {
		return cmdutil.UsageErrorf(cmd, "no runr runners found")
	}
	for _, server := range result.Array() {
		if server.Get("name").String() == o.server || server.Get("id").String() == o.server {
			o.environmentRef = server.Get("runr").String()
			o.runnerRef = server.Get("name").String()
			o.serverRef = server.Get("id").String()
			break
		}
	}
	if o.environmentRef == "" {
		return cmdutil.UsageErrorf(cmd, "server %s not found", o.server)
	}

	e, ok := o.State().Environments.HCloud[o.environmentRef]
	if !ok {
		return cmdutil.UsageErrorf(cmd, "environment %s not found", o.environmentRef)
	}

	var err error
	if o.scm, err = e.Setup.NewProvider(o.GetContext(), o.Logger()); err != nil {
		return err
	}

	return nil
}
func (o *CmdRmOptions) Run() error {
	if err := o.scm.RemoveRunner(o.runnerRef); err != nil {
		return err
	}

	// delete github
	cmd := o.hcloud.ServerDelete(o.serverRef)
	cmd.Stdout = o.IO().Out
	cmd.Stderr = o.IO().ErrOut
	return cmd.Run()
}
