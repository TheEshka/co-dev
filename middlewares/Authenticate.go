package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/common"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if !(len(tokenString) > 7) || strings.ToUpper(tokenString[0:6]) != "BEARER" {
			common.RespondError(w, http.StatusUnauthorized, auth.ErrWrongToken.Error())
			return
		}
		tokenString = tokenString[7:]
		token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, auth.KeyFunc)
		if err != nil {
			common.RespondError(w, http.StatusUnauthorized, auth.ErrWrongToken.Error())
			return
		}

		claims, ok := token.Claims.(*auth.Claims)
		if !ok || !token.Valid {
			common.RespondError(w, http.StatusUnauthorized, auth.ErrWrongToken.Error())
			return
		}
		ctx := auth.SetUserId(r.Context(), claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
