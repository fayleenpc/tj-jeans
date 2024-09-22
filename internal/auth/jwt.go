package auth

import (
	"context"
	"fmt"
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

func CreateJWT(access, secret []byte, userID int, role string, userName string) (string, string, error) {
	var secretTokenString, accessTokenString string
	var err error
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	secretToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"userRole":  role,
		"userName":  userName,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	secretTokenString, _ = secretToken.SignedString(secret)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"userRole":  role,
		"userName":  userName,
		"expiredAt": time.Now().Add(time.Hour * 1).Unix(),
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
		var (
			token       *jwt.Token
			tokenString string
		)
		log.Printf("request path : %v\n", r.RequestURI)
		// get the token from user request
		// tokenString = getTokenFromRequest(r)

		// get the token from user session
		tokenString = getTokenFromSession()
		// check blacklisted token db
		t, err := tokenStore.GetBlacklistTokenByString(tokenString)
		if t.Token == tokenString && err == nil {
			log.Printf("blacklisted token: %v", t)
			permissionDenied(w)
			return
		} else {
			if r.RequestURI == "/api/v1/refresh" || r.RequestURI == "/api/v1/logout" {
				// validate the JWT Access Token
				token, err = validateSecretToken(tokenString)
				if err != nil {
					log.Printf("failed to validate token: %v", err)
					permissionDenied(w)
					return
				}
			} else {
				// validate the JWT Access Token
				token, err = validateAccessToken(tokenString)
				if err != nil {
					log.Printf("failed to validate token: %v", err)
					permissionDenied(w)
					return
				}
			}
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// if is we need to fetch the userID from the DB (id from the token)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, _ := strconv.Atoi(str)
		u, err := userStore.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}
		str_role := claims["userRole"].(string)
		if str_role == "" {
			log.Printf("failed to get userRole: %v", err)
			permissionDenied(w)
			return
		}
		// set context "userID" to the userID
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		ctx = context.WithValue(ctx, UserRoleKey, u.Role)
		ctx = context.WithValue(ctx, UserNameKey, u.FirstName+" "+u.LastName)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenAuth_Cookie, err := r.Cookie("secretTokenCookie")
	if err != nil {
		log.Printf("cookie got here : %v\n", tokenAuth_Cookie)
	}
	log.Println(tokenAuth)
	if tokenAuth != "" {
		return tokenAuth
	}
	return ""
}

func getTokenFromSession() string {

	return session.SessionStoreClient.Get()
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
	str := claims["userID"].(string)

	userID, _ := strconv.Atoi(str)
	u, err := userStore.GetUserByID(userID)
	if err != nil {
		log.Fatal(err)
	}
	return u.ID
}

func GetUserNameFromSession(v string) string {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	str_role := claims["userName"].(string)
	if str_role == "" {
		log.Fatal(err)
	}
	return str_role
}

func GetUserRoleFromSession(v string) string {
	token, err := validateAccessToken(v)
	// if is we need to fetch the userID from the DB (id from the token)
	claims := token.Claims.(jwt.MapClaims)
	str_role := claims["userRole"].(string)
	if str_role == "" {
		log.Fatal(err)
	}
	return str_role
}
