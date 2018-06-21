package version

import (
	"fmt"

	"github.com/sniperkit/hub/plugin/vcs/local/git"
)

const ProgramName string = "xhub"

var (
	Version      string = "2.4.0"
	CommitHash   string = ""
	CommitID     string = ""
	CommitUnix   string = ""
	BuildVersion string = "2015.6.2-6-gfd7e2d1-dev"
	BuildTime    string = "2015-06-16-0431 UTC"
	BuildCount   string = ""
	BuildUnix    string = ""
	BranchName   string = ""
)

func FullVersion() (string, error) {
	gitVersion, err := git.Version()
	if err != nil {
		gitVersion = "git version (unavailable)"
	}
	return fmt.Sprintf("%s\nhub version %s", gitVersion, Version), err
}
