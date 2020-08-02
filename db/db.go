package db

import (
	"Mlops/model"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func DatabaseInit() {
	db, err = gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic("error opening db")
	}
	defer db.Close()

	db.AutoMigrate(&models.Stats{}, &models.User{}, &models.Model{})
}

func Manager() *gorm.DB {
	return db

}
