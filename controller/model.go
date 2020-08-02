package controller

import (
	"Mlops/config"
	"Mlops/db"
	"Mlops/lib"
	http_response "Mlops/lib/http/model"
	models "Mlops/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

const ModelCreated string = "Model created successfully."
const ModelNotExist string = "Model does not exist."
const ModelDeleted string = "Model deleted successfully."

var configuration = config.GetConfig()

func CreateModel(c echo.Context) error {
	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic("error opening Db")
	}
	defer Db.Close()

	params := new(CreateModelParams)

	if err := c.Bind(params); err != nil {
		return c.String(404, "error binding")
	}
	userContext := c.Get("user").(*jwt.Token)
	claims := userContext.Claims.(jwt.MapClaims)
	email := claims["email"].(string)


	user := new(models.User)

	if recordNotFound := Db.Where("email= ?", email).Find(&user).RecordNotFound(); recordNotFound {
		return BadRequestResponse(c, lib.InvalidUser)
	}

	//contID, err := CreateNewContainer(params.ModelName,user.FullName)
	//
	//if err != nil {
	//	return c.String(404, "error creating container")
	//}
	newModel := models.Model{

		Name:            params.ModelName,
		Path:            "",
		TargetVariable:  params.TargetVariable,
		ModelType:       params.ModelType,
		PredictionType:  params.PredictionType,
		Features:        "",
		Status:          "STARTED",
		TrainingDetails: "",
		ParamVersion:    params.ModelName,
		PreProcessID:    params.ModelName,
		Threshold:       0,
		ContainerID:     string(user.Model.ID),
		UserID:          user.Model.ID,
	}

	Db.Save(&newModel)

	return DataResponse(c, http.StatusCreated, newModel)
}

func Models(c echo.Context) error {
	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic("error opening Db")
	}
	defer Db.Close()

	//get token claims
	userModel, notFound :=GetUser(c)
	if notFound	{
		return BadRequestResponse(c, lib.JwtAccountNotExist)
	}
	//userContext := c.Get("user").(*jwt.Token)
	//claims := userContext.Claims.(jwt.MapClaims)
	//// get email from claims
	//email := claims["email"].(string)
	////get instance of gorm table
	//userModel := new(models.User)
	//find the user model
	exists := Db.Where("email= ?", userModel.Email).Find(&userModel).RecordNotFound()
	if exists {
		return BadRequestResponse(c, lib.AccountNotExist)
	}

	pageNum := c.QueryParam("page")
	idP, _ := strconv.Atoi(pageNum)
	page := int64(idP)
	offset := int((page - 1) * lib.Limit)
	var model []models.Model
	var totalCount int64
	Db.Where("user_id = ?", userModel.Model.ID).Limit(lib.Limit).Offset(offset).Order("created_at desc").Find(&model)
	//Get Total
	Db.Model(&model).Where("user_id = ?", userModel.Model.ID).Count(&totalCount)

	response := make([]http_response.Model, 0)

	for _, value := range model {

		modelDetails := new(http_response.Model)
		modelDetails.ID = value.ID
		modelDetails.Name = value.Name
		modelDetails.Path = value.Path
		modelDetails.TargetVariable = value.TargetVariable
		modelDetails.ModelType = value.ModelType
		modelDetails.PredictionType = value.PredictionType
		modelDetails.Features = value.Features
		modelDetails.Status = value.Status
		modelDetails.TrainingDetails = value.TrainingDetails
		modelDetails.ParamVersion = value.ParamVersion
		modelDetails.PreProcessID = value.PreProcessID
		modelDetails.CreatedAt = value.CreatedAt
		modelDetails.UpdatedAt = value.UpdatedAt
		response = append(response, *modelDetails)

	}
	//return PaginatedDataResponse(c, http.StatusOK, response, totalCount, page)

	return PaginatedDataResponse(c, http.StatusOK, response, totalCount, page)
}

func ViewRunningModel(c echo.Context) error {
	Db := db.Manager()
	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic("error opening Db")
	}
	var model models.Model

	//user := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	//email := claims["email"].(string)
	user, notFound := GetUser(c)
	if notFound	{
		return BadRequestResponse(c, lib.JwtAccountNotExist)
	}
	email := user.Email
	model_id := c.Param("model_id")

	userModel := new(models.User)
	exists := Db.Where("email = ?", email).Find(&userModel).RecordNotFound()
	if exists {
		return BadRequestResponse(c, lib.AccountNotExist)
	}

	Db.Where("user_id = ? AND _id = ?", userModel.Model.ID, model_id).Find(&model)

	validModel := Db.Where("id = ? AND user_id = ? AND status = 'running'", model_id, userModel.Model.ID).Find(&model).RecordNotFound()
	if validModel {
		return BadRequestResponse(c, ModelNotExist)
	}
	modelDetails := new(http_response.Model)
	modelDetails.ID = model.ID
	modelDetails.Name = model.Name
	modelDetails.Path = model.Path
	modelDetails.TargetVariable = model.TargetVariable
	modelDetails.ModelType = model.ModelType
	modelDetails.PredictionType = model.PredictionType
	modelDetails.Features = model.Features
	modelDetails.Status = model.Status
	modelDetails.TrainingDetails = model.TrainingDetails
	modelDetails.ParamVersion = model.ParamVersion
	modelDetails.PreProcessID = model.PreProcessID
	modelDetails.UserID = model.UserID
	modelDetails.CreatedAt = model.CreatedAt
	modelDetails.UpdatedAt = model.UpdatedAt

	return DataResponse(c, http.StatusOK, modelDetails)
}

func PausingARunningModel(c echo.Context) error {
	Db := db.Manager()
	Db, er := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if er != nil {
		panic("error opening Db")
	}

	var model models.Model
	model_id := c.Param("model_id")

	userContext, _ := GetUser(c)
	email := userContext.Email

	user := new(models.User)

	// Throws unauthorized error
	exists := Db.Where("email = ?", email).Find(&user).RecordNotFound()
	if exists {
		return BadRequestResponse(c, lib.AccountNotExist)
	}

	validModel := Db.Where("id = ? AND user_id = ? AND status = 'running'", model_id, user.Model.ID).Find(&model).RecordNotFound()
	if validModel {
		return BadRequestResponse(c, ModelNotExist)
	}

	err := PauseContainer(model.ContainerID)

	if err != nil {
		return BadRequestResponse(c, err.Error())
	}

	//Update New Status To Database
	model.Status = "stopped"

	Db.Save(&model)

	return DataResponse(c, http.StatusOK, model)
}


func DeleteModel(c echo.Context) error {
	Db := db.Manager()
	//Db, err := gorm.Open("postgres", ".user="+configuration.DB_USERNAME+ "password="+configuration.DB_PASSWORD+ "dbname="+configuration.DB_NAME +"sslmode=disable")

	Db, err := gorm.Open("postgres", ".user=raynardomongbale password=raynard dbname=mlops sslmode=disable")
	if err != nil {
		panic("error opening Db")
	}
	defer Db.Close()

	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(jwt.MapClaims)
	//email := claims["email"].(string)

	user, notFound := GetUser(c)
	if notFound	{
		return BadRequestResponse(c, lib.JwtAccountNotExist)
	}
	email := user.Email

	model_id := c.Param("model_id")
	userModel := new(models.User)

	exists := Db.Where("email = ?", email).Find(&userModel).RecordNotFound()
	if exists {
		return InternalError(c, lib.AccountNotExist)
	}

	model := new(models.Model)
	// Soft Delete Models
	Db.Where("user_id = ? AND _id = ?", user.ID, model_id).Delete(&model)

	return DataResponse(c, http.StatusAccepted, ModelDeleted)
}
