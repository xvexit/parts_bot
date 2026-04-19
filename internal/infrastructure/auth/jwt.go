package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJWTManager() *JWTManager {
	return &JWTManager{
		accessSecret:  []byte("access-secret"),
		refreshSecret: []byte("refresh-secret"),
	}
}

type claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

var secret = []byte("super-secret-key")


func (j *JWTManager) GenerateAccessToken(userID int64) (string, error) {
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.accessSecret)
}

func (j *JWTManager) GenerateRefreshToken(userID int64) (string, error) {
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.refreshSecret)
}

func (j *JWTManager) ParseAccessToken(tokenStr string) (int64, error) {
	return j.parse(tokenStr, j.accessSecret)
}

func (j *JWTManager) ParseRefreshToken(tokenStr string) (int64, error) {
	return j.parse(tokenStr, j.refreshSecret)
}

func (j *JWTManager) parse(tokenStr string, secret []byte) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	c := token.Claims.(*claims)
	return c.UserID, nil
}