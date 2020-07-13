package db

import (
	"Mlops/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var db *gorm.DB

func DatabaseInit (next echo.HandlerFunc) echo.HandlerFunc {
	db, err := gorm.Open("sqlite3", "./database/database.db")
	if err != nil {
		panic("error opening db")
	}
	defer db.Close()

	db.AutoMigrate(&model.User{})

	return func(c echo.Context) error {

		return next(c)
	}
}


