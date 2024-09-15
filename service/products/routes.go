package products

import (
	"fmt"
	"net/http"

	"github.com/fayleenpc/tj-jeans/service/auth"
	"github.com/fayleenpc/tj-jeans/types"
	"github.com/fayleenpc/tj-jeans/utils"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", auth.WithLogger(auth.WithRateLimiter(h.handleGetProducts))).Methods("GET")
	router.HandleFunc("/products", auth.WithRateLimiter(auth.WithJWTAuth(h.handleCreateProduct, h.userStore))).Methods("POST")
	router.HandleFunc("/products", auth.WithRateLimiter(auth.WithJWTAuth(h.handleUpdateProduct, h.userStore))).Methods("PATCH")
	router.HandleFunc("/products", auth.WithRateLimiter(auth.WithJWTAuth(h.handleDeleteProduct, h.userStore))).Methods("DELETE")
}

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

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
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
		utils.WriteJSON(w, http.StatusOK, payload)
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
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
		utils.WriteJSON(w, http.StatusOK, payload)
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
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
		utils.WriteJSON(w, http.StatusOK, payload)
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}
