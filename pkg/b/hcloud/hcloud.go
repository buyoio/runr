package hcloud

import (
	"os/exec"

	"github.com/buyoio/b/pkg/binaries"
	"github.com/buyoio/b/pkg/binaries/hcloud"
	"github.com/buyoio/b/pkg/binary"
	"github.com/buyoio/goodies/cmdutil"
	"github.com/buyoio/runr/pkg/cli/state"
	"github.com/spf13/cobra"
)

type Hcloud struct {
	*binary.Binary
}

func Binary(options *binaries.BinaryOptions) *Hcloud {
	hcloud := &Hcloud{
		Binary: hcloud.Binary(options),
	}
	return hcloud
}

func BinaryCmd(o *state.Runr, cmd *cobra.Command) (*Hcloud, error) {
	hcloud := Binary(&binaries.BinaryOptions{
		Context: o.GetContext(),
		Envs: map[string]string{
			"HCLOUD_TOKEN": o.State().HCloud.Token,
		},
	})
	if !hcloud.BinaryExists() {
		return nil, cmdutil.UsageErrorf(cmd, "hcloud not found, please install it with `runr b -ia`")
	}
	return hcloud, nil
}

func (hcloud *Hcloud) HasValidToken() error {
	_, err := hcloud.Exec("server", "list")
	return err
}

func (h *Hcloud) ServerCreate(opts *state.HCloudServerCreateOpts) (*exec.Cmd, error) {
	args, stdin, err := opts.Args()
	if err != nil {
		return nil, err
	}

	c := h.Cmd(args...)
	c.Stdin = stdin
	return c, nil
}

func (h *Hcloud) ServerDelete(server string) *exec.Cmd {
	c := h.Cmd("server", "delete", server)
	return c
}

func (h *Hcloud) ServerList() *exec.Cmd {
	c := h.Cmd("server", "list", "--selector", "runr")
	return c
}
