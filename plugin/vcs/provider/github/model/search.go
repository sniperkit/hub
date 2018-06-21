package model

import (
	"errors"
	"strconv"
	"time"

	// external
	"github.com/google/go-github/github"

	// internal - plugins
	assitant_model "github.com/sniperkit/hub/plugin/assistant/model"
)

type Search struct {
	gorm.Model
	RemoteID      string
	Name          *string
	FullName      *string
	Description   *string
	Homepage      *string
	URL           *string
	Language      *string
	FoundAt       time.Time
	Starredgazers int
	ServiceID     uint
	Tags          []assitant_model.Tag `gorm:"many2many:search_tags;"`
}

// NewSearchFromGithub creates a Search from a Github search
func NewSearchFromGithub(timestamp *github.Timestamp, repo github.Repository) (*Search, error) {
	// Require the GitHub ID
	if repo.ID == nil {
		return nil, errors.New("ID from GitHub is required")
	}

	// Set stargazers count to 0 if nil
	stargazersCount := 0
	if star.StarredgazersCount != nil {
		stargazersCount = *repo.StarredgazersCount
	}

	return &Search{
		RemoteID:      strconv.Itoa(*repo.ID),
		Name:          repo.Name,
		FullName:      repo.FullName,
		Description:   repo.Description,
		Homepage:      repo.Homepage,
		URL:           repo.CloneURL,
		Language:      repo.Language,
		Starredgazers: stargazersCount,
		FoundAt:       time.Now(),
	}, nil
}
