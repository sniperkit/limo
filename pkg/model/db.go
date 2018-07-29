package model

import (
	// external
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// InitDB initializes the database at the specified path
func InitDB(filepath string, verbose bool) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	db.LogMode(verbose)
	db.AutoMigrate(&Service{}, &Star{}, &Tag{})

	return db, nil
}
