package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TechGG1/chat/server/internal/handler"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
)

var ErrMalformedToken = errors.New("malformed jwt token")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		jwtSecret := os.Getenv("JWT_SECRET")
		if len(authHeader) != 2 {
			handleAuthenticationErr(w, ErrMalformedToken)
			return
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				handleAuthenticationErr(w, err)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				var props = "JWTProps"
				ctx := context.WithValue(r.Context(), props, claims)
				// Access context values in handlers like this
				// props, _ := r.Context().Value(props).(jwt.MapClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				handleAuthenticationErr(w, err)
				return
			}
		}
	})
}

func handleAuthenticationErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := handler.ErrorResponse{Message: err.Error(), Status: false, Code: http.StatusUnauthorized}
	data, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Write(data)
}

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
