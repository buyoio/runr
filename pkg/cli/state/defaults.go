package state

import (
	"context"
	"fmt"
	"log/slog"

	"dario.cat/mergo"
	"github.com/buyoio/goodies/ptr"
)

const (
	defaultMapKey = "_"
	defaultScrub  = "********"
)

var defaultConfig = &State{
	Logs: &Logs{
		Path:    ptr.To[string](".logs"),
		Verbose: false,
		Level:   ptr.To[string]("debug"),
	},
	Runners: map[string]*Runner{
		defaultMapKey: {
			User: "root",
			Auth: "ssh-agent",
			Setup: &RunnerSetup{
				SCMPlatform: SCMPlatform{
					Platform: "github",
				},
				User:          "runner",
				RunnerWorkDir: "/home/runner",
				Packages:      []string{},
			},
			Docker: map[string]*DockerRunner{
				defaultMapKey: {
					Dockerfile: &Data{
						File: File(".github/Dockerfile"),
					},
					Quantity: 1,
				},
			},
		},
	},
	Environments: &Environments{
		HCloud: map[string]*HCloudServerCreateOpts{
			defaultMapKey: {
				Setup: &SCMPlatform{
					Platform: "github",
				},
				Type:     "cpx41",
				Image:    "ubuntu-20.04",
				Location: ptr.To[string]("hel1"),
			},
		},
	},
}

func (state *State) aggregate(ctx context.Context, logger *slog.Logger) (*State, error) {
	// todo maybe generate deep copy files instead of reflect
	org := defaultConfig.DeepCopy()
	dst := state.DeepCopy()

	if err := mergo.Merge(dst, org); err != nil {
		return nil, err
	}

	dst.lookupEnvs()
	if err := dst.mergeSecrets(); err != nil {
		return nil, err
	}
	if err := dst.aggregateRunners(ctx, logger); err != nil {
		return nil, err
	}
	if err := dst.aggregateEnvironments(ctx, logger); err != nil {
		return nil, err
	}

	return dst, nil
}

func (dst *State) aggregateRunners(ctx context.Context, logger *slog.Logger) error {
	for name, runner := range dst.Runners {
		if name == defaultMapKey {
			continue
		}

		err := mergo.Merge(runner, dst.Runners[defaultMapKey])
		if err != nil {
			return err
		}
		runner.Name = &name

		if err := dst.mergeSCMPlatform(&runner.Setup.SCMPlatform); err != nil {
			return err
		}

		if runner.Setup.Installimage != nil {
			runner.Setup.Installimage.setGoTemplate(ctx, logger, &GoTemplateOptions{
				Runner:      runner,
				SCMPlatform: &runner.Setup.SCMPlatform,
			})
		}

		defaultDocker, ok := runner.Docker[defaultMapKey]
		for k, docker := range runner.Docker {
			if k == defaultMapKey {
				continue
			}
			if ok {
				err := mergo.Merge(docker, defaultDocker)
				if err != nil {
					return err
				}
			}
			if docker.Dockerfile != nil {
				docker.Dockerfile.setGoTemplate(ctx, logger, &GoTemplateOptions{
					Runner:      runner,
					SCMPlatform: &runner.Setup.SCMPlatform,
				})
			}
		}
		delete(runner.Docker, defaultMapKey)
	}
	delete(dst.Runners, defaultMapKey)
	return nil
}

func (dst *State) aggregateEnvironments(ctx context.Context, logger *slog.Logger) error {
	for name, env := range dst.Environments.HCloud {
		if name == defaultMapKey {
			continue
		}

		err := mergo.Merge(env, dst.Environments.HCloud[defaultMapKey])
		if err != nil {
			return err
		}

		if err := dst.mergeSCMPlatform(env.Setup); err != nil {
			return err
		}

		if env.UserData != nil {
			env.UserData.setGoTemplate(ctx, logger, &GoTemplateOptions{
				SCMPlatform: env.Setup,
			})
		}
	}
	delete(dst.Environments.HCloud, defaultMapKey)
	return nil
}

func (dst *State) mergeSCMPlatform(scm *SCMPlatform) error {
	switch scm.Platform {
	case "github":
		if err := mergo.Merge(scm, dst.Github.SCMPlatform()); err != nil {
			return err
		}
	case "gitlab":
		if err := mergo.Merge(scm, dst.Gitlab.SCMPlatform()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported SCM platform: %s", scm.Platform)
	}

	if scm.Token == "" {
		return fmt.Errorf("%s token not set", scm.Platform)
	}

	return nil
}
