package controller

import (
	"github.com/labstack/echo"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		//user := c.Get("user").(*jwt.Token)
		//claims := user.Claims.(jwt.MapClaims)
		//fmt.Println("admin middleware")
		//isAdmin := claims["is_admin"].(bool)
		//if isAdmin == false {
		//	fmt.Println("admin not authorized")
		////	return echo.ErrUnauthorized
		//}
		return next(c)
	}
}
