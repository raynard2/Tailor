package routes



import (

	"Mlops/controller"
	"Mlops/db"
	"Mlops/lib"
	"Mlops/model"
	"github.com/dgrijalva/jwt-go"
	"log"


	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"

)







func createModel(c echo.Context) error {
	db := db.Manager()
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
	db := db.Manager()
	db,_  = gorm.Open("sqlite3", "./test.db")
	defer db.Close()
	var users []model.User
	db.Find(&users)
	return c.JSONPretty(200, users, "")
}




func New () *echo.Echo {

	e := echo.New()
	api := e.Group("/v1")


	e.GET("/token", controller.Home)
	e.POST("/login", controller.Login)
	e.POST("/create", createModel)
	api.GET("/getusers", getModel)


	return e
}