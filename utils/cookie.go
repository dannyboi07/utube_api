package utils

import (
	"net/http"
	"utube/common"
	"utube/schema"
)

func CreateAccessTokenCookie(user schema.UserForToken) (*http.Cookie, error) {
	token, expiresIn, err := CreateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		MaxAge:   expiresIn,
		Path:     common.AccessTokenPath,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}, nil
}

func CreateRefreshTokenCookie(user schema.UserForToken) (*http.Cookie, error) {
	token, expiresIn, err := CreateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     "refreshToken",
		Value:    token,
		MaxAge:   expiresIn,
		Path:     common.RefreshTokenPath,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}, nil
}
