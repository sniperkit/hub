package cmd

import (
	"fmt"
	"regexp"

	"github.com/sniperkit/hub/pkg/utils"

	"github.com/sniperkit/hub/plugin/vcs/provider/github"
)

var cmdCherryPick = &Command{
	Run:          cherryPick,
	GitExtension: true,
	Usage: `
cherry-pick <COMMIT-URL>
cherry-pick <USER>@<SHA>
`,
	Long: `Cherry-pick a commit from a fork on GitHub.

## See also:

hub-am(1), hub(1), git-cherry-pick(1)
`,
}

func init() {
	CmdRunner.Use(cmdCherryPick)
}

func cherryPick(command *Command, args *Args) {
	if args.IndexOfParam("-m") == -1 && args.IndexOfParam("--mainline") == -1 {
		transformCherryPickArgs(args)
	}
}

func transformCherryPickArgs(args *Args) {
	if args.IsParamsEmpty() {
		return
	}

	ref := args.LastParam()
	project, sha, refspec, isPrivate := parseCherryPickProjectAndSha(ref)
	if project != nil {
		args.ReplaceParam(args.IndexOfParam(ref), sha)

		tmpName := "_hub-cherry-pick"
		remoteName := tmpName

		if remote := gitRemoteForProject(project); remote != nil {
			remoteName = remote.Name
		} else {
			args.Before("git", "remote", "add", remoteName, project.GitURL("", "", isPrivate))
		}

		fetchArgs := []string{"git", "fetch", "-q", "--no-tags", remoteName}
		if refspec != "" {
			fetchArgs = append(fetchArgs, refspec)
		}
		args.Before(fetchArgs...)

		if remoteName == tmpName {
			args.Before("git", "remote", "rm", remoteName)
		}
	}
}

func parseCherryPickProjectAndSha(ref string) (project *github.Project, sha, refspec string, isPrivate bool) {
	shaRe := "[a-f0-9]{7,40}"

	var mainProject *github.Project
	localRepo, mainProjectErr := github.LocalRepo()
	if mainProjectErr == nil {
		mainProject, mainProjectErr = localRepo.MainProject()
	}

	url, err := github.ParseURL(ref)
	if err == nil {
		projectPath := url.ProjectPath()

		commitRegex := regexp.MustCompile(fmt.Sprintf("^commit/(%s)", shaRe))
		if matches := commitRegex.FindStringSubmatch(projectPath); len(matches) > 0 {
			sha = matches[1]
			project = url.Project
			return
		}

		pullRegex := regexp.MustCompile(fmt.Sprintf(`^pull/(\d+)/commits/(%s)`, shaRe))
		if matches := pullRegex.FindStringSubmatch(projectPath); len(matches) > 0 {
			pullId := matches[1]
			sha = matches[2]
			utils.Check(mainProjectErr)
			project = mainProject
			refspec = fmt.Sprintf("refs/pull/%s/head", pullId)
			return
		}
	}

	ownerWithShaRegexp := regexp.MustCompile(fmt.Sprintf("^(%s)@(%s)$", OwnerRe, shaRe))
	if matches := ownerWithShaRegexp.FindStringSubmatch(ref); len(matches) > 0 {
		utils.Check(mainProjectErr)
		project = mainProject
		project.Owner = matches[1]
		sha = matches[2]
	}

	return
}
