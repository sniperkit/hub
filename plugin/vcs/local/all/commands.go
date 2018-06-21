package vcs_abstract

import (
	"io/ioutil"

	"github.com/Masterminds/vcs"
)

func CloneRepository(remote, prefixPath, dirName string) (*vcs.Repo, error) {
	local, err := ioutil.TempDir(prefixPath, dirName)
	if err != nil {
		return nil, err
	}
	return vcs.NewRepo(remote, local)
}
