package cart

import (
	"fmt"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/logger"
	"github.com/fayleenpc/tj-jeans/internal/ratelimiter"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
	tokenStore   types.TokenStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore, tokenStore types.TokenStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore, tokenStore: tokenStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", logger.WithLogger(ratelimiter.WithRateLimiter(auth.WithJWTAuth(h.handleCheckout, h.userStore, h.tokenStore)))).Methods("POST")
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

	// userID := auth.GetUserIDFromContext(r.Context())
	userID := auth.GetUserIDFromSession(session.SessionStoreClient.Get(), h.userStore)
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
