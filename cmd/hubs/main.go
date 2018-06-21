// +build go1.8

package main

import (
	"os"

	// core
	"github.com/sniperkit/hub/pkg/ui"

	// plugins
	"github.com/sniperkit/hub/plugin/vcs/provider/github"
	"github.com/sniperkit/hub/plugin/vcs/provider/github/cmd"

	_ "github.com/sniperkit/hub/plugin/vcs/provider/bitbucket"
	_ "github.com/sniperkit/hub/plugin/vcs/provider/bitbucket/cmd"

	_ "github.com/sniperkit/hub/plugin/vcs/provider/gitlab"
	_ "github.com/sniperkit/hub/plugin/vcs/provider/gitlab/cmd"
)

func main() {
	defer github.CaptureCrash()

	err := cmd.CmdRunner.Execute()
	if !err.Ran {
		ui.Errorln(err.Error())
	}
	os.Exit(err.ExitCode)
}
