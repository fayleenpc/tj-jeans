package cart

import (
	"fmt"
	"net/http"

	"github.com/fayleenpc/tj-jeans/service/auth"
	"github.com/fayleenpc/tj-jeans/types"
	"github.com/fayleenpc/tj-jeans/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithLogger(auth.WithRateLimiter(auth.WithJWTAuth(h.handleCheckout, h.userStore)))).Methods("POST")
}

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
