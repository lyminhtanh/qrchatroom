package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/revel/revel"
	"os"
)

type Product struct {
	gorm.Model
	Code string
	Price uint
}

func Connect() (*gorm.DB, error) {
	dbUri := "D:\\sqlite\\data\\test.db"
	if !revel.DevMode {
		dbUri = os.Getenv("DB_URI")
		if dbUri == "" {
			panic("empty db uri")
		}
	}
	db, err := gorm.Open("sqlite3", dbUri)
	db.LogMode(true)
	return db, err
}



