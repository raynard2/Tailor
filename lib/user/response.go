package user

type UserResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Channel  string `json:"channel"`
	Active   bool   `json:"active"`
}




type LoginResponse struct {
	User    UserResponse `json:"user"`
	Token   string       `json:"token"`
	Success bool         `json:"success"`
}
