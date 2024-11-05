package cart

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/ratelimiter"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
	tokenStore   types.TokenStore
	redisStore   *redis.Client
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore, tokenStore types.TokenStore, redisStore *redis.Client) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore, tokenStore: tokenStore, redisStore: redisStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleCheckout), h.userStore, h.tokenStore)).Methods("POST")
	router.HandleFunc("/orders", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetOrders), h.userStore, h.tokenStore)).Methods("GET")
	router.HandleFunc("/orders/{order_id}", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetOrderByID), h.userStore, h.tokenStore)).Methods("GET")
	router.HandleFunc("/orders/{order_id}/update", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleUpdateOrderByID), h.userStore, h.tokenStore)).Methods("PATCH")
	router.HandleFunc("/orders/{order_id}/delete", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleDeleteOrderByID), h.userStore, h.tokenStore)).Methods("DELETE")
	router.HandleFunc("/order_items", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetOrderItems), h.userStore, h.tokenStore)).Methods("GET")
	router.HandleFunc("/order_items/{order_item_id}", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetOrderItemByID), h.userStore, h.tokenStore)).Methods("GET")
	router.HandleFunc("/order_items/{order_item_id}/update", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleUpdateOrderItemByID), h.userStore, h.tokenStore)).Methods("PATCH")
	router.HandleFunc("/order_items/{order_item_id}/delete", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleDeleteOrderItemByID), h.userStore, h.tokenStore)).Methods("DELETE")
}

func (h *Handler) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetOrders")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	if userRole == "admin" {
		orders, err := h.store.GetOrders()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, orders)
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}
}

func (h *Handler) handleGetOrderByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetOrderByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	orderID, err := strconv.Atoi(mux.Vars(r)["order_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		order, err := h.store.GetOrderByID(orderID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, order)
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleUpdateOrderByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateOrderByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	orderID, err := strconv.Atoi(mux.Vars(r)["order_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		var payload types.Order
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldOrder, err := h.store.GetOrderByID(orderID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		if oldOrder.ID != payload.ID {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		updatedOrderID, err := h.store.UpdateOrder(payload)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": updatedOrderID, "old_order": oldOrder, "updated_order": payload})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleDeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteOrderByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	orderID, err := strconv.Atoi(mux.Vars(r)["order_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldOrder, err := h.store.GetOrderByID(orderID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		deletedOldOrderID, err := h.store.DeleteOrderByID(oldOrder.ID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": deletedOldOrderID, "deleted_order": oldOrder})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleGetOrderItems(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetOrderItems")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	if userRole == "admin" {
		orderItems, err := h.store.GetOrderItems()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, orderItems)
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}
}

func (h *Handler) handleGetOrderItemByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetOrderItemByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	orderItemID, err := strconv.Atoi(mux.Vars(r)["order_item_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		orderItem, err := h.store.GetOrderItemsByID(orderItemID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, orderItem)
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleUpdateOrderItemByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateOrderItemByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	orderItemID, err := strconv.Atoi(mux.Vars(r)["order_item_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		var payload types.OrderItem
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldOrderItem, err := h.store.GetOrderItemsByID(orderItemID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		if oldOrderItem.ID != payload.ID {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		updatedOrderItemID, err := h.store.UpdateOrderItem(payload)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": updatedOrderItemID, "old_order_item": oldOrderItem, "updated_order_item": payload})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *Handler) handleDeleteOrderItemByID(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteOrderItemByID")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	orderItemID, err := strconv.Atoi(mux.Vars(r)["order_item_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldOrderItem, err := h.store.GetOrderItemsByID(orderItemID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		deletedOldOrderID, err := h.store.DeleteOrderItemByID(oldOrderItem.ID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": deletedOldOrderID, "deleted_order_item": oldOrderItem})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

// handleCheckout godoc
//
//	@Summary		Checkout a products using JWT Token ( accessToken )
//	@Description	Checkout a products using JWT Token ( accessToken ), with login credentials ( role admin & customer )
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.CartCheckoutPayload
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/api/v1/cart/checkout [post]
func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleCheckout")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)

	userID := auth.GetUserIDFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errv := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errv))
		return
	}
	// get product
	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ps, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	orderId, totalPrice, err := h.createOrder(ps, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// helper payment
	newPs := []types.Product(ps)
	for i := 0; i < len(cart.Items); i++ {
		newPs[i].Quantity = cart.Items[i].Quantity
	}
	// d, err := json.Marshal(map[string]any{
	// 	"total_price": totalPrice,
	// 	"order_id":    orderId,
	// 	"items":       newPs,
	// })
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, err)
	// }
	// span.SetTag("payload checkout", d)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderId,
		"items":       newPs,
	})
}
