package controller

import (
	"Mlops/db"
	"Mlops/lib"
	userLib "Mlops/lib/user"
	"Mlops/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func Login(c echo.Context) error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./database/database.db")
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
	c.JSON(http.StatusAccepted, "Correct Password")
	//create cookie
	c.SetCookie(userLib.CreateCookie(user))
	//generate token
	rawtoken, token, err := userLib.GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error generating token for client")
	}
	c.JSON(http.StatusOK, rawtoken)
	//parse data into response structure
	response := userLib.LoginResponse{
		Success: true,
		User: userLib.UserResponse{
			FullName: user.FullName,
			Email:    user.Email,
			Active:   true,
			Channel:  user.Email,
		},
		Token: token,
		IsAdmin: user.IsAdmin,
	}


	return c.JSONPretty(200, response, "")
}



func CreateUser(c echo.Context) error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./database/database.db")

	defer db.Close()
	params := new(userLib.CreateUserParams)



	if  err := c.Bind(params);err != nil {
		log.Println("error binding params")
		return err
	}

	user := new(model.User)
	exist := db.Where("email= ?", params.Email).Find(&user).RecordNotFound()
	if exist == false{
		return c.JSON(http.StatusConflict, "email exist already")
	}


		user.FullName= params.FullName
		user.Email =   params.Email
		user.Password =  lib.CreateHashFromPassword(params.Password)
		user.IsAdmin = params.IsAdmin



	db.Save(&user)

	exist = db.Where("email= ?", user.Email).Find(&user).RecordNotFound()
	if exist == true{
		return c.JSON(http.StatusNotModified, "Error saving user details")
	}

	return c.JSONPretty(http.StatusCreated, user, "")
}




func GetUsers(c echo.Context) error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./database/database.db")
	defer db.Close()
	var users []model.User
	db.Find(&users)
	return c.JSONPretty(200, users, "")
}

func DeleteUser(c echo.Context)	error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./database/database.db")
	defer db.Close()

	user := new(model.User)

	if err := c.Bind(user); err != nil {
		return c.String(http.StatusInternalServerError, "error binding params")
	}
	err := db.Delete(&user,"email= ?",user.Email)

	if err.RowsAffected < 1 {
		log.Println("invalid user")
		return c.JSONPretty(http.StatusBadRequest, err, "")
	}
	var users []model.User
	db.Find(&users)
	return c.JSONPretty(http.StatusOK, users, "")
}




func DeleteUserPermanent(c echo.Context) error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./database/database.db")
	defer db.Close()

	user := new(model.User)

	_ = db.Delete(&user)

var users []model.User

	db.Find(&users)

	return c.JSONPretty(http.StatusOK, users, "")
}




func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["is_admin"].(bool)
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
