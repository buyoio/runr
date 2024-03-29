package linux

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	gssh "golang.org/x/crypto/ssh"

	"github.com/buyoio/goodies/ssh"
)

type Linux struct {
	ctx    context.Context
	logger *slog.Logger
	remote *ssh.Client
}

func NewLinux(ctx context.Context, logger *slog.Logger, remote *ssh.Client) *Linux {
	return &Linux{
		ctx:    ctx,
		logger: logger.With(slog.String("os", "linux")),
		remote: remote,
	}
}

func (linux *Linux) Prerun() error {
	return linux.remote.Run(fmt.Sprintf(`cat <<EOF | sudo tee /usr/local/bin/lok8s
#!/usr/bin/env bash
set -euo pipefail
%s`, lok8sCommand))
}

func (linux *Linux) executeScript(name string, args ...string) error {
	cmd := fmt.Sprintf(
		"%s; %s %s",
		lok8sCommand,
		name,
		strings.Join(args, " "),
	)
	linux.logger.Debug(fmt.Sprintf("Executing script: %s", name), "args", args)
	return linux.remote.Run(cmd)
}

func (linux *Linux) EnsureSystem() error {
	packages := []string{
		"apt-transport-https",
		"ca-certificates",
		"gnupg",
		"curl",
		"iptables",
		"curl",
		"vim",
		"fail2ban",
	}
	err := linux.executeScript("system::ensure", strings.Join(packages, " "))
	linux.remote.PossibleRebootWait()
	return err
}

func (linux *Linux) EnsureUser(username string, password string) error {
	return linux.executeScript("system::user", username, password)
}

func (linux *Linux) EnsureDocker(username string) error {
	return linux.executeScript("system::docker", username)
}

// 0 3 1 * *
func (linux *Linux) DockerCronPruneSystem() error {
	return linux.executeScript("docker::cron")
}

func (linux *Linux) DockerStopAll() error {
	return linux.executeScript("docker::stop")
}

func (linux *Linux) DockerBuildRunner(dockerfile string) error {
	return linux.executeScript("docker::build", dockerfile)
}

func (linux *Linux) DockerStartRunner(name string, token string, labels string, orgrepo string, path string) error {
	return linux.executeScript("docker::start", name, token, labels, orgrepo, path)
}

func (linux *Linux) HetznerInstallimage(installimage string) error {
	err := linux.executeScript("hetzner::install", installimage)
	if err != nil {
		if e, ok := err.(*gssh.ExitError); ok && e.ExitStatus() == 127 {
			// nothing to do
			return nil
		}
		return err
	}
	linux.remote.PossibleRebootWait()

	return nil
}
