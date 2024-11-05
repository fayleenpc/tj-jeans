package cart

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type HandlerClient struct {
	client pb.OrderServiceClient
}

func NewHandlerClient(client pb.OrderServiceClient) *HandlerClient {
	return &HandlerClient{
		client: client,
	}
}

type HandlerHTTP struct {
	client types.OrderService
}

func NewHandlerHTTP(client types.OrderService) *HandlerHTTP {
	return &HandlerHTTP{client: client}
}

func (h *HandlerHTTP) RegisterRoutes(mux *http.ServeMux) {
	// router.HandleFunc("/cart/checkout", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleCheckout), h.userStore, h.tokenStore)).Methods("POST")
	mux.HandleFunc("GET /api/v1/orders", h.handleGetOrders_Proto)
	mux.HandleFunc("PATCH /api/v1/orders/{order_id}/update", h.handleUpdateOrderByID_Proto)
	mux.HandleFunc("DELETE /api/v1/orders/{order_id}/delete", h.handleDeleteOrderByID_Proto)
	mux.HandleFunc("GET /api/v1/order_items", h.handleGetOrderByID_Proto)

	// router.HandleFunc("/orders/{order_id}", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetOrderByID), h.userStore, h.tokenStore)).Methods("GET")
	// router.HandleFunc("/order_items/{order_item_id}", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleGetOrderItemByID), h.userStore, h.tokenStore)).Methods("GET")
	// router.HandleFunc("/order_items/{order_item_id}/update", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleUpdateOrderItemByID), h.userStore, h.tokenStore)).Methods("PATCH")
	// router.HandleFunc("/order_items/{order_item_id}/delete", auth.WithJWTAuth(ratelimiter.WithRateLimiter(h.handleDeleteOrderItemByID), h.userStore, h.tokenStore)).Methods("DELETE")
}

func (h *HandlerHTTP) handleGetOrders_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetOrders_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	if userRole == "admin" {
		orders, err := h.client.GetOrders(r.Context(), &pb.GetOrdersRequest{})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, orders.GetOrders())
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}
}

func (h *HandlerHTTP) handleGetOrderByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetOrderByID_Proto")
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
		order, err := h.client.GetOrderByID(r.Context(), &pb.GetOrderByIDRequest{Id: int32(orderID)})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, order.GetOrder())
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}
func (h *HandlerHTTP) handleUpdateOrderByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateOrderByID_Proto")
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
		var payload pb.Order
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldOrder, err := h.client.GetOrderByID(r.Context(), &pb.GetOrderByIDRequest{Id: int32(orderID)})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		if int(oldOrder.GetOrder().GetId()) != int(payload.GetId()) {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		updatedOrderID, err := h.client.UpdateOrder(r.Context(), &pb.UpdateOrderRequest{Order: &payload})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": updatedOrderID.GetUpdatedCount(), "old_order": oldOrder.GetOrder(), "updated_order": &payload})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}
func (h *HandlerHTTP) handleDeleteOrderByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteOrderByID_Proto")
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
		oldOrder, err := h.client.GetOrderByID(r.Context(), &pb.GetOrderByIDRequest{Id: int32(orderID)})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		deletedOldOrderID, err := h.client.DeleteOrderByID(r.Context(), &pb.DeleteOrderByIDRequest{Id: oldOrder.GetOrder().GetId()})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": deletedOldOrderID.GetDeletedCount(), "deleted_order": oldOrder.GetOrder()})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}
