package db

import (
	"Mlops/model"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error


func DatabaseInit (){
	db, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		panic("error opening db")
	}
	defer db.Close()

	db.AutoMigrate(&model.User{})
}





func Manager() *gorm.DB {
	return db

}
