package version

import (
	"fmt"
)

type Info struct {
	// Version is the semantic version of the build.
	// https://semver.org/
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return fmt.Sprintf("%s+%s", info.Version, info.GitCommit)
}
