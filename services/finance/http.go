package finance

import (
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/types"
)

type HandlerHTTP struct {
	client types.UserService
}

func NewHandlerHTTP(client types.UserService) *HandlerHTTP {
	return &HandlerHTTP{client: client}
}

func (h *HandlerHTTP) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/users", h.handleGetFinance_Proto)
}

func (h *HandlerHTTP) handleGetFinance_Proto(w http.ResponseWriter, r *http.Request) {

}
