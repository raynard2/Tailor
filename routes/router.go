package routes

import (
	"Mlops/config"
	"Mlops/controller"
	"Mlops/lib"
	//"net/http"

	"github.com/labstack/echo/middleware"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
)

var SignedKey = config.GetHmacSignKey()


func New() *echo.Echo {

	e := echo.New()
	api := e.Group("/v1")
	authGroup := api.Group("/oauth")
	adminGroup := api.Group("/admin")
	cookieGroup := api.Group("/cookie")
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

	api.GET("/index", controller.Home)
	api.POST("/register", controller.CreateUser)

	//authGroup
	authGroup.POST("/login", controller.Login)

	//admin login required
	adminGroup.GET("/getusers", controller.GetUsers)
	adminGroup.DELETE("/deleteuser", controller.DeleteUser)
	adminGroup.DELETE("/deleteusers", controller.DeleteUserPermanent)

	return e
}
