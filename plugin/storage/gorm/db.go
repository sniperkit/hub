package storage_gorm

import (
	"errors"

	"github.com/jinzhu/gorm"

	config "github.com/sniperkit/hub/plugin/config"
)

var (
	ErrUnsupportedStorageAdapter = errors.New("Not supported storage adapter")
)

type Storage interface {
	NewWithConfig(conf *config.Config) (*gorm.DB, error)
	NewWithDSN(dsn string, verbose bool) (*gorm.DB, error)
}

// InitDB initializes the database with a dsn address...
func InitDB(dsn string, autoMigrate, verbose bool) (*gorm.DB, error) {

	var err error
	var storage *gorm.DB

	switch storeBackend {
	case "mysql":
		backend := &StorageMSSQL{}
		storage, err = backend.NewWithDSN(dsn, debug, verbose)

	case "postgres":
		backend := &StoragePOSTGRES{}
		storage, err = backend.NewWithDSN(dsn, debug, verbose)

	case "mssql":
		backend := &StorageMSSQL{}
		storage, err = backend.NewWithDSN(dsn, debug, verbose)

	case "sqlite", "sqlite3":
		backend := &StorageSQLITE3{}
		storage, err = backend.NewWithDSN(dsn, debug, verbose)

	default:
		panic(ErrUnsupportedStorageAdapter)
	}

	/*
		if autoMigrate {
			db.AutoMigrate(&Service{}, &Star{}, &Tag{})
		}
	*/

	return storage, err
}

// NewWithConfig initializes the database with a config struct...
func NewWithConfig(conf *config.Storage) (*gorm.DB, error) {

	var err error
	var storage *gorm.DB

	switch storeBackend {
	case "mysql":
		backend := &StorageMYSQL{}
		storage, err = backend.NewWithConfig(conf)

	case "postgres":
		backend := &StoragePOSTGRES{}
		storage, err = backend.NewWithConfig(conf)

	case "mssql":
		backend := &StorageMSSQL{}
		storage, err = backend.NewWithConfig(conf)

	case "sqlite", "sqlite3":
		backend := &StorageSQLITE3{}
		storage, err = backend.NewWithConfig(conf)

	default:
		panic(ErrUnsupportedStorageAdapter)
	}

	/*
		if conf.AutoMigrate {
			db.AutoMigrate(&Service{}, &Star{}, &Tag{})
		}
	*/

	return storage, err
}
