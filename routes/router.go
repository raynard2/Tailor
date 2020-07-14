package routes

import (
	"Mlops/controller"
	"Mlops/db"
	"Mlops/lib"
	"Mlops/lib/user"
	"Mlops/model"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)
type UserResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Channel  string `json:"channel"`
	Active   bool   `json:"active"`
}

type LoginResponse struct {
	Success bool         `json:"success"`
	User    UserResponse `json:"user"`
	Token   string       `json:"token"`
}

var APP_SECRET string = "APP_SECRET"

func Login(c echo.Context) error {
	db := db.DbManager()
	params := new(user.LoginParams)
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
		c.JSON(http.StatusUnauthorized, "Wrong Password")
	}
	c.SetCookie(CreateCookie(user))
	token, err := GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error generating token for client")
	}

	response := LoginResponse{
		Success: true,
		User:    UserResponse{
			FullName: user.FullName,
			Email: user.Email,
			Active: true,
			Channel: user.Email,
		},
		Token:   token,
	}

	return c.JSONPretty(200, response, "")
}
func GenerateToken (user *model.User) (string,error) {
	now := time.Now()
	claims :=	jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"issued_at": now.Unix(),
		"expire_at": now.Add(time.Hour * 72).Unix(),
	}
	 	rawtoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		token,err := rawtoken.SignedString(APP_SECRET)
	return token,err
}
func CreateCookie (user *model.User) *http.Cookie {
	cookie := http.Cookie{
		Name: "mlops",
		Value: "mlops_id",
		Expires: time.Now().Add(30 * time.Minute),
	}
	return &cookie
}




func New() *echo.Echo {
	e := echo.New()
	//DBInit()
	e.GET("/", controller.Home)
	dbGroup := e.Group("/db")
	dbGroup.Use(db.DatabaseInit)
	dbGroup.POST("/user/", Login)

	return e
}
