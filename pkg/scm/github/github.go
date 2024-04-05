package github

import (
	"fmt"
	"strings"

	"github.com/buyoio/b/pkg/binary"
	"github.com/google/go-github/v59/github"
)

func (p *provider) RunnerVersion() (string, error) {
	return binary.GithubLatest(&binary.Binary{
		GitHubRepo: "actions/runner",
	})
}

func (p *provider) RunnerToken() (string, error) {
	if p.runnerToken != nil {
		return p.runnerToken.GetToken(), nil
	}

	var resp *github.Response
	var err error
	if p.options.Repository != nil {
		p.runnerToken, resp, err = p.github.Actions.CreateRegistrationToken(
			p.ctx,
			p.options.Owner,
			*p.options.Repository,
		)
		if err != nil {
			p.logger.Error("Failed to create registration token", "response", resp)
			return "", err
		}
	} else {
		p.runnerToken, resp, err = p.github.Actions.CreateOrganizationRegistrationToken(
			p.ctx,
			p.options.Owner,
		)
		if err != nil {
			p.logger.Error("Failed to create organization registration token", "response", resp)
			return "", err
		}
	}

	return p.runnerToken.GetToken(), nil
}

func (p *provider) RunnerSetup() ([]string, error) {
	version, err := p.RunnerVersion()
	if err != nil {
		return nil, err
	}
	version = version[1:]

	token, err := p.RunnerToken()
	if err != nil {
		return nil, err
	}

	owner := p.Organization()
	repo := p.Repository()
	if repo != "" {
		owner += "/" + repo
	}

	return strings.Split(strings.TrimSpace(fmt.Sprintf(`
curl -Lo actions-runner-linux-x64-%[1]s.tar.gz https://github.com/actions/runner/releases/download/v%[1]s/actions-runner-linux-x64-%[1]s.tar.gz
tar xzf ./actions-runner-linux-x64-%[1]s.tar.gz
rm -f ./actions-runner-linux-x64-%[1]s.tar.gz
chown user:user -R .
runuser -u user -- ./config.sh --url "https://github.com/%[2]s" --token "%[3]s" --labels "$(uname),$(uname -m),$(uname -n),%[4]s"
./svc.sh install user
./svc.sh start
`, version, owner, token, strings.Join(p.options.Labels, ","))), "\n"), nil
}

func (p *provider) TeamSSHKeys(team string) ([]string, error) {
	keys := make([]string, 0)

	users, _, err := p.github.Teams.ListTeamMembersBySlug(p.ctx, p.options.Owner, team, &github.TeamListTeamMembersOptions{
		Role: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		sshKeys, _, err := p.github.Users.ListKeys(p.ctx, *user.Login, &github.ListOptions{
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		for _, key := range sshKeys {
			keys = append(keys, *key.Key)
		}
	}

	return keys, nil
}

func (p *provider) UserSSHKeys(user string) ([]string, error) {
	keys := make([]string, 0)

	sshKeys, _, err := p.github.Users.ListKeys(p.ctx, user, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return nil, err
	}
	for _, key := range sshKeys {
		keys = append(keys, *key.Key)
	}

	return keys, nil
}

func (p *provider) RemoveRunner(name string) error {
	var isOrga bool
	if p.options.Repository == nil {
		isOrga = true
	}

	var runners *github.Runners
	var err error
	// filter by name after https://github.com/google/go-github/pull/3094 is merged and released
	// so we ignore pagination for now
	opts := &github.ListOptions{
		PerPage: 100,
	}

	if isOrga {
		runners, _, err = p.github.Actions.ListOrganizationRunners(p.ctx, p.options.Owner, opts)
	} else {
		runners, _, err = p.github.Actions.ListRunners(p.ctx, p.options.Owner, *p.options.Repository, opts)
	}
	if err != nil {
		return err
	}
	for _, runner := range runners.Runners {
		if strings.HasSuffix(*runner.Name, name) {
			if isOrga {
				_, err = p.github.Actions.RemoveOrganizationRunner(p.ctx, p.options.Owner, runner.GetID())
			} else {
				_, err = p.github.Actions.RemoveRunner(p.ctx, p.options.Owner, *p.options.Repository, runner.GetID())
			}
			if err != nil {
				return err
			}
		}
	}
	// runner not found
	// no error as the runner is not there
	return nil
}
