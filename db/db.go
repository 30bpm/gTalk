package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "groomer:groomer5959@tcp(db-groomer-talk.c89btdfo8aq1.ap-northeast-2.rds.amazonaws.com:3306)/groomer_talk?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Migrator().DropTable(&Message{})
	db.Migrator().CreateTable(&Message{})
	DB = db
}
