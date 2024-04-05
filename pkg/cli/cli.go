package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/buyoio/goodies/streams"
	"github.com/buyoio/goodies/templates"
	"github.com/buyoio/runr/pkg/cli/state"

	// "github.com/buyoio/runr/pkg/cmd/config"
	bhcloud "github.com/buyoio/runr/pkg/b/hcloud"
	"github.com/buyoio/runr/pkg/cmd/hcloud"
	"github.com/buyoio/runr/pkg/cmd/runner"
	"github.com/buyoio/runr/pkg/cmd/version"
	"github.com/buyoio/runr/pkg/i18n"

	"github.com/buyoio/b/pkg/binaries"
	"github.com/buyoio/b/pkg/binary"
	b "github.com/buyoio/b/pkg/cli"
)

func NewRunrCommand() *cobra.Command {
	cmd, o := state.NewState(
		&cobra.Command{
			Use:              "runr",
			TraverseChildren: true,
			SilenceUsage:     true,
			// SilenceErrors:         true,
			DisableFlagsInUseLine: true,
			Short:                 i18n.T("Runr makes it easy to set up runners for CI/CD pipelines"),
			Long: i18n.T(templates.LongDesc(`
				Runr makes it easy to set up runners for CI/CD pipelines.

				For more information, visit https://docs.runr.cloud
			`)),
		},
	)

	bCmd := b.NewCmdBinary(&b.CmdBinaryOptions{
		IO:       &streams.IO{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		NoConfig: true,
		Binaries: []*binary.Binary{
			bhcloud.Binary(&binaries.BinaryOptions{
				Context: cmd.Context(),
			}).Binary,
		},
	})
	bCmd.Short = i18n.T("Manage all required binaries")
	templates.AddGroup(cmd, "local", "Local environment",
		bCmd,
	)

	templates.AddGroup(cmd, "ci", "Runr",
		runner.NewCmdRunner(o),
		hcloud.NewCmdHcloud(o),
	)

	cmd.AddCommand(version.NewCmdVersion(o))

	return cmd
}
