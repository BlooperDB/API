package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func Initialize() {
	connection, err := gorm.Open("postgres", "host=postgres user=blooper dbname=blooper sslmode=disable password=ZThnie2mffo2cEAA5E2bytnKW3IgA9vZ")

	if err != nil {
		panic("failed to connect database")
	}

	defer connection.Close()

	db = connection

	migrations()
}

func migrations() {
	db.AutoMigrate(&Blueprint{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Rating{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Version{})
}
