package models

type LoginRequest struct {
	UserID       string `json:"userId" validate:"required"`
	Password     string `json:"password" validate:"required"`
	VerifiedCode string `json:"verifiedCode" validate:"required"`
}

type LoginResponse struct {
	ATK string `json:"atk"`
	RTK string `json:"rtk"`
}

type RefreshAtkResponse struct {
	ATK string `json:"atk"`
}
