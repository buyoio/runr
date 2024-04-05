package oses

import (
	"context"
	"log/slog"
	"strings"

	"github.com/buyoio/goodies/ssh"
	"github.com/buyoio/runr/pkg/OSes/linux"
)

type OS interface {
	Prerun() error

	EnsureSystem() error
	EnsureUser(username string, password string) error
	EnsureDocker(username string) error
	DockerCronPruneSystem() error
	DockerStopAll() error
	DockerBuildRunner(dockerfile string) error
	DockerStartRunner(name string, token string, labels string, orgrepo string, path string) error

	Exec(command string) error

	// Hetzner specific
	HetznerInstallimage(installimage string) error
}

func GetOS(ctx context.Context, logger *slog.Logger, remote *ssh.Client) (OS, *ssh.RemoteInfo, error) {
	info, err := remote.RemoteInfo()
	if err != nil {
		return nil, nil, err
	}
	switch strings.ToLower(info.OS) {
	case "linux":
		linux := linux.NewLinux(ctx, logger, remote)
		if err := linux.Prerun(); err != nil {
			return nil, info, err
		}
		return linux, info, nil
	}
	return nil, info, nil
}
