package controller



// Update With Fields For Model Struct
type CreateModelParams struct {
	ModelName      string `json:"model_name"`
	Threshold      uint   `json:"threshold"`
	TargetVariable string `json:"target_variable"`
	ModelType      string `json:"model_type"`
	PredictionType string `json:"prediction_type"`
}

