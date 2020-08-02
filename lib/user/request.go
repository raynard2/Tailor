package user

type LoginParams struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type CreateUserParams struct {
	Password string	`json:"password" validate:"required"`
	ConfirmPassword string	`json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	IsAdmin  bool   `json:"is_admin"`
	FullName string `json:"full_name" validate:"required"`
}

type SetPasswordParams struct {
	Hash     string `json:"hash" validate:"required"`
	Password string `json:"password" validate:"required"`
}
