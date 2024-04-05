/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"github.com/spf13/cobra"

	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/goodies/output"
	"github.com/buyoio/goodies/templates"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/i18n"
)

var (
	// Magic variables set by goreleaser
	version = "v0.0.1" // x-release-please-version
	commit  = ""
	date    = "01.01.1970"
)

var info = &Info{
	Version:   version,
	GitCommit: commit,
	BuildDate: date,
}

// VersionOptions is a struct to support version command
type VersionOptions struct {
	*state.Runr
}

// NewOptions returns initialized Options
func NewVersionOptions(runr *state.Runr) *VersionOptions {
	return &VersionOptions{
		Runr: runr,
	}

}

// NewCmdVersion returns a cobra command for fetching versions
func NewCmdVersion(runr *state.Runr) *cobra.Command {
	o := NewVersionOptions(runr)
	cmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("Print version information"),
		Long:  i18n.T("Print version information for the current context."),
		Example: templates.Examples(i18n.T(`
			# Print the client and server versions for the current context
			runr version
		`)),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML(), output.OptionFormat())
	return cmd
}

// Complete completes all the required options
func (o *VersionOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Validate validates the provided options
func (o *VersionOptions) Validate() error {
	return nil
}

// Run executes version command
func (o *VersionOptions) Run() error {
	return o.IO().Print(&info)
}
