package routes

import (
	"Mlops/config"
	"Mlops/controller"
	"Mlops/db"
	"Mlops/lib"
	"net/http"

	"Mlops/model"

	"github.com/labstack/echo/middleware"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
)

var SignedKey = config.GetHmacSignKey()

func getModel(c echo.Context) error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./database/database.db")
	defer db.Close()
	var users []model.User
	db.Find(&users)
	return c.JSONPretty(200, users, "")
}

func New() *echo.Echo {

	e := echo.New()
	api := e.Group("/v1")
	authGroup := api.Group("/oauth")
	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")
	//middlewares
	// admin log middleware
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339_nano} ${status} ${latency_human} ${uri} ${remote_ip} ${method}] +\n`,
	}))
	// jwt config middleware
	authGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    SignedKey,
	}))
	// cookie middleware
	cookieGroup.Use(lib.CheckCookie)

	cookieGroup.GET("/index", controller.Home)
	authGroup.POST("/login", controller.Login)

	e.POST("/create", controller.CreateModel)
	adminGroup.GET("/getusers", getModel)

	return e
}
