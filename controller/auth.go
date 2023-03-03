package controller

import (
	"encoding/json"
	"net/http"
	"utube/db"
	"utube/models"
	"utube/schema"
	"utube/utils"

	"github.com/golang-jwt/jwt/v4"
)

var refToken string = ""

// /auth/register
func Register(w http.ResponseWriter, r *http.Request) {
	jDec := json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		userRegister schema.UserRegister
		statusCode   int
		err          error
	)
	statusCode, err = utils.JsonParseErr(jDec.Decode(&userRegister))
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to parse register user request: ", err)
		return
	}

	statusCode, err = userRegister.Validate()
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("User register request validation failed: ", err)
		return
	}

	var userExists bool
	userExists, err = db.ActorExistsByEmail(*userRegister.Email)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to check for existing actor, err:", err)
		return
	}

	if userExists {
		utils.WriteApiErrMessage(w, http.StatusBadRequest, "Account already exists, please login")
		utils.Log.Println("Register request, user already exists")
		return
	}

	var hashedPw string
	hashedPw, err = utils.HashPassword(*userRegister.Password)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to hash password, errr:", err)
		return
	}

	var user models.Actor = models.Actor{
		Email:    *userRegister.Email,
		Password: hashedPw,
	}

	user, err = db.InsertActor(user)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to insert new user into db, err:", err)
		return
	}

	utils.WriteApiMessage(w, 0, "Account created")
}

// /auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	jDec := json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		userLogin  schema.UserLogin
		statusCode int
		err        error
	)
	statusCode, err = utils.JsonParseErr(jDec.Decode(&userLogin))
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to parse user login request, err:", err)
		return
	}

	statusCode, err = userLogin.Validate()
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("User login request validation failed:", err)
		return
	}

	var (
		user  models.Actor
		found bool
	)
	user, found, err = db.SelectUserByEmail(*userLogin.Email)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to get user by email, err:", err)
		return
	} else if !found {
		utils.WriteApiErrMessage(w, http.StatusBadRequest, "Check your email/password")
		utils.Log.Println("User doesn't exist for login")
		return
	}

	if isCorrectPw := utils.VerifyPassword(user.Password, *userLogin.Password); !isCorrectPw {
		utils.WriteApiErrMessage(w, http.StatusBadRequest, "Check your email/password")
		utils.Log.Println("Wrong password for user email: ", *userLogin.Email)
		return
	}

	var userTokenDetails schema.UserForToken = schema.UserForToken{
		Id:    user.Id,
		Email: user.Email,
	}

	var accessTokenCookie *http.Cookie
	accessTokenCookie, err = utils.CreateAccessTokenCookie(userTokenDetails)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to create access token cookie, err:", err)
		return
	}

	var refreshTokenCookie *http.Cookie
	refreshTokenCookie, err = utils.CreateRefreshTokenCookie(userTokenDetails)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to create refresh token cookie, err:", err)
		return
	}

	refToken = refreshTokenCookie.Value

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	utils.WriteApiMessage(w, 0, "Logged in")
}

func RefreshTokenOldBrokenForSomeGodForsakenReason(w http.ResponseWriter, r *http.Request) {
	var (
		refreshTokenCookie *http.Cookie
		err                error
	)
	refreshTokenCookie, err = r.Cookie("refreshToken")
	if err != nil {
		utils.WriteApiErrMessage(w, http.StatusForbidden, "Session expired")
		utils.Log.Println("Missing refresh token")
		return
	}

	var (
		mapClaims  jwt.MapClaims
		statusCode int
	)
	utils.Log.Println("token", refreshTokenCookie.Value)
	utils.Log.Println("ref tk eq", refToken == refreshTokenCookie.Value)
	mapClaims, statusCode, err = utils.VerifyJwtToken(refreshTokenCookie.Value)
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to verify refresh token, err:", err)
		return
	}

	var userDetails map[string]interface{}
	userDetails, statusCode, err = utils.ParseJwtClaims(mapClaims)
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to parse refresh token claims, err:", err)
		return
	}

	var (
		userId uint64
		ok     bool
	)
	userId, ok = userDetails["userId"].(uint64)
	if !ok {
		utils.Log.Println("Typo failed")
	}
	var email string = userDetails["email"].(string)

	var userForToken schema.UserForToken = schema.UserForToken{
		Id:    userId,
		Email: email,
	}
	var accessTokenCookie *http.Cookie
	accessTokenCookie, err = utils.CreateAccessTokenCookie(userForToken)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to create access token, err:", err)
		return
	}

	http.SetCookie(w, accessTokenCookie)
	w.WriteHeader(http.StatusOK)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var (
		refreshTokenCookie *http.Cookie
		err                error
	)
	refreshTokenCookie, err = r.Cookie("refreshToken")
	if err != nil {
		utils.WriteApiErrMessage(w, http.StatusForbidden, "Session expired")
		utils.Log.Println("Missing refresh token")
		return
	}

	var (
		mapClaims  jwt.MapClaims
		statusCode int
	)
	mapClaims, statusCode, err = utils.VerifyJwtToken(refreshTokenCookie.Value)
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to verify refresh token, err:", err)
		return
	}

	var userDetails map[string]interface{}
	userDetails, statusCode, err = utils.ParseJwtClaims(mapClaims)
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to parse refresh token claims, err:", err)
		return
	}

	var (
		userId uint64
	)
	userId, _ = userDetails["id"].(uint64)
	var email string = userDetails["email"].(string)

	var userForToken schema.UserForToken = schema.UserForToken{
		Id:    userId,
		Email: email,
	}
	var accessTokenCookie *http.Cookie
	accessTokenCookie, err = utils.CreateAccessTokenCookie(userForToken)
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "")
		utils.Log.Println("Failed to create access token, err:", err)
		return
	}

	http.SetCookie(w, accessTokenCookie)
	w.WriteHeader(http.StatusOK)
}
