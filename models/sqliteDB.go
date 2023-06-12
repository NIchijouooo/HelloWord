package models

/**
20230605
*/
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "file:./em.db?cache=shared&mode=rwc&_pragma=journal_mode=WAL&_pragma=legacy_alter_table=ON"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}
