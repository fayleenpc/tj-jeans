package user

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fayleenpc/tj-jeans/config"
	"github.com/fayleenpc/tj-jeans/service/auth"
	"github.com/fayleenpc/tj-jeans/types"
	"github.com/fayleenpc/tj-jeans/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store types.UserStore
	sync.WaitGroup
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.WithLogger(auth.WithRateLimiter(h.handleLogin))).Methods("POST")
	router.HandleFunc("/register", auth.WithLogger(auth.WithRateLimiter(h.handleRegister))).Methods("POST")
	router.HandleFunc("/verify", h.handleVerify).Methods("GET")
}

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
	token, err := auth.CreateJWT(secret, u.ID, u.Role)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tokenCookie := &http.Cookie{
		Name:     "token_cookie",
		Value:    token,
		Expires:  time.Unix(config.Envs.JWTExpirationInSeconds, 0),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
		Secure:   true,
		Path:     "/",
	}

	http.SetCookie(w, tokenCookie)

	utils.WriteJSON(w, http.StatusSeeOther, map[string]string{"token": token})

	// fmt.Fprintf(w, "Registration successful. Please check your email to verify.")

	// if u.Email == "jackwdy1@gmail.com" {
	// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token, "admin": "true"})
	// } else {
	// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	// }

}

func (h *Handler) handleVerify(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleVerify")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)

	token := r.URL.Query().Get("token")
	// verify_email := r.URL.Query().Get("verify_user")
	email, exists := auth.EmailVerificationTokens[token]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid or expired token"))
		return
	}

	if err := h.store.UpdateVerifiedUserByEmail(email); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email"))
		return
	}
	log.Println(auth.EmailVerificationTokens)

	delete(auth.EmailVerificationTokens, token)
	// Here you can update the user's status in your database or mark the email as verified.
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Email %s has been verified successfully!\n", email))
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	return
}

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
	tokenVerification, err := auth.GenerateToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error generating token verification"))
		return
	}
	// verification here
	auth.EmailVerificationTokens[tokenVerification] = payload.Email
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

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"verify_url": config.Envs.PublicHost + "/api/v1/verify?token=" + tokenVerification})
}
