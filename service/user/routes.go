package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/logger"
	"github.com/fayleenpc/tj-jeans/internal/mailer"
	"github.com/fayleenpc/tj-jeans/internal/ratelimiter"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", logger.WithLogger(ratelimiter.WithRateLimiter(h.handleLogin))).Methods("POST")
	router.HandleFunc("/login", logger.WithLogger(ratelimiter.WithRateLimiter(h.handleLoginCookie))).Methods("GET")
	router.HandleFunc("/register", logger.WithLogger(ratelimiter.WithRateLimiter(h.handleRegister))).Methods("POST")
	router.HandleFunc("/verify", h.handleVerify).Methods("GET")
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

	if err := h.store.UpdateVerifiedUserByEmail(email); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email"))
		return
	}
	log.Println(mailer.EmailVerificationTokens)

	delete(mailer.EmailVerificationTokens, token)
	// Here you can update the user's status in your database or mark the email as verified.
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Email %s has been verified successfully!\n", email))
}

func (h *Handler) handleLoginCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("Set-Cookie")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	log.Printf("Req Cookie is : %v\n", c.Value)
	log.Printf("Res Cookie is : %v\n", w.Header().Get("Set-Cookie"))
	utils.WriteJSON(w, http.StatusOK, nil)
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
	log.Println(payload)
	//validate the payload

	if err := utils.Validate.Struct(payload); err != nil {
		errv := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errv))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	log.Println(u)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	access := []byte(config.Envs.JWTRefresh)
	accessToken, secretToken, err := auth.CreateJWT(access, secret, u.ID, u.Role, u.FirstName+" "+u.LastName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	c := &http.Cookie{
		Name:     "secretTokenCookie",
		Value:    secretToken,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, c)
	cc := &http.Cookie{
		Name:     "accessTokenCookie",
		Value:    accessToken,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cc)
	session.SessionStoreServer.Set(secretToken)
	session.SessionStoreClient.Set(accessToken)
	v, _ := session.SessionStoreServer.Get()
	log.Printf("session client %v\n", session.SessionStoreClient.String())
	log.Printf("session server %v\n", v)
	// log.Printf("session cookie %v\n", c.Value)
	r.Header.Set("Authorization", accessToken)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})

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
	_, err := h.store.GetUserByEmail(payload.Email)
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
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
		Verified:  false,
		Role:      "customer",
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"verify_url": config.Envs.PublicHost + ":" + config.Envs.Port + "/api/v1/verify?token=" + tokenVerification})
}
