package linux

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"

	gssh "golang.org/x/crypto/ssh"

	scp "github.com/bramvdbogaerde/go-scp"
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

func (linux *Linux) scpCopy(content string, dest string) error {
	var buff bytes.Buffer
	_, err := buff.Write([]byte(content))
	if err != nil {
		return err
	}
	sshClient, err := linux.remote.RawClient()
	if err != nil {
		return err
	}
	scpClient, err := scp.NewClientBySSH(sshClient)
	if err != nil {
		return err
	}
	return scpClient.CopyFile(linux.ctx, &buff, dest, "0755")
}

func (linux *Linux) Prerun() error {
	err := linux.scpCopy(lok8sScript, "/usr/local/bin/lok8s")
	if err != nil {
		return err
	}

	return linux.remote.Run(`
		curl -fsSL https://min.arg.sh > /usr/local/bin/argsh
		sudo chmod +x /usr/local/bin/argsh
		chmod +x /usr/local/bin/lok8s
	`)
}

func (linux *Linux) Exec(command string) error {
	linux.logger.Debug("Executing command", "command", command)
	return linux.remote.Run(command)
}

func (linux *Linux) executeScript(name string, args ...string) error {
	cmd := fmt.Sprintf(
		"lok8s %s %s",
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
	err := linux.executeScript("system-ensure", strings.Join(packages, " "))
	linux.remote.PossibleRebootWait()
	return err
}

func (linux *Linux) EnsureUser(username string, password string) error {
	return linux.executeScript("system-user", username, password)
}

func (linux *Linux) EnsureDocker(username string) error {
	return linux.executeScript("system-docker", username)
}

// 0 3 1 * *
func (linux *Linux) DockerCronPruneSystem() error {
	return linux.executeScript("docker-cron")
}

func (linux *Linux) DockerStopAll() error {
	return linux.executeScript("docker-stop")
}

func (linux *Linux) DockerBuildRunner(dockerfile string) error {
	err := linux.scpCopy(dockerfile, "/tmp/Dockerfile")
	if err != nil {
		return err
	}
	return linux.executeScript("docker-build", "/tmp/Dockerfile")
}

func (linux *Linux) DockerStartRunner(name string, token string, labels string, orgrepo string, path string) error {
	return linux.executeScript("docker-start", name, orgrepo, token, labels, path)
}

func (linux *Linux) HetznerInstallimage(installimage string) error {
	err := linux.executeScript("hetzner-install", installimage)
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
