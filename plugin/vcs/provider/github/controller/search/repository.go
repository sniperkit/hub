package ghs

import (
	"sync"

	"context"
	"github.com/google/go-github/github"
)

type Repo struct {
	search     *Search
	maxItem    int
	printCount int
}

func NewRepo(s *Search) *Repo {
	return &Repo{s, 0, 0}
}

func (r *Repo) Search(c context.Context) (<-chan []github.Repository, <-chan error) {
	var wg sync.WaitGroup
	reposBuff := make(chan []github.Repository, 1)
	errChan := make(chan error, 1)

	// 1st search

	repos, lastPage, max, err := r.search.First(c)
	if err != nil {
		Debug("Error First()\n")
		errChan <- err
		return reposBuff, errChan
	}
	r.maxItem = max
	// notify main thread of first search result
	reposBuff <- repos

	// 2nd - 10th search
	go func() {
		for page := 2; page < lastPage+1; page++ {
			Debug("sub thread start %d\n", page)
			wg.Add(1)
			go func(p int) {
				// notify main thread of 2nd - 10th search result
				rs, err := r.search.Exec(c, p)
				if err != nil {
					Debug("sub thread error %d\n", p)
					errChan <- err
				}
				reposBuff <- rs
				wg.Done()
				Debug("sub thread end %d\n", p)
			}(page)
		}
		Debug("sub thread wait...\n")
		wg.Wait()
		Debug("sub thread wakeup!!\n")
		close(reposBuff)
	}()

	Debug("main thread return\n")

	return reposBuff, errChan
}

func (r *Repo) Print(repos []github.Repository) (bool, int) {
	Debug("repos length %d\n", len(repos))
	repoNameMaxLen := 0
	for _, repo := range repos {
		repoNamelen := len(*repo.FullName)
		if repoNamelen > repoNameMaxLen {
			repoNameMaxLen = repoNamelen
		}
	}
	printLine := func(repo *github.Repository) {
		if repo.FullName != nil {
			Printf("%v", *repo.FullName)
		}
		Printf("    ")
		paddingLen := repoNameMaxLen - len(*repo.FullName)
		for i := 0; i < paddingLen; i++ {
			Printf(" ")
		}
		if repo.Description != nil {
			Printf("%v", *repo.Description)
		}
		Printf("\n")
	}
	for _, repo := range repos {
		printLine(&repo)
		r.printCount++
		Debug("printCount %d, r.maxItem %d\n", r.printCount, r.maxItem)
		if r.printCount >= r.maxItem {
			return true, r.printCount
		}
	}
	return false, r.printCount
}
