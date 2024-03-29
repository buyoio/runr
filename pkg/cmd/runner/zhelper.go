package runner

import (
	"github.com/buyoio/goodies/ssh"
	"github.com/buyoio/runr/pkg/cli/state"
)

func getRunnerSSHClient(o *state.Runr, name string) (*ssh.Client, error) {
	runner, err := o.State().ExpandRunner(name)
	if err != nil {
		return nil, err
	}

	var host string
	if runner.HCloud != nil {
		host = *runner.HCloud
	} else if runner.Host != nil {
		host = *runner.Host
	}

	return ssh.NewClient(o.GetContext(), o.Logger(), &ssh.ClientOptions{
		Host: host,
		User: runner.User,
	}), nil
}
