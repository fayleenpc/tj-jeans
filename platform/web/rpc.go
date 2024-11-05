package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/loadbalancer"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/gorilla/mux"
)

var (
	lb = loadbalancer.NewLoadBalancer()
)

func logout(w http.ResponseWriter, r *http.Request, callback func()) {
	callback()
}

func login(w http.ResponseWriter, r *http.Request, callback func(status int, responseLogin types.ResponseLogin)) {
	var payload types.LoginUserPayload
	var response types.ResponseLogin
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := json.Marshal(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resBody, status, err := utils.CraftJSON("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/login", m, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("location %v, got response %+v, payload %v\n", r.URL.Path, response, payload)

	session.SetJWTAccessToken(w, response.AccessToken)
	session.SetJWTSecretToken(w, response.SecretToken)
	callback(status, response)
}

func register(w http.ResponseWriter, r *http.Request, callback func(status int, responseRegister types.ResponseRegister)) {
	var payload types.RegisterUserPayload
	var response types.ResponseRegister
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := json.Marshal(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resBody, status, err := utils.CraftJSON("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/register", m, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("location %v, got response %+v\n", r.URL.Path, response)
	callback(status, response)
}

func refresh(w http.ResponseWriter, r *http.Request, callback func(status int, responseRefresh types.ResponseRefreshToken)) {
	var payload types.RefreshTokenPayload
	var response types.ResponseRefreshToken

	if r.Header.Get("Authorization") != "" && r.Header.Get("Authorization-X") != "" {
		payload.AccessToken = r.Header.Get("Authorization")
		payload.SecretToken = r.Header.Get("Authorization-X")
	} else {
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}
	log.Printf("[check payload /service/refresh] : %+v\n", payload)
	m, err := json.Marshal(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resBody, status, err := utils.CraftJSON("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/refresh", m, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("location %v, got response %s\n", r.URL.Path, response)
	callback(status, response)
}

func createProducts(w http.ResponseWriter, r *http.Request, callback func(status int, response types.ResponseProduct)) {
	var payload types.Product
	var response types.ResponseProduct

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := json.Marshal(payload)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	// resBody, status, err := utils.CraftJSON("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/products", m, r)
	// do loadbalancer for frequent products
	resBody, status, err := utils.CraftJSON("POST", lb.GetBackend()+"/api/v1/products", m, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("location %v, got response %+v\n", r.URL.Path, response)
	callback(status, response)
}

func updateProducts(w http.ResponseWriter, r *http.Request, callback func(status int, response types.ResponseProduct)) {
	var payload types.Product
	var response types.ResponseProduct

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := json.Marshal(payload)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resBody, status, err := utils.CraftJSON("PATCH", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/products", m, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("location %v, got response %+v\n", r.URL.Path, response)
	callback(status, response)
}

func deleteProducts(w http.ResponseWriter, r *http.Request, callback func(status int, response types.ResponseProduct)) {
	var payload types.Product
	var response types.ResponseProduct

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := json.Marshal(payload)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resBody, status, err := utils.CraftJSON("DELETE", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/products", m, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("location %v, got response %+v\n", r.URL.Path, response)
	callback(status, response)
}

func handleInvoice(w http.ResponseWriter, r *http.Request, req types.CartCheckoutPayload, res types.ResponseCart, callback func(req types.InvoicePayload, responseInvoice types.InvoiceResponse)) {
	var invoicePayload types.InvoicePayload
	var invoiceResponse types.InvoiceResponse

	_ = req
	// invoices payment gateway
	invoicePayload.Payment.Type = "dana"
	invoicePayload.Payment.Amount = res.Total
	invoicePayload.Customer.Name = auth.GetUserNameFromSession(r.Header.Get("Authorization"))
	invoicePayload.Customer.Email = auth.GetUserEmailFromSession(r.Header.Get("Authorization"))
	invoicePayload.Customer.PhoneNumber = auth.GetUserPhoneNumberFromSession(r.Header.Get("Authorization"))
	invoicePayload.Items = res.Items

	marshalInvoice, err := json.Marshal(invoicePayload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	resBodyInvoice, statusInvoice, err := utils.CraftJSON("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/payment/invoices", marshalInvoice, r)
	if err != nil {
		utils.WriteError(w, statusInvoice, err)
		return
	}

	if err := json.Unmarshal(resBodyInvoice, &invoiceResponse); err != nil {
		utils.WriteError(w, statusInvoice, err)
		return
	}

	// log.Printf("location %v, set invoice %+v\n", r.URL.Path, invoicePayload)
	// log.Printf("location %v, got invoice %+v\n", r.URL.Path, invoiceResponse)
	// log.Printf("response final at /payment/invoices %v\n", string(resBodyInvoice))

	callback(invoicePayload, invoiceResponse)

}

func handleCart(w http.ResponseWriter, r *http.Request, callback func(req types.CartCheckoutPayload, res types.ResponseCart)) {
	var cart types.CartCheckoutPayload
	var response types.ResponseCart

	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	marshalCart, err := json.Marshal(cart)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	resBodyCart, statusCart, err := utils.CraftJSON("POST", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/cart/checkout", marshalCart, r)
	if err != nil {
		utils.WriteError(w, statusCart, err)
		return
	}

	if err := json.Unmarshal(resBodyCart, &response); err != nil {
		utils.WriteError(w, statusCart, err)
		return
	}
	callback(cart, response)
}

func getOrderItems(r *http.Request) ([]types.OrderItem, error) {
	var payload []types.OrderItem
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/order_items", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func getOrderItemByID(r *http.Request) (*types.OrderItem, error) {
	orderItemID := mux.Vars(r)["order_item_id"]
	var payload types.OrderItem
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/order_items"+fmt.Sprintf("/%v", orderItemID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func updateOrderItemByID(r *http.Request) (any, error) {
	orderItemID := mux.Vars(r)["order_item_id"]
	var payload types.OrderItem
	var response struct {
		ID               int             `json:"updated_id"`
		OldOrderItem     types.OrderItem `json:"old_order_item"`
		UpdatedOrderItem types.OrderItem `json:"updated_order_item"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		return nil, err
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/order_items"+fmt.Sprintf("/%v/update", orderItemID), bytes.NewBuffer(m))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func deleteOrderItemByID(r *http.Request) (any, error) {
	orderItemID := mux.Vars(r)["order_item_id"]
	var response struct {
		ID      int              `json:"deleted_id"`
		Deleted *types.OrderItem `json:"deleted_order_item"`
	}
	req, err := http.NewRequest("DELETE", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/order_items"+fmt.Sprintf("/%v/delete", orderItemID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func getOrders(r *http.Request) ([]types.Order, error) {
	var payload []types.Order
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/orders", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func getOrderByID(r *http.Request) (*types.Order, error) {
	orderID := mux.Vars(r)["order_id"]
	var payload types.Order
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/orders"+fmt.Sprintf("/%v", orderID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func updateOrderByID(r *http.Request) (any, error) {
	orderID := mux.Vars(r)["order_id"]
	var payload types.Order
	var response struct {
		ID           int         `json:"updated_id"`
		OldOrder     types.Order `json:"old_order"`
		UpdatedOrder types.Order `json:"updated_order"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		return nil, err
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/orders"+fmt.Sprintf("/%v/update", orderID), bytes.NewBuffer(m))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func deleteOrderByID(r *http.Request) (any, error) {
	orderID := mux.Vars(r)["order_id"]
	var response struct {
		ID      int          `json:"deleted_id"`
		Deleted *types.Order `json:"deleted_order"`
	}
	req, err := http.NewRequest("DELETE", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/orders"+fmt.Sprintf("/%v/delete", orderID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func getCustomers(r *http.Request) ([]types.User, error) {
	var payload []types.User
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/users", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func getCustomerByID(r *http.Request) (*types.User, error) {
	customerID := mux.Vars(r)["customer_id"]
	var payload types.User
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/users"+fmt.Sprintf("/%v", customerID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func updateCustomerByID(r *http.Request) (any, error) {
	customerID := mux.Vars(r)["customer_id"]
	var payload types.User
	var response struct {
		ID              int        `json:"updated_id"`
		OldCustomer     types.User `json:"old_customer"`
		UpdatedCustomer types.User `json:"updated_customer"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		return nil, err
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/users"+fmt.Sprintf("/%v/update", customerID), bytes.NewBuffer(m))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func deleteCustomerByID(r *http.Request) (any, error) {
	customerID := mux.Vars(r)["customer_id"]
	var response struct {
		ID      int         `json:"deleted_id"`
		Deleted *types.User `json:"deleted_customer"`
	}
	req, err := http.NewRequest("DELETE", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/users/"+fmt.Sprintf("/%v/delete", customerID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func getProducts(r *http.Request) ([]types.Product, error) {
	var payload []types.Product
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/products", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func getProductByID(r *http.Request) (*types.Product, error) {
	productID := mux.Vars(r)["product_id"]
	var payload types.Product
	req, err := http.NewRequest("GET", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/products"+fmt.Sprintf("/%v", productID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func updateProductByID(r *http.Request) (any, error) {
	productID := mux.Vars(r)["product_id"]
	var payload types.Product
	var response struct {
		ID             int           `json:"updated_id"`
		OldProduct     types.Product `json:"old_product"`
		UpdatedProduct types.Product `json:"updated_product"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		return nil, err
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/products"+fmt.Sprintf("/%v/update", productID), bytes.NewBuffer(m))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func deleteProductByID(r *http.Request) (any, error) {
	productID := mux.Vars(r)["product_id"]
	var response struct {
		ID      int            `json:"deleted_id"`
		Deleted *types.Product `json:"deleted_product"`
	}
	req, err := http.NewRequest("DELETE", config.Envs.PublicHost+":"+config.Envs.Port+"/api/v1/products/"+fmt.Sprintf("/%v/delete", productID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
