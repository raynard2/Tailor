package lib

import (
	"Mlops/db"
	userLib "Mlops/lib/user"
	"Mlops/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func CheckCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		db := db.Manager()
		db, _ = gorm.Open("sqlite3", "./database/database.db")
		defer db.Close()
		params := new(userLib.LoginParams)
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusInternalServerError, "error binding params")
		}
		user := new(model.User)
		log.Println(params)
		db.Where("email =?", params.Email).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, "Username Invalid")
		}
		cookie, err := c.Cookie("mlops_cookie")
		if err != nil {
			return err
		}

		if cookie.Value == string(user.ID) {
			return next(c)
			log.Println(user)

		}

		return c.String(200, "Cookie Match")


	}

}


