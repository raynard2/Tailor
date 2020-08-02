package controller

import (
	"Mlops/db"
	"Mlops/lib"
	http2 "Mlops/lib/http"
	userLib "Mlops/lib/user"
	"Mlops/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

var validate = validator.New()

func Login(c echo.Context) error {

	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic(lib.DbError)
	}
	defer Db.Close()

	params := new(userLib.LoginParams)

	err = validate.Struct(params)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return http2.ValidationErrorResponse(c, errors)
	}

	if err := c.Bind(params); err != nil {
		return InternalError(c, err.Error())

	}

	user := new(models.User)

	Db.Where("email =?", params.Email).First(&user)
	if user.ID == 0 {
		return BadRequestResponse(c, lib.InvalidUser)
	}

	//auth the params with user
	auth := lib.CompareHashAndPassword(user.Password, params.Password)
	if !auth {
		return BadRequestResponse(c, lib.WrongPassword)
	}

	//generate token
	token, err := userLib.GenerateToken(user)
	if err != nil {
		return InternalError(c, err.Error())
	}

	//parse data into response structure
	response := userLib.LoginResponse{
		Success: true,
		User: userLib.UserResponse{
			FullName: user.FullName,
			Email:    user.Email,
			Active:   true,
			Channel:  user.Email,
		},
		Token:   token,
		IsAdmin: user.IsAdmin,
	}

	return DataResponse(c, http.StatusOK, response)
}

func CreateUser(c echo.Context) error {

	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic(lib.DbError)
	}
	defer Db.Close()

	params := new(userLib.CreateUserParams)

	if err := c.Bind(params); err != nil {
		return InternalError(c, lib.ErrorBinding)
	}

	if params.Password !=	params.ConfirmPassword {
		return BadRequestResponse(c, "Input Password Error")
	}

	user := new(models.User)

	exist := Db.Where("email= ?", params.Email).Find(&user).RecordNotFound()
	if exist == false {
		return BadRequestResponse(c, lib.AccountExists)
	}

	user.FullName = params.FullName
	user.Email = params.Email
	user.Password = lib.CreateHashFromPassword(params.Password)
	user.IsAdmin = params.IsAdmin
	user.Active = true

	Db.Save(&user)

	response := userLib.CreateUserResonse{
		User:    userLib.UserResponse{
			FullName: user.FullName,
			Email: user.Email,
			Channel: user.Email,
			Active: user.Active,
		},
		Success: true,
		IsAdmin: params.IsAdmin,
	}
	exist = Db.Where("email= ?", user.Email).Find(&user).RecordNotFound()
	if exist == true {
		return InternalError(c, lib.ErrorCreatingAccount)
	}

	return DataResponse(c, http.StatusCreated, response)
}

func GetUser(c echo.Context) (*models.User, bool) {
	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic(lib.DbError)
	}
	defer Db.Close()

	user := new(models.User)
	claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	UserDoesNotExist := Db.Where("id = ?", claims["user_id"]).First(&user).RecordNotFound()
	return user, UserDoesNotExist
}

func GetUsers(c echo.Context) error {
	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic(lib.DbError)
	}
	defer Db.Close()


	admin, notFound := GetUser(c)
	if notFound {
		return BadRequestResponse(c, lib.JwtAccountNotExist)
	}

	// Throws unauthorized error
	exists := Db.Where("email = ?", admin.Email).Find(&admin).RecordNotFound()

	if exists {
		return BadRequestResponse(c, lib.AdminPrivilegeError)
	}

	isAdmin := admin.IsAdmin
	var users []models.User
	if isAdmin == false {
		return BadRequestResponse(c, lib.AdminPrivilegeError)
	}
	Db.Find(&users)
	return DataResponse(c, http.StatusOK, users)
}

func DeleteUser(c echo.Context) error {
	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic(lib.DbError)
	}
	defer Db.Close()
	email := c.QueryParam("email")

	user := new(models.User)


	if err := Db.Delete(&user, "email= ?", email); err.RowsAffected < 1 {
		log.Println("invalid user")
		return c.JSONPretty(http.StatusBadRequest, err, "")
	}

	return DataResponse(c, http.StatusAccepted, user)
}

func DeleteUserPermanent(c echo.Context) error {
	Db := db.Manager()

	Db, _ = gorm.Open("postgres", "name=raynardomongbale password=raynard dbname=mlops dbssl=disable")

	defer Db.Close()

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusInternalServerError, "error binding params")
	}

	_ = Db.Delete(&user)


	return DataResponse(c, http.StatusAccepted, user)
}

func UpdateUser(c echo.Context) error {

	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic(lib.DbError)
	}
	defer Db.Close()

	params := new(userLib.CreateUserParams)

	user := new(models.User)

	if err := c.Bind(params); err != nil {
		return c.String(http.StatusInternalServerError, lib.ErrorBinding)
	}

	exist := Db.Where("email= ?", params.Email).Find(&user).RecordNotFound()
	if exist  {
		return BadRequestResponse(c, lib.AccountNotExist)
	}

	Db.Model(&user).Update(map[string]interface{}{
		user.FullName : params.FullName,
		user.Email : params.Email,
		user.Password : lib.CreateHashFromPassword(params.Password),

	})
	Db.Save(&user)

	response := userLib.CreateUserResonse{
		User:    userLib.UserResponse{
			FullName: user.FullName,
			Email: user.Email,
			Channel: user.Email,
			Active: user.Active,
		},
		Success: true,
		IsAdmin: params.IsAdmin,
	}


	return DataResponse(c, http.StatusOK, response)
}


