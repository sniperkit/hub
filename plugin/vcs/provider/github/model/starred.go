package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	// external
	"github.com/blevesearch/bleve"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/skratchdot/open-golang/open"
	// internal
	// internal - plugins
)

// Starred represents a starred repository
type Starred struct {
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
	Tags          []Tag `gorm:"many2many:star_tags;"`
}

// StarredResult wraps a star and an error
type StarredResult struct {
	Starred *Starred
	Error   error
}

// NewStarredFromGithub creates a Starred from a Github star
func NewStarredFromGithub(timestamp *github.Timestamp, star github.Repository) (*Starred, error) {
	// Require the GitHub ID
	if star.ID == nil {
		return nil, errors.New("ID from GitHub is required")
	}

	// Set stargazers count to 0 if nil
	stargazersCount := 0
	if star.StarredgazersCount != nil {
		stargazersCount = *star.StarredgazersCount
	}

	starredAt := time.Now()
	if timestamp != nil {
		starredAt = timestamp.Time
	}

	return &Starred{
		RemoteID:      strconv.Itoa(*star.ID),
		Name:          star.Name,
		FullName:      star.FullName,
		Description:   star.Description,
		Homepage:      star.Homepage,
		URL:           star.CloneURL,
		Language:      star.Language,
		Starredgazers: stargazersCount,
		StarredredAt:  starredAt,
	}, nil
}

// CreateOrUpdateStarred creates or updates a star and returns true if the star was created (vs updated)
func CreateOrUpdateStarred(db *gorm.DB, star *Starred, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing Starred
	if db.Where("remote_id = ? AND service_id = ?", star.RemoteID, service.ID).First(&existing).RecordNotFound() {
		star.ServiceID = service.ID
		err := db.Create(star).Error
		return err == nil, err
	}
	star.ID = existing.ID
	star.ServiceID = service.ID
	star.CreatedAt = existing.CreatedAt
	return false, db.Save(star).Error
}

// FindStarredByID finds a star by ID
func FindStarredByID(db *gorm.DB, ID uint) (*Starred, error) {
	var star Starred
	if db.First(&star, ID).RecordNotFound() {
		return nil, fmt.Errorf("star '%d' not found", ID)
	}
	return &star, db.Error
}

// FindStarreds finds all stars
func FindStarreds(db *gorm.DB, match string) ([]Starred, error) {
	var stars []Starred
	if match != "" {
		db.Where("full_name LIKE ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match))).Order("full_name").Find(&stars)
	} else {
		db.Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FindUntaggedStarreds finds stars without any tags
func FindUntaggedStarreds(db *gorm.DB, match string) ([]Starred, error) {
	var stars []Starred
	if match != "" {
		db.Raw(`
			SELECT *
			FROM STARS S
			WHERE S.DELETED_AT IS NULL
			AND S.FULL_NAME LIKE ?
			AND S.ID NOT IN (
				SELECT STAR_ID
				FROM STAR_TAGS
			) ORDER BY S.FULL_NAME`,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
	} else {
		db.Raw(`
			SELECT *
			FROM STARS S
			WHERE S.DELETED_AT IS NULL
			AND S.ID NOT IN (
				SELECT STAR_ID
				FROM STAR_TAGS
			) ORDER BY S.FULL_NAME`).Scan(&stars)
	}
	return stars, db.Error
}

// FindStarredsByLanguageAndOrTag finds stars with the specified language and/or the specified tag
func FindStarredsByLanguageAndOrTag(db *gorm.DB, match string, language string, tagName string, union bool) ([]Starred, error) {
	operator := "AND"
	if union {
		operator = "OR"
	}

	var stars []Starred
	if match != "" {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM STARS S, TAGS T 
			INNER JOIN STAR_TAGS ST ON S.ID = ST.STAR_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(S.FULL_NAME) LIKE ? 
			AND (LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ?) 
			GROUP BY ST.STAR_ID 
			ORDER BY S.FULL_NAME`, operator),
			fmt.Sprintf("%%%s%%", strings.ToLower(match)),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&stars)
	} else {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM STARS S, TAGS T 
			INNER JOIN STAR_TAGS ST ON S.ID = ST.STAR_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ? 
			GROUP BY ST.STAR_ID 
			ORDER BY S.FULL_NAME`, operator),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&stars)
	}
	return stars, db.Error
}

// FindStarredsByLanguage finds stars with the specified language
func FindStarredsByLanguage(db *gorm.DB, match string, language string) ([]Starred, error) {
	var stars []Starred
	if match != "" {
		db.Where("full_name LIKE ? AND lower(language) = ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match)),
			strings.ToLower(language)).Order("full_name").Find(&stars)
	} else {
		db.Where("lower(language) = ?",
			strings.ToLower(language)).Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FuzzyFindStarredsByName finds stars with approximate matching for full name and name
func FuzzyFindStarredsByName(db *gorm.DB, name string) ([]Starred, error) {
	// Try each of these, and as soon as we hit, return
	// 1. Exact match full name
	// 2. Exact match name
	// 3. Case-insensitive full name
	// 4. Case-insensitive name
	// 5. Case-insensitive like full name
	// 6. Case-insensitive like name
	var stars []Starred
	db.Where("full_name = ?", name).Order("full_name").Find(&stars)
	if len(stars) == 0 {
		db.Where("name = ?", name).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("lower(full_name) = ?", strings.ToLower(name)).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("lower(name) = ?", strings.ToLower(name)).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("full_name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FindPrunableStarreds finds all stars that weren't updated during the last successful update
func FindPrunableStarreds(db *gorm.DB, service *Service) ([]Starred, error) {
	var stars []Starred
	db.Where("service_id = ? AND updated_at < ?", service.ID, service.LastSuccess).Order("full_name").Find(&stars)
	return stars, db.Error
}

// FindLanguages finds all languages
func FindLanguages(db *gorm.DB) ([]string, error) {
	var languages []string
	db.Table("stars").Order("language").Pluck("distinct(language)", &languages)
	return languages, db.Error
}

// AddTag adds a tag to a star
func (star *Starred) AddTag(db *gorm.DB, tag *Tag) error {
	star.Tags = append(star.Tags, *tag)
	return db.Save(star).Error
}

// LoadTags loads the tags for a star
func (star *Starred) LoadTags(db *gorm.DB) error {
	// Make sure star exists in database, or we will panic
	var existing Starred
	if db.Where("id = ?", star.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("star '%d' not found", star.ID)
	}
	return db.Model(star).Association("Tags").Find(&star.Tags).Error
}

// RemoveAllTags removes all tags for a star
func (star *Starred) RemoveAllTags(db *gorm.DB) error {
	return db.Model(star).Association("Tags").Clear().Error
}

// RemoveTag removes a tag from a star
func (star *Starred) RemoveTag(db *gorm.DB, tag *Tag) error {
	return db.Model(star).Association("Tags").Delete(tag).Error
}

// HasTag returns whether a star has a tag. Note that you must call LoadTags first -- no reason to incur a database call each time
func (star *Starred) HasTag(tag *Tag) bool {
	if len(star.Tags) > 0 {
		for _, t := range star.Tags {
			if t.Name == tag.Name {
				return true
			}
		}
	}
	return false
}

// Index adds the star to the index
func (star *Starred) Index(index bleve.Index, db *gorm.DB) error {
	if err := star.LoadTags(db); err != nil {
		return err
	}
	return index.Index(fmt.Sprintf("%d", star.ID), star)
}

// OpenInBrowser opens the star in the browser
func (star *Starred) OpenInBrowser(preferHomepage bool) error {
	var URL string
	if preferHomepage && star.Homepage != nil && *star.Homepage != "" {
		URL = *star.Homepage
	} else if star.URL != nil && *star.URL != "" {
		URL = *star.URL
	} else {
		if star.Name != nil {
			return fmt.Errorf("no URL for star '%s'", *star.Name)
		}
		return errors.New("no URL for star")
	}
	return open.Starredt(URL)
}

// Delete soft-deletes a star
func (star *Starred) Delete(db *gorm.DB) error {
	return db.Delete(&star).Error
}
