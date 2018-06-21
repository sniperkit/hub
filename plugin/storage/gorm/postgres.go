// +build with-postgres

package storage_gorm

import (
	"fmt"

	// Use the postgres dialect
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// plugin
	"github.com/sniperkit/hub/pkg/config"
)

type StoragePOSTGRES struct {
	db   *gorm.DB
	conf *config.Storage
}

func (s *StoragePOSTGRES) NewWithConfig(conf *config.Config) (*gorm.DB, error) {

	dbConfig := conf.Config.Storage
	db, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
	if err != nil {
		return nil, err
	}

	if verbose {
		db.LogMode(verbose)
	}

	s.db = db

	return s, nil
}

func (s *StoragePOSTGRES) NewWithDSN(dsn string, verbose bool) (*gorm.DB, error) {

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if verbose {
		db.LogMode(verbose)
	}
	s.db = db

	return s, nil
}
