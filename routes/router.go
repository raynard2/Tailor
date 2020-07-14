package routes



import (
	"Mlops/lib"
	"Mlops/lib/user"
	"Mlops/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"

)
var db *gorm.DB
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


func GenerateToken (user *model.User) (*jwt.Token, string,error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"user_id": user.ID,
		"issued_at": time.Now(),
		"expire_at": time.Now().Add(time.Minute * 72).Unix(),
	}
	rawtoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawtoken.SignedString([]byte("AppSecret"))
	if err != nil {
		log.Println("error generating token")
		return  rawtoken,"",err
	}
	log.Println(rawtoken)
	return rawtoken,token,err
}
func CreateCookie (user *model.User) *http.Cookie {
	cookie := http.Cookie{
		Name: "mlops",
		Value: "mlops_id",
		Expires: time.Now().Add(30 * time.Minute),
	}
	return &cookie
}



func Login(c echo.Context) error {
	//db := db.Manager()
	db,_  = gorm.Open("sqlite3", "./test.db")
	defer db.Close()
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
		return c.JSON(http.StatusUnauthorized, "Wrong Password")
	}
	c.JSON(http.StatusUnauthorized, "Correct Password")
	c.SetCookie(CreateCookie(user))
	rawtoken,token, err := GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error generating token for client")
	}
	c.JSON(http.StatusOK, rawtoken)
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
	log.Println(response)
	return c.JSONPretty(200, response, "")
}



func createModel(c echo.Context) error {

	db,_  = gorm.Open("sqlite3", "./test.db")
	defer db.Close()
	params  := new(model.User)
	_ = c.Bind(params)
	password := params.Password
	hashpassword, _ := lib.CreateHashFromPassword(password)
	user := model.User{
		FullName: "Raynard",
		Email:    params.Email,
		Password: hashpassword,
		Active:   true,
		Token:    "",
	}

	if recordExist := db.NewRecord(user); recordExist == false {
		log.Println("user doesnt exist !st")
	}
	db.Create(&user)
	db.Save(&user)
	if recordExist := db.NewRecord(user); recordExist == false {
		log.Println("user doesnt exist 2nd")
	}

	log.Println(user)
	//users := db.Find(&user)
	return c.JSONPretty(200, user, "")
}

func getModel(c echo.Context) error {
	db,_  = gorm.Open("sqlite3", "./test.db")
	defer db.Close()
	var users []model.User
	db.Find(&users)
	return c.JSONPretty(200, users, "")
}
type JwtClaims struct {
	ModelUser string `json:"model_user"`
	jwt.StandardClaims
}

func token(c echo.Context) error {
	claims := jwt.MapClaims{
	"email": "user",
	"user_id": "user.ID",
	"issued_at": time.Now(),
	"expire_at": time.Now().Add(time.Hour * 72).Unix(),
}
	rawtoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawtoken.SignedString([]byte("AppSecret"))
	if err != nil {
		log.Println("error generating token")
		return  err
	}


	return c.JSONPretty(333,token,"")
}

func  Dbinit()   {
	db,_  = gorm.Open("sqlite3", "./test.db")
	defer db.Close()
	db.AutoMigrate(&model.User{})
}


func New () *echo.Echo {

	e := echo.New()
	//	api := e.Group("/v1")

	e.GET("/", getModel)
	e.GET("/token", token)
	e.POST("/create", createModel)
	e.POST("/login", Login)
	e.Start(":8080")

	return e
}