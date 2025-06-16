package version

import (
	"fmt"
)

// this will be set at build time
var (
	Version        = "dev"
	CommitHash     = "n/a"
	BuildTimestamp = "n/a"
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}
