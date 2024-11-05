package products

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
	Client pb.ProductServiceClient
}

func NewHandlerClient(client pb.ProductServiceClient) *HandlerClient {
	return &HandlerClient{Client: client}
}

type HandlerHTTP struct {
	client types.ProductService
}

func NewHandlerHTTP(client types.ProductService) *HandlerHTTP {
	return &HandlerHTTP{client: client}
}

func (h *HandlerHTTP) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/products", h.handleGetProducts_Proto)
	mux.HandleFunc("POST /api/v1/products", h.handleUpdateProduct_Proto)
	mux.HandleFunc("PATCH /api/v1/products", h.handleUpdateProduct_Proto)
	mux.HandleFunc("DELETE /api/v1/products", h.handleDeleteProduct_Proto)
	mux.HandleFunc("GET /api/v1/products/{product_id}", h.handleGetProductByID_Proto)
	mux.HandleFunc("PATCH /api/v1/products/{product_id}/update", h.handleUpdateProductByID_Proto)
	mux.HandleFunc("DELETE /api/v1/products/{product_id}/delete", h.handleDeleteProductByID_Proto)

}

func (h *HandlerHTTP) handleGetProductByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetProductByID_Proto")
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
		product, err := h.client.GetProductByID(r.Context(), &pb.GetProductByIDRequest{Id: int32(productID)})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, product.GetProduct())
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *HandlerHTTP) handleUpdateProductByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateProductByID_Proto")
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
		var payload pb.Product
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldProduct, err := h.client.GetProductByID(r.Context(), &pb.GetProductByIDRequest{Id: int32(productID)})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		if int(oldProduct.GetProduct().GetId()) != int(payload.GetId()) {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		updatedProductID, err := h.client.UpdateProduct(r.Context(), &pb.UpdateProductRequest{Product: &payload})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": updatedProductID.GetUpdatedCount(), "old_product": oldProduct.GetProduct(), "updated_product": &payload})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *HandlerHTTP) handleDeleteProductByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteProductByIDProto")
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
		oldProduct, err := h.client.GetProductByID(r.Context(), &pb.GetProductByIDRequest{Id: int32(productID)})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		deletedOldProductID, err := h.client.DeleteProductByID(r.Context(), &pb.DeleteProductByIDRequest{Id: oldProduct.GetProduct().GetId()})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": deletedOldProductID.GetDeletedCount(), "deleted_product": oldProduct.GetProduct()})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
	}

}

func (h *HandlerHTTP) handleGetProducts_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetProducts_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)

	ps, err := h.client.GetProducts(r.Context(), &pb.GetProductsRequest{})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// PublishForProducts(ps)
	// SubscribeForProduct(w, r)

	utils.WriteJSON(w, http.StatusOK, ps.GetProducts())
}

func (h *HandlerHTTP) handleCreateProduct_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleCreateProduct_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	var payload pb.Product
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if userRole == "admin" {
		responseCreated, err := h.client.CreateProduct(r.Context(), &pb.CreateProductRequest{})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		payload.Id = int32(responseCreated.GetId())
		utils.WriteJSON(w, http.StatusOK, map[string]any{"created_id": responseCreated.GetId(), "created_product": &payload})
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}

func (h *HandlerHTTP) handleUpdateProduct_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateProduct_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	var payload pb.Product
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if userRole == "admin" {
		responseUpdated, err := h.client.UpdateProduct(r.Context(), &pb.UpdateProductRequest{Product: &payload})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		payload.Id = int32(responseUpdated.GetUpdatedCount())
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": responseUpdated.GetUpdatedCount(), "updated_product": &payload})
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}
func (h *HandlerHTTP) handleDeleteProduct_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteProduct_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	var payload pb.Product
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if userRole == "admin" {
		responseDeleted, err := h.client.DeleteProduct(r.Context(), &pb.DeleteProductRequest{Product: &payload})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		payload.Id = int32(responseDeleted.GetDeletedCount())
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": responseDeleted.GetDeletedCount(), "deleted_product": &payload})
	} else {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
	}

}
