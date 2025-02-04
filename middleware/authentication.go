package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ARF-DEV/image-processing-api/utils/httputils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if !strings.Contains(authorization, "Bearer ") {
			fmt.Println(authorization)
			httputils.SendResponse(w, httputils.ErrUnauthorized.Error(), nil, nil, httputils.ErrUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authorization, "Bearer ")
		claims := jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(viper.GetString("SECRET_KEY")), nil
		})

		if err != nil {
			httputils.SendResponse(w, err.Error(), nil, nil, httputils.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
