package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/avatar-api-gin"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Database Auto Migrate
	database.AutoMigrate(&Avatar{})
	DB = database
}
