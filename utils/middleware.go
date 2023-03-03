package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func MultipartFormDataRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			WriteApiErrMessage(w, http.StatusBadRequest, "Invalid content type")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func JsonRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			WriteApiErrMessage(w, http.StatusBadRequest, "Invalid content type")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			accessTokenCookie *http.Cookie
			err               error
		)
		accessTokenCookie, err = r.Cookie("accessToken")
		if err != nil {
			WriteApiErrMessage(w, http.StatusUnauthorized, "Missing access token")
			return
		}

		var (
			jwtClaims  jwt.MapClaims
			statusCode int
		)
		jwtClaims, statusCode, err = VerifyJwtToken(accessTokenCookie.Value)
		if err != nil {
			WriteApiErrMessage(w, statusCode, err.Error())
			return
		}

		var userDetails map[string]interface{}
		userDetails, statusCode, err = ParseJwtClaims(jwtClaims)
		if err != nil {
			WriteApiErrMessage(w, statusCode, err.Error())
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userDetails", userDetails)))
	})
}
