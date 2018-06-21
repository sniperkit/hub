package model

import (
	"time"
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
	Starredgazers int
	StarredredAt  time.Time
	ServiceID     uint
	Tags          []Tag `gorm:"many2many:search_tags;"`
}
