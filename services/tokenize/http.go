package tokenize

import (
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HandlerClient struct {
	client pb.TokenServiceClient
}

func NewHandlerClient(client pb.TokenServiceClient) *HandlerClient {
	return &HandlerClient{client: client}
}

type HandlerHTTP struct {
	client types.TokenService
}

func NewHandlerHTTP(client types.TokenService) *HandlerHTTP {
	return &HandlerHTTP{client: client}
}

func (h *HandlerHTTP) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/blacklisted_tokens", h.handleGetBlacklistedTokens_Proto)
}

func (h *HandlerHTTP) handleGetBlacklistedTokens_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetBlacklistedTokens_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)

	resp, err := h.client.GetBlacklistedTokens(r.Context(), &pb.GetBlacklistedTokensRequest{})
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp.GetTokens())
}
