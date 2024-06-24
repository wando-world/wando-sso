package models

type CreateUserRequest struct {
	Nickname     string  `json:"nickname" validate:"required"`
	UserID       string  `json:"userId" validate:"required,min=5,max=99,alphanum"`
	Password     string  `json:"password" validate:"required,min=8,max=99,alphanum"`
	Email        *string `json:"email" validate:"omitempty,email"`
	VerifiedCode string  `json:"verifiedCode" validate:"required,alphanum,min=6,max=10"`
}

type FindSelfResponse struct {
	Nickname string  `json:"nickname"`
	UserID   string  `json:"userId"`
	Email    *string `json:"email"`
}
