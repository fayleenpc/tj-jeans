package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/ratelimiter"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store      types.UserStore
	tokenStore types.TokenStore
	redisStore *redis.Client
}

func NewHandler(store types.UserStore, tokenStore types.TokenStore, redisStore *redis.Client) *Handler {
	return &Handler{store: store, tokenStore: tokenStore, redisStore: redisStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetUsers), h.store, h.tokenStore)).Methods("GET")
	router.HandleFunc("/users/{user_id}", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetUserByID), h.store, h.tokenStore)).Methods("GET")
	router.HandleFunc("/users/{user_id}/update", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleUpdateUserByID), h.store, h.tokenStore)).Methods("PATCH")
	router.HandleFunc("/users/{user_id}/delete", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleDeleteUserByID), h.store, h.tokenStore)).Methods("DELETE")
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetUsers")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())

	if userRole == "admin" {
		redisString := h.redisStore.Get(r.Context(), "handleGetUsers")
		if result, err := redisString.Result(); err == nil {
			utils.WriteJSON(w, http.StatusOK, result)
		}
		users, err := h.store.GetUsers()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		redisStatus := h.redisStore.Set(r.Context(), "handleGetUsers", users, 3600)
		if redisStatus.Err() != nil {
			log.Fatalf("Could not set key: %v", err)
		}

		utils.WriteJSON(w, http.StatusOK, users)
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetUserByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		user, err := h.store.GetUserByID(userID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, user)
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleUpdateUserByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateUserByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		var payload types.User
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldUser, err := h.store.GetUserByID(userID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		if oldUser.ID != payload.ID {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		updatedUserID, err := h.store.UpdateUser(payload)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": updatedUserID, "old_user": oldUser, "updated_user": payload})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleDeleteUserByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteUserByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldUser, err := h.store.GetUserByID(userID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		deletedUserID, err := h.store.DeleteUserByID(oldUser.ID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": deletedUserID, "deleted_user": oldUser})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission, denied"))
	}

}
