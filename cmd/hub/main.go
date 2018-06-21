// +build go1.8

package main

import (
	"os"

	"github.com/sniperkit/hub/pkg/commands"
	"github.com/sniperkit/hub/pkg/ui"

	"github.com/sniperkit/hub/plugin/vcs/provider/github"
)

func main() {
	defer github.CaptureCrash()

	err := commands.CmdRunner.Execute()
	if !err.Ran {
		ui.Errorln(err.Error())
	}
	os.Exit(err.ExitCode)
}
