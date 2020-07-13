package controller

import "github.com/labstack/echo"

func Home(c echo.Context) error {
	return c.JSON(200, "we are live")
}
