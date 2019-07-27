package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user-id"`
	jwt.StandardClaims
}

type ctxKey int

const userIDKey = ctxKey(1)

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIDKey).(string)
	return val, ok
}

var key = []byte("MySecretKey")

func KeyFunc(*jwt.Token) (interface{}, error) {
	return key, nil
}

func CreateToken(userID string) (string, error) {
	claims := Claims{
		userID,
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
