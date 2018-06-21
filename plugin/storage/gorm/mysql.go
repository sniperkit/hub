// +build with-mysql

package storage_gorm

import (
	"fmt"

	// Use the mysql dialect
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	// plugin
	"github.com/sniperkit/hub/pkg/config"
)

type StorageMYSQL struct {
	db   *gorm.DB
	conf *config.Storage
}

func (s *StorageMYSQL) NewWithConfig(conf *config.Config) (*gorm.DB, error) {

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

func (s *StorageMYSQL) NewWithDSN(dsn string, verbose bool) (*gorm.DB, error) {

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

/*
	dbConfig := config.Config.DB
	if config.Config.DB.Adapter == "mysql" {
		DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		// DB = DB.Set("gorm:table_options", "CHARSET=utf8")
	} else if config.Config.DB.Adapter == "postgres" {
		DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name))
	} else if config.Config.DB.Adapter == "sqlite" {
		DB, err = gorm.Open("sqlite3", fmt.Sprintf("%v/%v", os.TempDir(), dbConfig.Name))
	} else {
		panic(errors.New("not supported database adapter"))
	}
*/
