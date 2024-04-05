package state

import (
	"log/slog"

	"github.com/buyoio/goodies/streams"
	"github.com/spf13/cobra"
)

// type File string
type File string

type DockerRunner struct {
	Dockerfile *Data    `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`
	Labels     []string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Quantity   int      `json:"quantity,omitempty" yaml:"quantity,omitempty"`
}

type RunnerSetup struct {
	SCMPlatform `json:",inline" yaml:",inline"`

	User          string   `json:"user"`
	Packages      []string `json:"packages,omitempty" yaml:"packages,omitempty"`
	RunnerWorkDir string   `json:"runner_work_dir,omitempty" yaml:"runner_work_dir,omitempty"`
	Installimage  *Data    `json:"installimage,omitempty" yaml:"installimage,omitempty"`
}

type Runner struct {
	// will be filled with the name of the runner from map
	Name *string `json:"-" yaml:"-"`

	Host     *string `json:"host,omitempty" yaml:"host,omitempty"`
	HCloud   *string `json:"hcloud,omitempty" yaml:"hcloud,omitempty"`
	User     string  `json:"user,omitempty" yaml:"user,omitempty"`
	Auth     string  `json:"auth,omitempty" yaml:"auth,omitempty"`
	Password *string `json:"password,omitempty" yaml:"password,omitempty"` // Not saved to config

	Setup  *RunnerSetup             `json:"setup,omitempty" yaml:"setup,omitempty"`
	Docker map[string]*DockerRunner `json:"docker,omitempty" yaml:"docker,omitempty"`
	Pre    *string                  `json:"pre,omitempty" yaml:"pre,omitempty"`
	Post   *string                  `json:"post,omitempty" yaml:"post,omitempty"`
}

type Arg struct {
	File  *File    `json:"file,omitempty" yaml:"file,omitempty"`
	Exec  *File    `json:"exec,omitempty" yaml:"exec,omitempty"`
	Value *string  `json:"value,omitempty" yaml:"value,omitempty"`
	Envs  []string `json:"envs,omitempty" yaml:"envs,omitempty"`
}

type Args []*Arg

type Data struct {
	File File `json:"file,omitempty" yaml:"file,omitempty"`
	Args Args `json:"args,omitempty" yaml:"args,omitempty"`

	tmpl *GoTemplate `json:"-" yaml:"-"`
}

type HCloudServerCreateOpts struct {
	Setup *SCMPlatform `json:"setup,omitempty" yaml:"setup,omitempty"`
	// should not be provided by the config file
	Name             string   `json:"-"`
	Type             string   `json:"type"`
	Image            string   `json:"image"`
	SSHKeys          []string `json:"ssh_keys,omitempty" yaml:"ssh_keys,omitempty"`
	Location         *string  `json:"location,omitempty" yaml:"location,omitempty"`
	Datacenter       *string  `json:"datacenter,omitempty" yaml:"datacenter,omitempty"`
	UserData         *Data    `json:"user_data"`
	StartAfterCreate bool     `json:"start_after_create,omitempty" yaml:"start_after_create,omitempty"`
	Labels           []string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Automount        bool     `json:"automount,omitempty" yaml:"automount,omitempty"`
	Volumes          []string `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Networks         []string `json:"networks,omitempty" yaml:"networks,omitempty"`
	Firewalls        []string `json:"firewalls,omitempty" yaml:"firewalls,omitempty"`
	PlacementGroup   *string  `json:"placement_group,omitempty" yaml:"placement_group,omitempty"`
}

type Environments struct {
	HCloud map[string]*HCloudServerCreateOpts `json:"hcloud,omitempty" yaml:"hcloud,omitempty"`
}

type HCloud struct {
	Token string `json:"token,omitempty" yaml:"token,omitempty"` // Not saved to config
}

type Hetzner struct {
	Token string `json:"token,omitempty" yaml:"token,omitempty"` // Not saved to config
}

type Github struct {
	Token        string  `json:"token,omitempty" yaml:"token,omitempty"` // Not saved to config
	Organization string  `json:"organization,omitempty" yaml:"organization,omitempty"`
	Repository   *string `json:"repository,omitempty" yaml:"repository,omitempty"`
}

type Gitlab struct {
	Token        string  `json:"token,omitempty" yaml:"token,omitempty"` // Not saved to config
	Organization string  `json:"organization,omitempty" yaml:"organization,omitempty"`
	Repository   *string `json:"repository,omitempty" yaml:"repository,omitempty"`
}

type SCMPlatform struct {
	Platform     string  `json:"platform,omitempty" yaml:"platform,omitempty"`
	Organization string  `json:"organization,omitempty" yaml:"organization,omitempty"`
	Repository   *string `json:"repository,omitempty" yaml:"repository,omitempty"`
	Token        string  `json:"token,omitempty" yaml:"token,omitempty"` // Not saved to config
}

type Logs struct {
	Path    *string `json:"path,omitempty" yaml:"path,omitempty"`
	Level   *string `json:"level,omitempty" yaml:"level,omitempty"`
	Verbose bool    `json:"verbose,omitempty" yaml:"verbose,omitempty"`
}

type State struct {
	Logs         *Logs              `json:"logs,omitempty" yaml:"logs,omitempty"`
	Runners      map[string]*Runner `json:"runners,omitempty" yaml:"runners,omitempty"`
	Environments *Environments      `json:"environments,omitempty" yaml:"environments,omitempty"`

	Github  *Github  `json:"github,omitempty" yaml:"github,omitempty"`
	Gitlab  *Gitlab  `json:"gitlab,omitempty" yaml:"gitlab,omitempty"`
	HCloud  *HCloud  `json:"hcloud,omitempty" yaml:"hcloud,omitempty"`
	Hetzner *Hetzner `json:"hetzner,omitempty" yaml:"hetzner,omitempty"`
}

type Runr struct {
	// This reflects the state of the config file itself
	fileState *State
	// This reflects the state of the config file after all defaults and secrets have been applied
	aggrState *State
	// This reflects the state of the config file after all file references have been loaded
	// expdState *State

	cmd    *cobra.Command
	io     *streams.IO
	logger *slog.Logger

	config string
	quiet  bool

	closer func() error
}
