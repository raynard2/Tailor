package http

import "time"

type User struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	RequestID uint      `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
