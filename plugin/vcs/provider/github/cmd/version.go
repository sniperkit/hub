package cmd

import (
	"github.com/sniperkit/hub/pkg/ui"
	"github.com/sniperkit/hub/pkg/utils"
	"github.com/sniperkit/hub/pkg/version"
)

var cmdVersion = &Command{
	Run:          runVersion,
	Usage:        "version",
	Long:         "Shows git version and hub client version.",
	GitExtension: true,
}

func init() {
	CmdRunner.Use(cmdVersion, "--version")
}

func runVersion(cmd *Command, args *Args) {
	output, err := version.FullVersion()
	if output != "" {
		ui.Println(output)
	}
	utils.Check(err)
	args.NoForward()
}
