package routes

import (
	"Mlops/controller"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"

)

var db *gorm.DB
type User struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
}



func user(c echo.Context) error {
	db, err := gorm.Open("sqlite3", "./database/database.db")
	if err != nil {
		panic("error opening db")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	return c.JSONPretty(200, users, "")
}





func New() *echo.Echo {
	e := echo.New()
	//DBInit()
	e.GET("/", controller.Home)
	dbGroup := e.Group("/db")
	dbGroup.Use(DatabaseInit)
	dbGroup.GET("/user", user)

	return e
}
