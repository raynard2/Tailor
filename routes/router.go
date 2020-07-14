package routes



import (
	"Mlops/config"
	"Mlops/controller"
	"Mlops/db"
	"Mlops/lib"
	"Mlops/model"

	"github.com/labstack/echo/middleware"

	"log"


	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"

)


var SignedKey = config.GetHmacSignKey()




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
	authGroup := api.Group("/oauth")
	adminGroup := e.Group("admin")

	//middlewares
	// admin log middleware
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339_nano} ${status} ${latency_human} ${uri} ${remote_ip} ${method}] +\n`,
	}))
	// jwt config middleware
	authGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod:	"HS512",
		SigningKey: SignedKey,
	}))




	e.GET("/index", controller.Home)
	authGroup.POST("/login", controller.Login)

	e.POST("/create", createModel)
	adminGroup.GET("/getusers", getModel)


	return e
}