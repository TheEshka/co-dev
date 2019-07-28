package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/common/errors"
	"github.com/misgorod/co-dev/common"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if !(len(tokenString) > 7) || strings.ToUpper(tokenString[0:6]) != "BEARER" {
			common.RespondError(w, errors.ErrWrongToken)
			return
		}
		tokenString = tokenString[7:]
		token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, auth.KeyFunc)
		if err != nil {
			common.RespondError(w, errors.ErrWrongToken)
			return
		}

		claims, ok := token.Claims.(*auth.Claims)
		if !ok || !token.Valid {
			common.RespondError(w, errors.ErrWrongToken)
			return
		}
		ctx := auth.SetUserID(r.Context(), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
