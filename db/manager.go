package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func Initialize(connection *gorm.DB) {
	db = connection

	migrations()
}

func migrations() {
	db.AutoMigrate(&Blueprint{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Rating{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&BlueprintTag{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Revision{})
}
