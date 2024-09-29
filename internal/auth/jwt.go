package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"
const UserNameKey contextKey = "userName"
const UserEmailKey contextKey = "userEmail"
const UserPhoneNumberKey contextKey = "userPhoneNumber"
const userAddressKey contextKey = "userAddress"

func WithCookie(handlerFunc http.HandlerFunc, store types.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(time.Second * 5)
		if vv, err := session.GetJWTSecretToken(r); err != nil {
			log.Printf("not found cookie at middleware - secret token - location %v\n", r.URL.Path)
		} else {
			// log.Printf("got cookie at middleware - secret token %v - location %v\n", v, r.URL.Path)
			t, err := store.GetBlacklistTokenByString(vv)
			r.Header.Set("Authorization-X", vv)
			log.Printf("header [secret] has been set :  %v - location %v\n", r.Header.Get("Authorization-X"), r.URL.Path)
			if err == nil {
				log.Printf("got blacklisted cookie at middleware - secret token %v - location %v\n", t.Token, r.URL.Path)
			}

		}

		if v, err := session.GetJWTAccessToken(r); err != nil {
			log.Printf("not found cookie at middleware - access token - location %v\n", r.URL.Path)
		} else {
			// log.Printf("got cookie at middleware - access token  %v - location %v\n", v, r.URL.Path)
			t, err := store.GetBlacklistTokenByString(v)
			r.Header.Set("Authorization", v)
			log.Printf("header [access] has been set :  %v - location %v\n", r.Header.Get("Authorization"), r.URL.Path)
			if err == nil {
				log.Printf("got blacklisted cookie at middleware - access token %v - location %v \n", t.Token, r.URL.Path)
			}
		}

		// log.Printf("header final [access] has been set :  %v - location %v\n", r.Header.Get("Authorization"), r.URL.Path)
		// log.Printf("header final [secret] has been set :  %v - location %v\n", r.Header.Get("Authorization-X"), r.URL.Path)

		handlerFunc(w, r)
	}
}

func CreateJWT(access, secret []byte, userID int, userRole string, userName string, userEmail string, userPhoneNumber string, userAddress string) (string, string, error) {
	var secretTokenString, accessTokenString string
	var err error
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	secretToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":          strconv.Itoa(userID),
		"userRole":        userRole,
		"userName":        userName,
		"userEmail":       userEmail,
		"userPhoneNumber": userPhoneNumber,
		"userAddress":     userAddress,
		"expiredAt":       time.Now().Add(expiration).Unix(),
	})
	secretTokenString, _ = secretToken.SignedString(secret)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":          strconv.Itoa(userID),
		"userRole":        userRole,
		"userName":        userName,
		"userEmail":       userEmail,
		"userPhoneNumber": userPhoneNumber,
		"userAddress":     userAddress,
		"expiredAt":       time.Now().Add(time.Minute * 1).Unix(),
	})
	accessTokenString, err = accessToken.SignedString(access)

	if secretTokenString == "" && accessTokenString != "" {
		return accessTokenString, "", nil
	}
	if err != nil {
		return "", "", err
	}
	return accessTokenString, secretTokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, userStore types.UserStore, tokenStore types.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request path : %v\n", r.RequestURI)
		var (
			JWT_SECRET_TOKEN  *jwt.Token
			JWT_ACCESS_TOKEN  *jwt.Token
			JWT_SECRET_STRING string
			JWT_ACCESS_STRING string
		)
		JWT_ACCESS_STRING, JWT_SECRET_STRING = getTokenFromRequest(r)

		t2, err := tokenStore.GetBlacklistTokenByString(JWT_ACCESS_STRING)
		if t2.Token == JWT_ACCESS_STRING && err == nil {
			log.Printf("blacklisted token: %v", t2)
			permissionDenied(w)
			return
		}

		t1, err := tokenStore.GetBlacklistTokenByString(JWT_SECRET_STRING)
		if t1.Token == JWT_SECRET_STRING && err == nil {
			log.Printf("blacklisted token: %v", t1)
			permissionDenied(w)
			return
		}
		log.Printf("JWT_ACCESS_TOKEN=%v\n", JWT_ACCESS_STRING)
		log.Printf("JWT_SECRET_TOKEN=%v\n", JWT_SECRET_STRING)
		// validate the JWT Access Token
		JWT_SECRET_TOKEN, err = validateSecretToken(JWT_SECRET_STRING)
		if err != nil {
			log.Printf("failed to validate secret token: %v", err)
			permissionDenied(w)
			return
		}

		// validate the JWT Secret Token
		JWT_ACCESS_TOKEN, err = validateAccessToken(JWT_ACCESS_STRING)
		if err != nil {
			log.Printf("failed to validate access token: %v", err)
			permissionDenied(w)
			return
		}

		if !JWT_ACCESS_TOKEN.Valid {
			log.Println("invalid access token")
			permissionDenied(w)
			return
		}
		if !JWT_SECRET_TOKEN.Valid {
			log.Println("invalid secret token")
			permissionDenied(w)
			return
		}

		// if is we need to fetch the userID from the DB (id from the token)
		claims := JWT_ACCESS_TOKEN.Claims.(jwt.MapClaims)
		userID_str := claims["userID"].(string)
		userRole := claims["userRole"].(string)
		userEmail := claims["userEmail"].(string)
		userPhoneNumber := claims["userPhoneNumber"].(string)
		userAddress := claims["userAddress"].(string)
		userID, _ := strconv.Atoi(userID_str)
		u, err := userStore.GetUserByID(userID)
		userName := u.FirstName + " " + u.LastName
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// set context "userID" to the userID
		ctx := r.Context()
		if userID == u.ID {
			ctx = context.WithValue(ctx, UserKey, userID)
		}
		if userRole == u.Role {
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
		}
		if userEmail == u.Email {
			ctx = context.WithValue(ctx, UserEmailKey, userEmail)
		}
		if userName == u.FirstName+" "+u.LastName {
			ctx = context.WithValue(ctx, UserNameKey, userName)
		}
		if userPhoneNumber == u.PhoneNumber {
			ctx = context.WithValue(ctx, UserPhoneNumberKey, userPhoneNumber)
		}
		if userAddress == u.Address {
			ctx = context.WithValue(ctx, userAddressKey, userAddress)
		}
		r = r.WithContext(ctx)
		log.Printf("USER DATA BACKEND ( userID=%v, userRole=%v, userEmail=%v, userName=%v, userPhoneNumber=%v, userAddress=%v )\n", userID, userRole, userEmail, userName, userPhoneNumber, userAddress)
		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) (string, string) {
	authHeader := r.Header.Get("Authorization")
	authXHeader := r.Header.Get("Authorization-X")
	if authHeader == "" && authXHeader == "" {
		return "", ""
	}
	return authHeader, authXHeader
}

func getTokenFromSession() string {

	return session.SessionStoreClient.Get()
}

func IsTokenExpired(v string) bool {
	token, err := validateAccessToken(v)
	if err != nil {
		return true
	}
	claims := token.Claims.(jwt.MapClaims)
	expiredAt := claims["expiredAt"].(float64)

	return time.Now().After(time.Unix(int64(expiredAt), 0))
}

func CheckExpired(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	log.Println("IsTokenExpired? ", IsTokenExpired(authHeader))
	return authHeader, IsTokenExpired(authHeader)
}

func ShouldRenew(w http.ResponseWriter, r *http.Request) (string, bool) {
	renewed := RenewAccessToken(r)
	if renewed == "" {
		return "", false
	}
	log.Printf("location : %v\n Renewed : %v\n", r.URL.Path, renewed)

	return renewed, true
}

func ShouldRevoke(w http.ResponseWriter, r *http.Request) bool {

	revoked := RevokeToken(r)
	log.Printf("location : %v\n Revoked : %v\n", r.URL.Path, revoked == http.StatusOK)

	return revoked == http.StatusOK
}

func RevokeToken(r *http.Request) int {
	var payload struct {
		AccessToken string `json:"access_token"`
		SecretToken string `json:"secret_token"`
	}
	accessToken, err := session.GetJWTAccessToken(r)
	if err != nil {
		return int(http.StatusBadRequest)
	}
	secretToken, err := session.GetJWTSecretToken(r)
	if err != nil {
		return int(http.StatusBadRequest)
	}
	payload.AccessToken = accessToken
	payload.SecretToken = secretToken

	m, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/service/logout", bytes.NewBuffer(m))

	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	log.Println("req /service/logout : ", res.Status)
	return int(res.StatusCode)
}

func RenewAccessToken(r *http.Request) string {
	var payload struct {
		AccessToken string `json:"access_token"`
		SecretToken string `json:"secret_token"`
	}
	var response struct {
		AccessToken string `json:"access_token"`
	}
	accessToken, err := session.GetJWTAccessToken(r)
	if err != nil {
		return ""
	}
	secretToken, err := session.GetJWTSecretToken(r)
	if err != nil {
		return ""
	}
	payload.AccessToken = accessToken
	payload.SecretToken = secretToken
	m, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/refresh", bytes.NewBuffer(m))

	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	log.Println("req /api/v1/refresh : ", response)
	return response.AccessToken
}

func validateAccessToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTRefresh), nil
	})
}

func validateSecretToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}

func GetUserNameFromContext(ctx context.Context) string {
	userName, ok := ctx.Value(UserNameKey).(string)
	if !ok {
		return ""
	}
	return userName
}

func GetUserRoleFromContext(ctx context.Context) string {
	userRole, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return ""
	}
	return userRole
}

func GetUserIDFromSession(v string, userStore types.UserStore) int {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	userID_str := claims["userID"].(string)
	if err != nil {
		log.Fatal(err)
	}
	userID, _ := strconv.Atoi(userID_str)
	u, err := userStore.GetUserByID(userID)
	if err != nil {
		log.Fatal(err)
	}
	return u.ID
}

func GetUserPhoneNumberFromSession(v string) string {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	userPhoneNumber := claims["userPhoneNumber"].(string)
	if userPhoneNumber == "" {
		log.Fatal(err)
	}
	return userPhoneNumber
}

func GetUserAddressFromSession(v string) string {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	userAddress := claims["userAddress"].(string)
	if userAddress == "" {
		log.Fatal(err)
	}
	return userAddress
}

func GetUserEmailFromSession(v string) string {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	userEmail := claims["userEmail"].(string)
	if userEmail == "" {
		log.Fatal(err)
	}
	return userEmail
}

func GetUserNameFromSession(v string) string {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	userName := claims["userName"].(string)
	if userName == "" {
		log.Fatal(err)
	}
	return userName
}

func GetUserRoleFromSession(v string) string {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	userRole := claims["userRole"].(string)
	if userRole == "" {
		log.Fatal(err)
	}
	return userRole
}

func BridgeAdmin(w http.ResponseWriter, r *http.Request) bool {
	switch {
	case r.Header.Get("Authorization") != "" && r.Header.Get("Authorization-X") != "":
		if GetUserRoleFromSession(r.Header.Get("Authorization")) == "admin" {
			_, ok := CheckExpired(r)
			if ok {
				if renewed, ok := ShouldRenew(w, r); ok {
					session.SetJWTAccessToken(w, renewed)
					return true
				}
				return true
			}
			return true
		} else {
			return false
		}

	default:
		return false
	}
}

func BridgeCommon(w http.ResponseWriter, r *http.Request) bool {
	switch {
	case r.Header.Get("Authorization") != "" && r.Header.Get("Authorization-X") != "":
		if _, ok := CheckExpired(r); ok {
			if renewed, ok := ShouldRenew(w, r); ok {
				session.SetJWTAccessToken(w, renewed)
				return true
			}

			return true
		} else {
			return true
		}
	default:
		return false
	}
}
