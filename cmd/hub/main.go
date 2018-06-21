// +build go1.8

package main

import (
	"os"

	"github.com/sniperkit/hub/pkg/ui"

	"github.com/sniperkit/hub/plugin/vcs/provider/github"
	"github.com/sniperkit/hub/plugin/vcs/provider/github/cmd"
)

func main() {
	defer github.CaptureCrash()

	err := cmd.CmdRunner.Execute()
	if !err.Ran {
		ui.Errorln(err.Error())
	}
	os.Exit(err.ExitCode)
}
