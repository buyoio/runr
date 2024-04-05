package runner

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/buyoio/goodies/ssh"
	"github.com/buyoio/goodies/streams"
	oses "github.com/buyoio/runr/pkg/OSes"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/buyoio/runr/pkg/scm"

	"github.com/jedib0t/go-pretty/v6/progress"
)

type RunnerOptions struct {
	*state.Runner

	Streams *streams.IO
	Tracker *progress.Tracker

	Force bool
}

type Runner struct {
	ctx     context.Context
	logger  *slog.Logger
	options *RunnerOptions

	provider scm.Provider
	remote   *ssh.Client
}

func NewRunner(ctx context.Context, logger *slog.Logger, options *RunnerOptions) (*Runner, error) {
	logger = logger.
		With(slog.String("pkg", "runner")).
		With("runner", options.Name)

	provider, err := scm.NewProvider(ctx, logger, &scm.ProviderOptions{
		SCMPlatform: scm.SCMPlatform(options.Setup.SCMPlatform),
	})
	if err != nil {
		return nil, err
	}

	client := ssh.NewClient(ctx, logger, &ssh.ClientOptions{
		Host:      *options.Host,
		User:      options.User,
		IOStreams: options.Streams,
	})

	return &Runner{
		ctx:      ctx,
		logger:   logger,
		options:  options,
		provider: provider,
		remote:   client,
	}, nil
}

func (r *Runner) Close() error {
	if r.remote != nil {
		r.remote.Close()
	}
	if r.options.Streams != nil && r.options.Streams.Closer != nil {
		r.options.Streams.Closer()
	}
	return nil
}

// todo not sure if this is the right place for this
// not quite runner specific
// also conformation from the user? as disks are wiped
func (r *Runner) HetznerInstallimage() error {
	if r.options.Setup == nil || r.options.Setup.Installimage == nil {
		return nil
	}

	if _, err := r.remote.Connect(); err != nil {
		return err
	}
	defer r.remote.Close()

	os, _, err := oses.GetOS(r.ctx, r.logger, r.remote)
	if err != nil {
		return err
	}

	r.logger.InfoContext(r.ctx, "Installing hetzner installimage")
	content, err := r.options.Setup.Installimage.Content()
	if err != nil {
		return err
	}
	if err := os.HetznerInstallimage(
		fmt.Sprintf(
			content,
			r.options.Name,
		),
	); err != nil {
		return err
	}

	return nil
}

func (r *Runner) Upgrade() error {
	if _, err := r.remote.Connect(); err != nil {
		return err
	}
	defer r.remote.Close()
	r.logger.InfoContext(r.ctx, "Connected to remote")

	os, info, err := oses.GetOS(r.ctx, r.logger, r.remote)
	if err != nil {
		return err
	}

	if info.HetznerRescueMode != "" {
		r.logger.WarnContext(r.ctx, "Hetzner rescue mode detected")
		if !r.options.Force {
			return fmt.Errorf("hetzner rescue mode detected, use --force to continue")
		}

		if err := r.HetznerInstallimage(); err != nil {
			return err
		}
	}

	r.logger.InfoContext(r.ctx, "Ensuring system")
	if err := os.EnsureSystem(); err != nil {
		return err
	}

	r.logger.InfoContext(r.ctx, "Ensuring user")
	if err := os.EnsureUser(
		r.options.Setup.User,
		"",
	); err != nil {
		return err
	}

	r.logger.InfoContext(r.ctx, "Ensuring docker")
	if err := os.EnsureDocker(r.options.User); err != nil {
		return err
	}

	r.logger.InfoContext(r.ctx, "Ensuring cron")
	if err := os.DockerCronPruneSystem(); err != nil {
		return err
	}

	if r.options.Pre != nil {
		r.logger.InfoContext(r.ctx, "Running exec pre")
		if err := os.Exec(*r.options.Pre); err != nil {
			return err
		}
	}

	r.logger.InfoContext(r.ctx, "Stopping all runners")
	if err := os.DockerStopAll(); err != nil {
		return err
	}

	r.logger.InfoContext(r.ctx, "Gathering runner token")
	scm := r.options.Setup.SCMPlatform
	token, err := r.provider.RunnerToken()
	if err != nil {
		return err
	}

	r.logger.InfoContext(r.ctx, "Provisioning runners")
	for name, docker := range r.options.Docker {
		r.logger.DebugContext(r.ctx, "Building docker image")
		file, err := docker.Dockerfile.Content()
		if err != nil {
			return err
		}
		if err := os.DockerBuildRunner(file); err != nil {
			return err
		}

		orgrepo := scm.Organization
		if scm.Repository != nil {
			orgrepo = orgrepo + "/" + *scm.Repository
		}

		var labels []string
		labels = append(labels, docker.Labels...)
		labels = append(labels, info.OS, info.Architecture, info.Hostname, name)

		for i := 0; i < docker.Quantity; i++ {
			r.logger.DebugContext(r.ctx, "Starting docker runner", "name", name, "index", i)
			if err := os.DockerStartRunner(
				fmt.Sprintf("%s-%d", name, i),
				token,
				strings.Join(labels, ","),
				orgrepo,
				r.options.Setup.RunnerWorkDir,
			); err != nil {
				return err
			}
		}
	}

	if r.options.Post != nil {
		r.logger.InfoContext(r.ctx, "Running exec post")
		if err := os.Exec(*r.options.Post); err != nil {
			return err
		}
	}

	return nil
}
