package state

import (
	"fmt"

	"github.com/buyoio/goodies/ptr"
)

func (state *State) ExpandRunners() ([]*Runner, error) {
	if state.Runners == nil {
		return nil, fmt.Errorf("no runners configured")
	}

	runners := []*Runner{}
	for name := range state.Runners {
		r, err := state.ExpandRunner(name)
		if err != nil {
			return nil, err
		}
		runners = append(runners, r)
	}

	return runners, nil
}

// ExpandRunner returns a runner configuration by name
// It will complete and expand the configuration
func (state *State) ExpandRunner(name string) (*Runner, error) {
	if state.Runners == nil {
		return nil, fmt.Errorf("no runners configured")
	}
	runner, ok := state.Runners[name]
	if !ok || runner == nil {
		return nil, fmt.Errorf("runner %s not found", name)
	}

	if runner.Host == nil && runner.HCloud == nil {
		return nil, fmt.Errorf("runner %s has no host or hcloud configuration", name)
	}

	if runner.Host == nil {
		// todo lookup hcloud instance and get host
		return nil, fmt.Errorf("hcloud not implemented")
	}

	// // ensure logs path is absolute
	// runner.Logs.Path = state.ResolvePath(runner.Logs.Path)

	return runner, nil
}

func (r *Runner) ScrubSecrets() interface{} {
	runner := r.DeepCopy()
	if runner.Setup != nil {
		runner.Setup.SCMPlatform.Token = defaultScrub
	}
	if runner.Password != nil {
		runner.Password = ptr.To[string](defaultScrub)
	}
	return runner
}
