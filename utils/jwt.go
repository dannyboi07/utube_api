package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"utube/common"
	"utube/schema"

	"github.com/golang-jwt/jwt/v4"
)

type jwtCustomClaims struct {
	*jwt.RegisteredClaims
	schema.UserForToken
}

func CreateAccessToken(userDetails schema.UserForToken) (string, int, error) {
	var token *jwt.Token = jwt.New(jwt.GetSigningMethod("RS256"))

	var createdTime time.Time = time.Now()
	var expireAtTime time.Time = createdTime.Add(time.Minute * 15)

	token.Claims = jwtCustomClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAtTime),
		},
		userDetails,
	}

	signedToken, err := token.SignedString(common.PrivateKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, int(expireAtTime.Sub(createdTime).Seconds()), nil
}

func CreateRefreshToken(userDetails schema.UserForToken) (string, int, error) {
	var token *jwt.Token = jwt.New(jwt.GetSigningMethod("RS256"))

	var createdTime time.Time = time.Now()
	var expireAtTime time.Time = createdTime.AddDate(0, 0, 7)

	token.Claims = jwtCustomClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAtTime),
		},
		userDetails,
	}

	signedToken, err := token.SignedString(common.PrivateKey)
	if err != nil {
		return "", 0, err
	}

	// Log.Println("")
	claims, _, er := VerifyJwtToken(signedToken)
	// if er != nil {
	Log.Println("Failed verification, err", er, claims)
	// }
	return signedToken, int(expireAtTime.Sub(createdTime).Seconds()), nil
}

func VerifyJwtToken(token string) (jwt.MapClaims, int, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if tokenAlg := t.Method.Alg(); tokenAlg != "RS256" {
			return nil, fmt.Errorf("Unexpected signing method: %s", tokenAlg)
		}

		return common.PublicKey, nil
	})

	if err != nil {
		Log.Println("Failed to parse JWT token, err:", err)
		return nil, http.StatusInternalServerError, errors.New("Sorry, something went wrong")
	}

	if parsedToken.Valid {
		return parsedToken.Claims.(jwt.MapClaims), http.StatusOK, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, http.StatusBadRequest, errors.New("Malformed token")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, http.StatusUnauthorized, errors.New("Token expired")
	}

	Log.Println("Couldn't handle this JWT token, err:", err)
	return nil, http.StatusInternalServerError, errors.New("Sorry, something went wrong")
}

func ParseJwtClaims(jwtClaims jwt.MapClaims) (map[string]interface{}, int, error) {
	var embeddedDetails map[string]interface{} = make(map[string]interface{})

	if userIdFloat, ok := jwtClaims["id"].(float64); ok {
		embeddedDetails["id"] = uint64(userIdFloat)
	} else {
		return nil, http.StatusUnauthorized, fmt.Errorf("Malformed token")
	}

	if userEmail, ok := jwtClaims["email"].(string); ok {
		embeddedDetails["email"] = userEmail
	} else {
		return nil, http.StatusUnauthorized, fmt.Errorf("Malformed token")
	}

	return embeddedDetails, 0, nil
}
