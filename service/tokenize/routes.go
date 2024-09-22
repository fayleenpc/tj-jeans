package tokenize

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.TokenStore
	userStore types.UserStore
}

func NewHandler(store types.TokenStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/refresh", h.handleRefresh).Methods("GET")
	router.HandleFunc("/logout", h.handleLogout).Methods("POST")
}

func validateRefreshToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
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
	// userID := auth.GetUserIDFromContext(r.Context())
	// log.Printf("/refresh/%v\n", userID)

	secretToken, exists := session.SessionStoreServer.Get()

	if !exists {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// check blacklisted token db
	// c, _ := r.Cookie("secretTokenCookie")

	t, err := h.store.GetBlacklistTokenByString(secretToken)
	if t.Token == secretToken && err == nil {
		log.Printf("blacklisted token: %v", t)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("blacklisted token jwt secret %v", t))
		return
	}

	// Blacklist the token
	if err := h.store.CreateBlacklistTokens(types.Token{Token: session.SessionStoreClient.Blacklist()}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := validateRefreshToken(secretToken)
	if err != nil || !token.Valid {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	id, _ := strconv.Atoi(claims["userID"].(string))
	role := claims["userRole"].(string)
	username := claims["userName"].(string)

	newAccessToken, _, err := auth.CreateJWT([]byte(config.Envs.JWTRefresh), []byte(nil), id, role, username)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	session.SessionStoreClient.Refresh(newAccessToken)
	v, _ := session.SessionStoreServer.Get()
	log.Printf("session client %v\n", session.SessionStoreClient.String())
	log.Printf("session server %v\n", v)
	// log.Printf("session cookie %v\n", c.Value)

	utils.WriteJSON(w, http.StatusOK, map[string]string{"access_token": newAccessToken})
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

	// if err := utils.ParseJSON(r, &request); err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }

	secretToken, exists := session.SessionStoreServer.Get()

	if !exists {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// check blacklisted token db
	// c, _ := r.Cookie("secretTokenCookie")
	t, err := h.store.GetBlacklistTokenByString(secretToken)
	if t.Token == secretToken && err == nil {
		log.Printf("blacklisted token: %v", t)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("blacklisted token jwt secret %v", t))
		return
	}

	// Blacklist the token
	if err := h.store.CreateBlacklistTokens(types.Token{Token: secretToken}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// Blacklist the token
	if err := h.store.CreateBlacklistTokens(types.Token{Token: session.SessionStoreClient.Blacklist()}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	v, _ := session.SessionStoreServer.Get()
	log.Printf("session client %v\n", session.SessionStoreClient.String())
	log.Printf("session server %v\n", v)
	// log.Printf("session cookie %v\n", c.Value)
	w.WriteHeader(http.StatusNoContent)
}
