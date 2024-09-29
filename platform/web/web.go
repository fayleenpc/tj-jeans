package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	"github.com/fayleenpc/tj-jeans/platform/web/views"
	"github.com/fayleenpc/tj-jeans/platform/web/views_admin"
	"github.com/gorilla/mux"
)

var (
	Client = session.NewSessionStoreClient("init")
	Server = session.NewSessionStoreServer()
)

type Handler struct {
	store types.TokenStore
}

func NewHandler(store types.TokenStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/", auth.WithCookie(h.showHomePage, h.store)).Methods("GET")
	router.HandleFunc("/products", auth.WithCookie(h.showProductsPage, h.store)).Methods("GET")
	router.HandleFunc("/gallery", auth.WithCookie(h.showGalleryPage, h.store)).Methods("GET")

	router.HandleFunc("/service", auth.WithCookie(h.showServicePage, h.store)).Methods("GET")
	router.HandleFunc("/service/register", auth.WithCookie(h.handleRegisterService, h.store)).Methods("POST")
	router.HandleFunc("/service/login", auth.WithCookie(h.handleLoginService, h.store)).Methods("POST")
	router.HandleFunc("/service/logout", auth.WithCookie(h.handleLogoutService, h.store)).Methods("POST")
	router.HandleFunc("/service/refresh", auth.WithCookie(h.handleRefreshService, h.store)).Methods("POST")

	router.HandleFunc("/cart/checkout", auth.WithCookie(h.handleCheckoutService, h.store)).Methods("POST")

	router.HandleFunc("/admin", auth.WithCookie(h.showAdminPage, h.store)).Methods("GET")
	router.HandleFunc("/admin/order_items", auth.WithCookie(h.showAdminOrderItemsPage, h.store)).Methods("GET")
	router.HandleFunc("/admin/order_items/{order_item_id}", auth.WithCookie(h.handleGetOrderItemByID, h.store)).Methods("GET")
	router.HandleFunc("/admin/order_items/{order_item_id}/update", auth.WithCookie(h.handleUpdateOrderItemByID, h.store)).Methods("PATCH")
	router.HandleFunc("/admin/order_items/{order_item_id}/delete", auth.WithCookie(h.handleDeleteOrderItemByID, h.store)).Methods("DELETE")

	router.HandleFunc("/admin/orders", auth.WithCookie(h.showAdminOrdersPage, h.store)).Methods("GET")
	router.HandleFunc("/admin/orders/{order_id}", auth.WithCookie(h.handleGetOrderByID, h.store)).Methods("GET")
	router.HandleFunc("/admin/orders/{order_id}/update", auth.WithCookie(h.handleUpdateOrderByID, h.store)).Methods("PATCH")
	router.HandleFunc("/admin/orders/{order_id}/delete", auth.WithCookie(h.handleDeleteOrderByID, h.store)).Methods("DELETE")

	router.HandleFunc("/admin/customers", auth.WithCookie(h.showAdminCustomersPage, h.store)).Methods("GET")
	router.HandleFunc("/admin/customers/{customer_id}", auth.WithCookie(h.handleGetCustomerByID, h.store)).Methods("GET")
	router.HandleFunc("/admin/customers/{customer_id}/update", auth.WithCookie(h.handleUpdateCustomerByID, h.store)).Methods("PATCH")
	router.HandleFunc("/admin/customers/{customer_id}/delete", auth.WithCookie(h.handleDeleteCustomerByID, h.store)).Methods("DELETE")

	router.HandleFunc("/admin/products", auth.WithCookie(h.showAdminProductsPage, h.store)).Methods("GET")
	router.HandleFunc("/admin/products", auth.WithCookie(h.handleCreateProducts, h.store)).Methods("POST")
	router.HandleFunc("/admin/products", auth.WithCookie(h.handleUpdateProducts, h.store)).Methods("PATCH")
	router.HandleFunc("/admin/products", auth.WithCookie(h.handleDeleteProducts, h.store)).Methods("DELETE")
	router.HandleFunc("/admin/products/{product_id}", auth.WithCookie(h.handleGetProductByID, h.store)).Methods("GET")
	router.HandleFunc("/admin/products/{product_id}/update", auth.WithCookie(h.handleUpdateProductByID, h.store)).Methods("PATCH")
	router.HandleFunc("/admin/products/{product_id}/delete", auth.WithCookie(h.handleDeleteProductByID, h.store)).Methods("DELETE")
}

func messageWhatsapp(r *http.Request, req types.InvoicePayload, responseInvoice types.InvoiceResponse) string {
	body := ""
	for k, v := range req.Items {
		body += fmt.Sprintf("\n\t\t%v - %v, %v, qty : %v, %v %v.\n", k+1, v.Name, v.Category, v.Quantity, v.Currency, v.Price)
	}
	messageForWhatsapp := fmt.Sprintf("\nAtas nama data penjual\n\tMerchant : %v\n\tNo Telp : %v\n\nAtas nama data pembeli (isi data dengan lengkap)\n\tNama : %v\n\tEmail: %v\n\tNo Telp : %v\n\tAlamat Tujuan Pengiriman : %v\n\tBarang Yang Dibeli : \n\t%v\nUntuk pembayaran totalnya seharga %v %v (silahkan konfirmasi jika data sudah benar),lalu akses url berikut %v\nMohon mengirimkan bukti transfer ketika sudah maka proses pengiriman dapat dilakukan, terima kasih.",
		"TJ Jeans",
		"0895-0520-8391",
		req.Customer.Name,
		req.Customer.Email,
		req.Customer.PhoneNumber,
		auth.GetUserAddressFromSession(r.Header.Get("Authorization")),
		body,
		responseInvoice.TransactionValues.Currency,
		int64(responseInvoice.TransactionValues.Total),
		responseInvoice.Payment.RedirectURL,
	)
	log.Println(messageForWhatsapp)
	return messageForWhatsapp
}

// Request is guarded by Authorization Header (access_token) for every commit
func (h *Handler) handleCheckoutService(w http.ResponseWriter, r *http.Request) {
	handleCart(w, r, func(req types.CartCheckoutPayload, responseCart types.ResponseCart) {
		// log.Printf("location %v, got response %+v\n", r.URL.Path, responseCart)
		handleInvoice(w, r, req, responseCart, func(req types.InvoicePayload, responseInvoice types.InvoiceResponse) {
			_ = messageWhatsapp(r, req, responseInvoice)
			utils.WriteJSON(w, http.StatusOK, responseInvoice)
		})
	})

}

// Request is guarded by Authorization Header (access_token) for every commit
func (h *Handler) handleCreateProducts(w http.ResponseWriter, r *http.Request) {
	createProducts(w, r, func(status int, response types.ResponseProduct) {
		utils.WriteJSON(w, status, response)
	})
}

// Request is guarded by Authorization Header (access_token) for every commit
func (h *Handler) handleUpdateProducts(w http.ResponseWriter, r *http.Request) {
	updateProducts(w, r, func(status int, response types.ResponseProduct) {
		utils.WriteJSON(w, status, response)
	})
}

// Request is guarded by Authorization Header (access_token) for every commit
func (h *Handler) handleDeleteProducts(w http.ResponseWriter, r *http.Request) {
	deleteProducts(w, r, func(status int, response types.ResponseProduct) {
		utils.WriteJSON(w, status, response)
	})
}

// Request is guarded by Authorization-X and Authorization Header or JSON POST (access_token, secret_token) for every commit
func (h *Handler) handleLogoutService(w http.ResponseWriter, r *http.Request) {
	if ok := auth.ShouldRevoke(w, r); ok {
		session.SetJWTAccessToken(w, "")
		session.SetJWTSecretToken(w, "")
		w.WriteHeader(http.StatusNoContent)
	}
}

// Request is guarded by Authorization-X and Authorization Header or JSON POST (access_token, secret_token) for every commit
func (h *Handler) handleRefreshService(w http.ResponseWriter, r *http.Request) {
	refresh(w, r, func(status int, responseRefresh types.ResponseRefreshToken) {
		utils.WriteJSON(w, status, responseRefresh)
	})

}

func (h *Handler) handleRegisterService(w http.ResponseWriter, r *http.Request) {
	register(w, r, func(status int, responseRegister types.ResponseRegister) {
		utils.WriteJSON(w, status, responseRegister)
	})

}

func (h *Handler) handleLoginService(w http.ResponseWriter, r *http.Request) {
	login(w, r, func(status int, responseLogin types.ResponseLogin) {
		utils.WriteJSON(w, status, responseLogin)
	})
}

func (h *Handler) showHomePage(w http.ResponseWriter, r *http.Request) {
	if auth.BridgeCommon(w, r) {
		views.Home(auth.GetUserNameFromSession(r.Header.Get("Authorization")), auth.GetUserRoleFromSession(r.Header.Get("Authorization"))).Render(r.Context(), w)
	} else {
		views.Home("", "").Render(r.Context(), w)
	}
}

func (h *Handler) showProductsPage(w http.ResponseWriter, r *http.Request) {
	ps, _ := getProducts(r)
	if auth.BridgeCommon(w, r) {
		views.Products(ps, auth.GetUserNameFromSession(r.Header.Get("Authorization")), auth.GetUserRoleFromSession(r.Header.Get("Authorization"))).Render(r.Context(), w)
	} else {
		views.Products(ps, "", "").Render(r.Context(), w)
	}

}

func (h *Handler) showGalleryPage(w http.ResponseWriter, r *http.Request) {
	if auth.BridgeCommon(w, r) {
		views.Gallery(auth.GetUserNameFromSession(r.Header.Get("Authorization")), auth.GetUserRoleFromSession(r.Header.Get("Authorization"))).Render(r.Context(), w)
	} else {
		views.Gallery("", "").Render(r.Context(), w)
	}

}

func (h *Handler) showServicePage(w http.ResponseWriter, r *http.Request) {
	if auth.BridgeCommon(w, r) {
		views.Login_Register(auth.GetUserNameFromSession(r.Header.Get("Authorization")), auth.GetUserRoleFromSession(r.Header.Get("Authorization"))).Render(r.Context(), w)

	} else {
		views.Login_Register("", "").Render(r.Context(), w)
	}

}

func (h *Handler) showAdminPage(w http.ResponseWriter, r *http.Request) {
	if auth.BridgeAdmin(w, r) {
		views_admin.Home(auth.GetUserNameFromSession(r.Header.Get("Authorization"))).Render(r.Context(), w)
	} else {
		views_admin.Error().Render(r.Context(), w)
	}

}

func (h *Handler) showAdminOrderItemsPage(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		orderItems, err := getOrderItems(r)

		if err != nil {
			log.Fatal(err)
		}
		views_admin.Order_Items(auth.GetUserNameFromSession(r.Header.Get("Authorization")), orderItems).Render(r.Context(), w)
	} else {
		views_admin.Error().Render(r.Context(), w)
	}

}

func (h *Handler) handleGetOrderItemByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		orderItem, err := getOrderItemByID(r)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(orderItem)
	} else {
		views_admin.Error().Render(r.Context(), w)
	}

}

func (h *Handler) handleUpdateOrderItemByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseUpdateOrderItemByID, err := updateOrderItemByID(r)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseUpdateOrderItemByID)

	} else {
		views_admin.Error().Render(r.Context(), w)
	}

}

func (h *Handler) handleDeleteOrderItemByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseDeleteOrderItemByID, err := deleteOrderItemByID(r)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseDeleteOrderItemByID)

	} else {
		views_admin.Error().Render(r.Context(), w)

	}

}

func (h *Handler) showAdminOrdersPage(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		orders, err := getOrders(r)

		if err != nil {
			log.Fatal(err)
		}
		views_admin.Orders(auth.GetUserNameFromSession(r.Header.Get("Authorization")), orders).Render(r.Context(), w)

	} else {
		views_admin.Error().Render(r.Context(), w)
	}
}

func (h *Handler) handleGetOrderByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		order, err := getOrderByID(r)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(order)

	} else {
		views_admin.Error().Render(r.Context(), w)

	}

}

func (h *Handler) handleUpdateOrderByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseUpdateOrderByID, err := updateOrderByID(r)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseUpdateOrderByID)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}

func (h *Handler) handleDeleteOrderByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseDeleteOrderByID, err := deleteOrderByID(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseDeleteOrderByID)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}

func (h *Handler) showAdminCustomersPage(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		customers, err := getCustomers(r)
		if err != nil {
			log.Fatal(err)
		}
		views_admin.Customers(auth.GetUserNameFromSession(r.Header.Get("Authorization")), customers).Render(r.Context(), w)
	} else {
		views_admin.Error().Render(r.Context(), w)
	}
}

func (h *Handler) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		customer, err := getCustomerByID(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(customer)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}

func (h *Handler) handleUpdateCustomerByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseUpdateCustomerByID, err := updateCustomerByID(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseUpdateCustomerByID)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}

func (h *Handler) handleDeleteCustomerByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseDeleteCustomerByID, err := deleteCustomerByID(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseDeleteCustomerByID)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}

func (h *Handler) showAdminProductsPage(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		products, err := getProducts(r)
		if err != nil {
			log.Fatal(err)
		}
		views_admin.Products(auth.GetUserNameFromSession(r.Header.Get("Authorization")), products).Render(r.Context(), w)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		product, err := getProductByID(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(product)
	} else {
		views_admin.Error().Render(r.Context(), w)
	}
}

func (h *Handler) handleUpdateProductByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseUpdateProductByID, err := updateProductByID(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseUpdateProductByID)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}

func (h *Handler) handleDeleteProductByID(w http.ResponseWriter, r *http.Request) {

	if auth.BridgeAdmin(w, r) {
		responseDeleteProductByID, err := deleteProductByID(r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(responseDeleteProductByID)
	} else {
		views_admin.Error().Render(r.Context(), w)

	}
}
