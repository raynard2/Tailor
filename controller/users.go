package controller

import (
	"Mlops/db"
	"Mlops/lib"
	userLib "Mlops/lib/user"
	"Mlops/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func Login(c echo.Context) error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./test.db")
	defer db.Close()
	params := new(userLib.LoginParams)
	if err := c.Bind(params); err != nil {
		return c.String(http.StatusInternalServerError, "error binding params")
	}

	user := new(model.User)

	db.Where("email =?", params.Email).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, "Username Invalid")
	}
	//auth the params with user
	auth := lib.CompareHashAndPassword(user.Password, params.Password)
	if !auth {
		return c.JSON(http.StatusUnauthorized, "Wrong Password")
	}
	c.JSON(http.StatusUnauthorized, "Correct Password")
	c.SetCookie(userLib.CreateCookie(user))
	rawtoken, token, err := userLib.GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error generating token for client")
	}
	c.JSON(http.StatusOK, rawtoken)

	response := userLib.LoginResponse{
		Success: true,
		User: userLib.UserResponse{
			FullName: user.FullName,
			Email:    user.Email,
			Active:   true,
			Channel:  user.Email,
		},
		Token: token,
	}
	log.Println(response)
	return c.JSONPretty(200, response, "")
}
