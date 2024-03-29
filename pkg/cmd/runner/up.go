package runner

import (
	"fmt"
	"sync"
	"time"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/logs"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/goodies/progress"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/buyoio/runr/pkg/runner"
	"github.com/spf13/cobra"
)

type CmdUpgradeOptions struct {
	*state.Runr
	runner []*state.Runner

	All   bool
	Name  string
	Force bool
}

func NewCmdUpgrade(runr *state.Runr) *cobra.Command {
	o := &CmdUpgradeOptions{
		Runr: runr,
	}
	cmd := &cobra.Command{
		Use:   "up NAME | --all",
		Short: i18n.T("Ensures that the runner"),
		Long:  i18n.T("Ensures that the runner is installed and up to date"),
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

func (o *CmdUpgradeOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", false, "Upgrade all runners")
	cmd.Flags().BoolVar(&o.Force, "force", false, "Force upgrade, this will wipe disks if necessary (Hetzner rescue mode)")
}

func (o *CmdUpgradeOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		o.Name = args[0]
	}
	return nil
}
func (o *CmdUpgradeOptions) Validate(cmd *cobra.Command) error {
	if o.Name == "" && !o.All {
		return cmdutil.UsageErrorf(cmd, i18n.T("NAME is required"))
	}

	if o.All {
		runners, err := o.State().ExpandRunners()
		if err != nil {
			return err
		}
		o.runner = runners
		return nil
	}

	r, err := o.State().ExpandRunner(o.Name)
	if err != nil {
		return err
	}
	o.runner = append(o.runner, r)

	return nil
}
func (o *CmdUpgradeOptions) Run() error {
	var wg sync.WaitGroup
	wg.Add(len(o.runner))

	pw := progress.NewWriter(progress.StyleUnkown, o.IO().Out)
	go pw.Render()
	defer pw.Stop()

	for _, r := range o.runner {
		streams, err := o.State().Logs.Raw(*r.Name)
		if err != nil {
			return err
		}

		go func() {
			logger := o.Logger()
			logs.ProgressTracker(logger, pw.AddTracker(
				fmt.Sprintf("Upgrading runner %s", *r.Name),
				12,
			))

			rr, err := runner.NewRunner(
				o.GetContext(),
				logger,
				&runner.RunnerOptions{
					Runner:  r,
					Force:   o.Force,
					Streams: streams,
				},
			)
			if err == nil {
				defer rr.Close()
				err = rr.Upgrade()
			}
			logs.ProgressDone(logger, fmt.Sprintf("Upgrade of runner %s done", *r.Name), err)
			wg.Done()
		}()
	}

	wg.Wait()
	time.Sleep(500 * time.Millisecond)
	return nil
}
