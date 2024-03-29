package state

import (
	"fmt"
	"os"
)

const (
	envPrefix = "LOK8S"
)

func (state *State) lookupEnvs() {
	if state.Github == nil {
		state.Github = &Github{}
	}
	if state.Github.Token == "" {
		state.Github.Token = os.Getenv("GITHUB_TOKEN")
	} else {
		os.Setenv("GITHUB_TOKEN", state.Github.Token)
	}

	if state.Gitlab == nil {
		state.Gitlab = &Gitlab{}
	}
	if state.Gitlab.Token == "" {
		state.Gitlab.Token = os.Getenv("GITLAB_TOKEN")
	} else {
		os.Setenv("GITLAB_TOKEN", state.Gitlab.Token)
	}

	if state.HCloud == nil {
		state.HCloud = &HCloud{}
	}
	if state.HCloud.Token == "" {
		state.HCloud.Token = os.Getenv("HCLOUD_TOKEN")
	} else {
		os.Setenv("HCLOUD_TOKEN", state.HCloud.Token)
	}

	if state.Hetzner == nil {
		state.Hetzner = &Hetzner{}
	}
	if state.Hetzner.Token == "" {
		state.Hetzner.Token = os.Getenv("HETZNER_TOKEN")
	} else {
		os.Setenv("HETZNER_TOKEN", state.Hetzner.Token)
	}

	if state.Runners != nil {
		for name, runner := range state.Runners {
			if runner.Password == nil {
				passwd := os.Getenv(
					fmt.Sprintf("%s_RUNNER_%s_PASSWORD", envPrefix, name),
				)
				if passwd != "" {
					runner.Password = &passwd
				}
			}
		}
	}
}
