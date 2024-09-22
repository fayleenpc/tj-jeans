package products

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/logger"
	"github.com/fayleenpc/tj-jeans/internal/ratelimiter"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store      types.ProductStore
	userStore  types.UserStore
	tokenStore types.TokenStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore, tokenStore types.TokenStore) *Handler {
	return &Handler{store: store, userStore: userStore, tokenStore: tokenStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", logger.WithLogger(ratelimiter.WithRateLimiter(h.handleGetProducts))).Methods("GET")
	router.HandleFunc("/products", ratelimiter.WithRateLimiter(auth.WithJWTAuth(h.handleCreateProduct, h.userStore, h.tokenStore))).Methods("POST")
	router.HandleFunc("/products", ratelimiter.WithRateLimiter(auth.WithJWTAuth(h.handleUpdateProduct, h.userStore, h.tokenStore))).Methods("PATCH")
	router.HandleFunc("/products", ratelimiter.WithRateLimiter(auth.WithJWTAuth(h.handleDeleteProduct, h.userStore, h.tokenStore))).Methods("DELETE")
}

// handleGetProducts godoc
//
//	@Summary		Get a products using JWT Token ( accessToken )
//	@Description	Get a products using JWT Token ( accessToken ), with login credentials
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]types.Product
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/products [get]
func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetProducts")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// PublishForProducts(ps)
	// SubscribeForProduct(w, r)

	utils.WriteJSON(w, http.StatusOK, ps)
}

// handleCreateProduct godoc
//
//	@Summary		Post a products using JWT Token ( accessToken )
//	@Description	Post a products using JWT Token ( accessToken ), with login credentials ( role admin )
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.Product
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/products [post]
func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// userRole := auth.GetUserRoleFromContext(r.Context())
	userID := auth.GetUserIDFromSession(session.SessionStoreClient.Get(), h.userStore)
	userRole := auth.GetUserRoleFromSession(session.SessionStoreClient.Get())
	// set context "userID" to the userID
	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.UserKey, userID)
	ctx = context.WithValue(ctx, auth.UserRoleKey, userRole)
	r = r.WithContext(ctx)
	var payload types.Product
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if userRole == "admin" {
		id, err := h.store.CreateProduct(payload)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		payload.ID = int(id)
		utils.WriteJSON(w, http.StatusOK, payload)
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}

// handleUpdateProduct godoc
//
//	@Summary		Update a products using JWT Token ( accessToken )
//	@Description	Update a products using JWT Token ( accessToken ), with login credentials ( role admin )
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.Product
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/products [patch]
func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	// userRole := auth.GetUserRoleFromContext(r.Context())
	userID := auth.GetUserIDFromSession(session.SessionStoreClient.Get(), h.userStore)
	userRole := auth.GetUserRoleFromSession(session.SessionStoreClient.Get())
	// set context "userID" to the userID
	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.UserKey, userID)
	ctx = context.WithValue(ctx, auth.UserRoleKey, userRole)
	r = r.WithContext(ctx)
	var payload types.Product
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if userRole == "admin" {
		id, err := h.store.UpdateProduct(payload)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		payload.ID = int(id)
		utils.WriteJSON(w, http.StatusOK, payload)
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}

// handleDeleteProduct godoc
//
//	@Summary		Delete a products using JWT Token ( accessToken )
//	@Description	Delete a products using JWT Token ( accessToken ), with login credentials ( role admin )
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.Product
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/products [delete]
func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// userRole := auth.GetUserRoleFromContext(r.Context())
	userID := auth.GetUserIDFromSession(session.SessionStoreClient.Get(), h.userStore)
	userRole := auth.GetUserRoleFromSession(session.SessionStoreClient.Get())
	// set context "userID" to the userID
	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.UserKey, userID)
	ctx = context.WithValue(ctx, auth.UserRoleKey, userRole)
	r = r.WithContext(ctx)
	var payload types.Product
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if userRole == "admin" {
		id, err := h.store.DeleteProduct(payload)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		payload.ID = int(id)
		utils.WriteJSON(w, http.StatusOK, payload)
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}
