package auth

type FindUserReq struct {
	UserId       string `json:"userId" validate:"required"`
	VerifiedCode string `json:"verifiedCode" validate:"required"`
}

type GenerateJWTRes struct {
	ATK string `json:"atk"`
	RTK string `json:"rtk"`
}

type RefreshAtkRes struct {
	ATK string `json:"atk"`
}
