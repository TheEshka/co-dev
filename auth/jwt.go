package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId string `json:"user-id"`
	jwt.StandardClaims
}

type ctxKey int

const userIdKey = ctxKey(1)

func SetUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func GetUserId(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIdKey).(string)
	return val, ok
}

var key = []byte("MySecretKey")

func KeyFunc(*jwt.Token) (interface{}, error) {
	return key, nil
}

func CreateToken(userId string) (string, error) {
	claims := Claims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return sign, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, KeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}
