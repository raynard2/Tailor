package routes

import (
	"Mlops/controller"
	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()
	e.GET("/", controller.Home)

	return e
}
