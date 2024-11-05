package products

import (
	"fmt"
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
	store      types.ProductStore
	userStore  types.UserStore
	tokenStore types.TokenStore
	redisStore *redis.Client
}

func NewHandler(store types.ProductStore, userStore types.UserStore, tokenStore types.TokenStore, redisStore *redis.Client) *Handler {
	return &Handler{store: store, userStore: userStore, tokenStore: tokenStore, redisStore: redisStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", ratelimiter.WithRateLimiter(h.handleGetProducts)).Methods("GET")
	router.HandleFunc("/products", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleCreateProduct), h.userStore, h.tokenStore)).Methods("POST")
	router.HandleFunc("/products", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleUpdateProduct), h.userStore, h.tokenStore)).Methods("PATCH")
	router.HandleFunc("/products", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleDeleteProduct), h.userStore, h.tokenStore)).Methods("DELETE")
	router.HandleFunc("/products/{product_id}", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetProductByID), h.userStore, h.tokenStore)).Methods("GET")
	router.HandleFunc("/products/{product_id}/update", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleUpdateProductByID), h.userStore, h.tokenStore)).Methods("PATCH")
	router.HandleFunc("/products/{product_id}/delete", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleDeleteProductByID), h.userStore, h.tokenStore)).Methods("DELETE")
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetProductByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	productID, err := strconv.Atoi(mux.Vars(r)["product_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		product, err := h.store.GetProductByID(productID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, product)
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleUpdateProductByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateProductByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	productID, err := strconv.Atoi(mux.Vars(r)["product_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		var payload types.Product
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldProduct, err := h.store.GetProductByID(productID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		if oldProduct.ID != payload.ID {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		updatedProductID, err := h.store.UpdateProduct(payload)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": updatedProductID, "old_product": oldProduct, "updated_product": payload})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleDeleteProductByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteProductByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	productID, err := strconv.Atoi(mux.Vars(r)["product_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldProduct, err := h.store.GetProductByID(productID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		deletedOldProductID, err := h.store.DeleteProductByID(oldProduct.ID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": deletedOldProductID, "deleted_product": oldProduct})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

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
	span := opentracing.GlobalTracer().StartSpan("handleCreateProduct")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
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
		utils.WriteJSON(w, http.StatusOK, map[string]any{"created_id": id, "created_product": payload})
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
	span := opentracing.GlobalTracer().StartSpan("handleUpdateProduct")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
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
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": id, "updated_product": payload})
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
	span := opentracing.GlobalTracer().StartSpan("handleDeleteProduct")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
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
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": id, "deleted_product": &payload})
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}
