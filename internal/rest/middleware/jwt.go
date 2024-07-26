package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

type Claims struct {
	Id   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type IJwt interface {
	GenerateATK(id uint, role string) (string, error)
	GenerateRTK(id uint) (string, error)
	ParseToken(tokenString string) (*Claims, error)
	GetAtkSecret() []byte
	GetRtkSecret() []byte
}

type JwtUtils struct {
	atkSecret []byte
	rtkSecret []byte
}

var (
	jwtInstance *JwtUtils
	jwtOnce     sync.Once
)

func NewJwtUtils(atkSecret, rtkSecret string) *JwtUtils {
	jwtOnce.Do(func() {
		jwtInstance = &JwtUtils{atkSecret: []byte(atkSecret), rtkSecret: []byte(rtkSecret)}
	})
	return jwtInstance
}

func (j *JwtUtils) GenerateATK(id uint, role string) (string, error) {
	claims := &Claims{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			Issuer:    "wando",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(j.GetAtkSecret())
}

func (j *JwtUtils) GenerateRTK(id uint) (string, error) {
	claims := &Claims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			Issuer:    "wando",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(j.GetRtkSecret())
}

func (j *JwtUtils) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return j.atkSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (j *JwtUtils) GetAtkSecret() []byte {
	return j.atkSecret
}

func (j *JwtUtils) GetRtkSecret() []byte {
	return j.rtkSecret
}
