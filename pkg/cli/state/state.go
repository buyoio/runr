package state

import (

	// "github.com/hetznercloud/hcloud-go/hcloud"

	"fmt"
	"os"
	"path"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/buyoio/goodies/git"
	"github.com/buyoio/goodies/streams"
	"github.com/buyoio/runr/pkg/i18n"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewState(cmd *cobra.Command) (*cobra.Command, *Runr) {
	streams := &streams.IO{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	l := &Runr{
		cmd:    cmd,
		quiet:  false,
		config: getDefaultConfigPath(),
		io:     streams,
		fileState: &State{
			Logs: &Logs{},
		},
	}
	cmd.PersistentPreRunE = l.persistentPreRunE
	cmd.PersistentPostRunE = l.persistentPostRunE

	config := l.config
	if cwd, err := os.Getwd(); err == nil {
		config, _ = filepath.Rel(cwd, config)
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&l.config, "config", config, "Path to the config file")
	flags.BoolVarP(&l.quiet, "quiet", "q", false, "Only print error messages")

	l.fileState.Logs.AddFlags(flags)

	return cmd, l
}

func (state *State) mergeSecrets() error {
	secretsFile := path.Join(getDefaultConfigPath(), secretFileName)
	if _, err := os.Stat(secretsFile); err != nil {
		return nil
	}

	if !git.IsIgnored(secretsFile) {
		return fmt.Errorf(i18n.T("Secrets file (%s) is not ignored by git"), secretsFile)
	}

	secrets := &State{}
	yamlFile, err := os.ReadFile(secretsFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, secrets)
	if err != nil {
		return fmt.Errorf(i18n.T("Failed to unmarshal secrets file (%s)\n\n%s"), secretsFile, err)
	}
	return mergo.Merge(state, secrets, mergo.WithOverride)
}
