package tokenize

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/mailer"
	"github.com/fayleenpc/tj-jeans/internal/ratelimiter"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store      types.TokenStore
	userStore  types.UserStore
	redisStore *redis.Client
}

func NewHandler(store types.TokenStore, userStore types.UserStore, redisStore *redis.Client) *Handler {
	return &Handler{store: store, userStore: userStore, redisStore: redisStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/refresh", ratelimiter.WithRateLimiter(h.handleRefresh)).Methods("POST")
	router.HandleFunc("/logout", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleLogout), h.userStore, h.store)).Methods("POST")
	router.HandleFunc("/blacklisted_tokens", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetBlacklistedTokens), h.userStore, h.store)).Methods("GET")
	router.HandleFunc("/verify", ratelimiter.WithRateLimiter(h.handleVerify)).Methods("GET")
	router.HandleFunc("/login", ratelimiter.WithRateLimiter(h.handleLogin)).Methods("POST")
	router.HandleFunc("/register", ratelimiter.WithRateLimiter(h.handleRegister)).Methods("POST")
}

func validateRefreshToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func (h *Handler) handleGetBlacklistedTokens(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetBlacklistedTokens")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	tokens, err := h.store.GetBlacklistedTokens()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tokens)
}

// func readCookieHandler(w http.ResponseWriter, r *http.Request) string {
// 	var request struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}
// 	if cookie, err := r.Cookie("secureTokenCookie"); err == nil {

// 		if err = s.Decode("secureTokenCookie", cookie.Value, &request.RefreshToken); err == nil {
// 			fmt.Fprintf(w, "The value of secureTokenCookie is %q", request.RefreshToken)
// 		}
// 	}
// 	return request.RefreshToken
// }

// handleRefresh godoc
//
//	@Summary		Refresh a JWT Token using secretToken
//	@Description	Refresh a JWT Token using secretToken then give accessToken
//	@Tags			tokenize
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.Token
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/refresh [get]
func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleRefresh")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	var payload types.RefreshTokenPayload
	if auth.GetUserIDFromContext(r.Context()) != 0 {
		log.Println("got user for refreshing token")
	}
	if r.Header.Get("Authorization") != "" && r.Header.Get("Authorization-X") != "" {
		payload.AccessToken = r.Header.Get("Authorization")
		payload.SecretToken = r.Header.Get("Authorization-X")
		log.Printf("Authorization : %v, Authorization-X: %v\n", r.Header.Get("Authorization"), r.Header.Get("Authorization-X"))
	} else {
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}
	log.Printf("%+v", payload)
	t, err := h.store.GetBlacklistTokenByString(payload.SecretToken)
	if t.Token == payload.SecretToken && err == nil {
		log.Printf("blacklisted token: %v", t)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("blacklisted token jwt secret %v", t))
		return
	}

	// Blacklist the token
	if _, err := h.store.CreateBlacklistTokens(types.Token{Token: payload.AccessToken}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := validateRefreshToken(payload.SecretToken)
	if err != nil || !token.Valid {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, _ := strconv.Atoi(claims["userID"].(string))
	userRole := claims["userRole"].(string)
	userName := claims["userName"].(string)
	userEmail := claims["userEmail"].(string)
	userPhoneNumber := claims["userPhoneNumber"].(string)
	userAddress := claims["userAddress"].(string)

	newAccessToken, _, err := auth.CreateJWT([]byte(config.Envs.JWTRefresh), []byte(nil), userID, userRole, userName, userEmail, userPhoneNumber, userAddress)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"access_token": newAccessToken, "secret_token": payload.SecretToken})
}

// handleLogout godoc
//
//	@Summary		Logout a JWT Token using secretToken
//	@Description	Logout a JWT Token using secretToken then give accessToken
//	@Tags			tokenize
//	@Accept			json
//	@Produce		json
//	@Success		204	{object}
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/logout [post]
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleLogout")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	var payload struct {
		SecretToken string `json:"secret_token"`
		AccessToken string `json:"access_token"`
	}

	if r.Header.Get("Authorization") != "" && r.Header.Get("Authorization-X") != "" {
		payload.AccessToken = r.Header.Get("Authorization")
		payload.SecretToken = r.Header.Get("Authorization-X")
		log.Printf("Authorization : %v, Authorization-X: %v\n", r.Header.Get("Authorization"), r.Header.Get("Authorization-X"))
	} else {
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}
	log.Printf("%+v", payload)
	t, err := h.store.GetBlacklistTokenByString(payload.SecretToken)
	if t.Token == payload.SecretToken && err == nil {
		log.Printf("blacklisted token: %v", t)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("blacklisted token jwt secret %v", t))
		return
	}

	// Blacklist the token
	if _, err := h.store.CreateBlacklistTokens(types.Token{Token: payload.SecretToken}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// Blacklist the token
	if _, err := h.store.CreateBlacklistTokens(types.Token{Token: payload.AccessToken}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleVerify godoc
//
//	@Summary		Verify a user to API
//	@Description	Verify a user to API using in-memory storage
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	mailer.EmailVerificationTokens
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/verify [get]
func (h *Handler) handleVerify(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleVerify")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)

	token := r.URL.Query().Get("token")
	// verify_email := r.URL.Query().Get("verify_user")
	email, exists := mailer.EmailVerificationTokens[token]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid or expired token"))
		return
	}

	if err := h.userStore.UpdateVerifiedUserByEmail(email); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email"))
		return
	}
	log.Println(mailer.EmailVerificationTokens)

	delete(mailer.EmailVerificationTokens, token)
	// Here you can update the user's status in your database or mark the email as verified.
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Email %s has been verified successfully!\n", email))
}

// handleLogin godoc
//
//	@Summary		Login a user to API
//	@Description	Login a user to API then Create JWT Token (accessToken, secretToken)
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.LoginUserPayload
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/login [post]
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleLogin")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)

	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate the payload

	if err := utils.Validate.Struct(payload); err != nil {
		errv := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errv))
		return
	}

	u, err := h.userStore.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid password"))
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	access := []byte(config.Envs.JWTRefresh)
	accessToken, secretToken, err := auth.CreateJWT(access, secret, u.ID, u.Role, u.FirstName+" "+u.LastName, u.Email, u.PhoneNumber, u.Address)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// session.SetJWTSecretToken(w, secretToken)
	// session.SetJWTAccessToken(w, accessToken)

	utils.WriteJSON(w, http.StatusOK, map[string]string{"access_token": accessToken, "secret_token": secretToken})

}

// handleRegister godoc
//
//	@Summary		Register a user to API
//	@Description	Register a user to API without JWT Token
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.RegisterUserPayload
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/register [post]
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleRegister")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)

	// get json payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", error))
		return
	}

	// check if the user exists
	_, err := h.userStore.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// get token verification
	tokenVerification, err := mailer.GenerateToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error generating token verification"))
		return
	}
	// verification here
	mailer.EmailVerificationTokens[tokenVerification] = payload.Email
	// if err := auth.SendVerificationEmail(payload.Email, tokenVerification); err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error sending verification email"))
	// 	return
	// }

	// if it doesn't create new user
	err = h.userStore.CreateUser(types.User{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Password:    hashedPassword,
		PhoneNumber: payload.PhoneNumber,
		Address:     payload.Address,
		Verified:    false,
		Role:        "customer",
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"verify_url": config.Envs.PublicHost + ":" + config.Envs.Port + "/api/v1/verify?token=" + tokenVerification})
}
