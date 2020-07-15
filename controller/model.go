package controller

import (
	"Mlops/db"
	"Mlops/lib"
	"Mlops/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
)

func CreateModel(c echo.Context) error {
	db := db.Manager()
	db, _ = gorm.Open("sqlite3", "./database/database.db")
	defer db.Close()
	params := new(CreateModelParams)
	_ = c.Bind(params)
	password := params.Password
	hashpassword, _ := lib.CreateHashFromPassword(password)
	exist := db.Where("email= ?", params.Email).Find(&params).RecordNotFound()
	if exist == false{
		return c.JSON(http.StatusConflict, "user exist")
	}

	user := model.Model{


		Name:            params.ModelName,
		Path:            "",
		TargetVariable:  params.TargetVariable,
		ModelType:       params.ModelType,
		PredictionType:  params.PredictionType,
		Features:        "",
		Status:          "STARTED",
		TrainingDetails: "",
		//ContainerID:     contID,
		ParamVersion:    params.ModelName,
		PreProcessID:    params.ModelName,
		Threshold:       params.Threshold,
		//UserID:          user.Model.ID,


	}

	db.Save(&user)
	exist = db.Where("email= ?", params.Email).Find(&params).RecordNotFound()
	if exist == true{
		return c.JSON(http.StatusNotModified, "Error saving user details")
	}

	return c.JSONPretty(http.StatusCreated, user, "")
}

