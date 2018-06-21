// +build go1.8

package main

import (
	"os"

	// ocre
	"github.com/sniperkit/hub/pkg/commands"
	"github.com/sniperkit/hub/pkg/ui"

	// plugins
	"github.com/sniperkit/hub/plugin/vcs/provider/bitbucket"
	. "github.com/sniperkit/hub/plugin/vcs/provider/github"
	. "github.com/sniperkit/hub/plugin/vcs/provider/gitlab"
	// "github.com/sniperkit/hub/plugin/vcs/provider/all"
)

func main() {
	defer github.CaptureCrash()

	err := commands.CmdRunner.Execute()
	if !err.Ran {
		ui.Errorln(err.Error())
	}
	os.Exit(err.ExitCode)
}
