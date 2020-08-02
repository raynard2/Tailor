package models

import (
	"github.com/jinzhu/gorm"
)

// Stats Struct
type Stats struct {
	gorm.Model
	Type string `json:"opt"`
}
