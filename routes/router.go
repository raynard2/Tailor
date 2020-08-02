package routes

import (
	"Mlops/config"
	"Mlops/controller"
	//"net/http"

	"github.com/labstack/echo/middleware"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
)

var SignedKey = config.GetHmacSignKey()

func New() *echo.Echo {

	e := echo.New()

	e.Use(middleware.Logger()) // Logger
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.GET("/index", controller.Home)

	// api group
	api := e.Group("/v1")
	// jwt config middleware
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    SignedKey,
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/v1/login" || c.Path() == "/v1/register" {
				return true
			}
			return false
		},
	}))
	// api
	api.POST("/register", controller.CreateUser)
	api.POST("/login", controller.Login)

	// admin group
	adminGroup := api.Group("/admin")
	// admin log middleware
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339_nano} ${status} ${latency_human} ${uri} ${remote_ip} ${method}] +\n`,
	}))
	// admin api
	adminGroup.GET("/getusers", controller.GetUsers)
	adminGroup.DELETE("/deleteuser", controller.DeleteUser)
	adminGroup.DELETE("/deleteusers", controller.DeleteUserPermanent)
	adminGroup.PUT("/updateuser", controller.UpdateUser)

	//ml models
	modelGroup := api.Group("/models")
	// ML models routes
	modelGroup.POST("/model", controller.CreateModel)
	modelGroup.GET("/models", controller.Models)
	modelGroup.GET("/models/:model_id", controller.ViewRunningModel)
	modelGroup.PATCH("/models/:model_id", controller.PausingARunningModel)
	modelGroup.DELETE("/model/:model_id", controller.DeleteModel)

	return e
}
