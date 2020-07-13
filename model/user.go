package model

import "github.com/jinzhu/gorm"
//user model
type User struct {
	gorm.Model

	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	Password  string `json:"password"`
	Active    bool   `json:"active"`
}

