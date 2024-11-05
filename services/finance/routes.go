package finance

import (
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	orderStore types.OrderStore
}

func NewHandler(orderStore types.OrderStore) *Handler {
	return &Handler{orderStore: orderStore}
}

func (h *Handler) RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc("/finance", h.handleFinance).Methods("GET")
	// Additional routes can be added here
}

func (h *Handler) handleFinance(w http.ResponseWriter, r *http.Request) {
	orderItems, err := h.orderStore.GetOrderItems()
	if err != nil {
		http.Error(w, "Unable to fetch order items", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	orders, err := h.orderStore.GetOrders()
	if err != nil {
		http.Error(w, "Unable to fetch orders", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Calculate revenue and send response
	result := calculateDashboard(orderItems, orders)
	utils.WriteJSON(w, http.StatusOK, result)
}

func calculateDashboard(orderItems []types.OrderItem, orders []types.Order) map[string]any {
	totalRevenue := float64(0)
	totalItems := 0

	// Calculate total revenue from order items
	for _, item := range orderItems {
		totalRevenue += item.Price * float64(item.Quantity)
		totalItems += item.Quantity
	}

	// // Optionally, if you want to calculate based on orders too, you can modify this part
	// // Here I'm calculating a sample value; you can adjust as per your requirements
	// for _, order := range orders {
	// 	// Just an example of processing each order
	// 	totalRevenue += order.Total
	// }

	return map[string]any{
		"total_revenue":    totalRevenue,
		"total_items_sold": totalItems,
		"order_count":      len(orders),
	}
}
