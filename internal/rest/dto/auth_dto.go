package dto

type LoginRequest struct {
	UserID       string `json:"userId" validate:"required"`
	Password     string `json:"password" validate:"required"`
	VerifiedCode string `json:"verifiedCode" validate:"required"`
}
