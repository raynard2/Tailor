package http

import (
	"time"
)

type Model struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Path            string    `json:"path"`
	TargetVariable  string    `json:"target_variable"`
	ModelType       string    `json:"model_type"`
	PredictionType  string    `json:"prediction_type"`
	Features        string    `json:"features"`
	Status          string    `json:"status"`
	TrainingDetails string    `json:"training_details"`
	ParamVersion    string    `json:"param_version"`
	PreProcessID    string    `json:"pre_process_id"`
	UserID          uint      `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
