package models

import "github.com/jinzhu/gorm"

type Model struct {
	gorm.Model
	Name            string `json:"name"`
	Path            string `json:"path"`
	TargetVariable  string `json:"target_variable"`
	ModelType       string `json:"model_type"`
	PredictionType  string `json:"prediction_type"`
	Features        string `json:"features"`
	Status          string `json:"status"`
	TrainingDetails string `json:"training_details"`
	ParamVersion    string `json:"param_version"`
	PreProcessID    string `json:"pre_process_id"`
	Threshold       uint   `json:"threshold"`
	ContainerID     string `json:"container_id"`
	UserID          uint   `json:"user_id"`
}
