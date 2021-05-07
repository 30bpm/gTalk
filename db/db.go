package db

import (
	"fmt"

	"github.com/groomer/gTalk/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var NOTE_DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/groomer_talk?charset=utf8mb4&parseTime=True&loc=Local", config.DB_USER, config.DB_PASSWORD, config.DB_HOST)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	DB = db

	// db.AutoMigrate(&Message{})
	// db.Migrator().DropTable(&Message{})
	// db.Migrator().CreateTable(&Message{})
}
