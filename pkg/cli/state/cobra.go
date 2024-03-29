package state

import (
	"fmt"
	"io"
	"os"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func (l *Runr) persistentPreRunE(cmd *cobra.Command, args []string) error {
	// Configfile handling
	file := cmdutil.GetFlagString(cmd, "config")
	if file == "" {
		file = getDefaultConfigPath()
	}

	if _, err := os.Stat(file); err != nil {
		return nil
	}
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	os.Setenv("LOK8S_CONFIG", file)

	err = yaml.Unmarshal(yamlFile, l.fileState)
	if err != nil {
		return fmt.Errorf(i18n.T("Failed to unmarshal config file (%s)\n\n%s"), file, err)
	}

	l.aggrState, err = l.fileState.aggregate(l.GetContext(), l.Logger())
	if err != nil {
		return fmt.Errorf(i18n.T("Failed to aggregate config file (%s)\n\n%s"), file, err)
	}

	// Verbose - Output relevant
	if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
		l.io.Out = io.Discard
	}
	// if verbose we switch to slog json output - no need for standard output
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		l.io.Out = io.Discard
	}

	l.io.OutFlags = output.FlagsForCommand(cmd)

	return nil
}

func (state *Runr) persistentPostRunE(cmd *cobra.Command, args []string) error {
	if state.closer != nil {
		state.closer()
	}
	return nil
}
