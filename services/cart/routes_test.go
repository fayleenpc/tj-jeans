package cart

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/gorilla/mux"
)

func TestProductsServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	tokenStore := &mockTokenStore{}
	productsStore := &mockProductsStore{}
	store := &mockOrderStore{}
	handler := NewHandler(store, productsStore, userStore, tokenStore)

	t.Run("should fail handle the cart/checkout", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/cart/checkout", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/cart/checkout", handler.handleCheckout)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
	t.Run("should correctly handle the cart/checkout", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/cart/checkout", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/cart/checkout", handler.handleCheckout)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockOrderStore struct{}

func (m *mockOrderStore) GetOrders() ([]types.Order, error) { return nil, nil }

func (m *mockOrderStore) GetOrdersByIDs(ordersIDs []int) ([]types.Order, error) { return nil, nil }

func (m *mockOrderStore) GetOrderByID(id int) (*types.Order, error) { return nil, nil }

func (m *mockOrderStore) CreateOrder(types.Order) (int64, error) { return 0, nil }

func (m *mockOrderStore) DeleteOrderByID(id int) (int64, error) { return 0, nil }

func (m *mockOrderStore) DeleteOrder(order types.Order) (int64, error) { return 0, nil }

func (m *mockOrderStore) UpdateOrder(order types.Order) (int64, error) { return 0, nil }

func (m *mockOrderStore) GetOrderItems() ([]types.OrderItem, error) { return nil, nil }

func (m *mockOrderStore) GetOrderItemsByIDs(ordersItemsIDs []int) ([]types.OrderItem, error) {
	return nil, nil
}

func (m *mockOrderStore) GetOrderItemsByID(int) (*types.OrderItem, error) { return nil, nil }

func (m *mockOrderStore) CreateOrderItem(types.OrderItem) error { return nil }

func (m *mockOrderStore) DeleteOrderItemByID(id int) (int64, error) { return 0, nil }

func (m *mockOrderStore) DeleteOrderItem(orderItem types.OrderItem) (int64, error) { return 0, nil }

func (m *mockOrderStore) UpdateOrderItem(orderItem types.OrderItem) (int64, error) { return 0, nil }

type mockProductsStore struct{}

func (m *mockProductsStore) GetProducts() ([]types.Product, error) {
	return []types.Product{}, nil
}
func (m *mockProductsStore) GetProductsByIDs([]int) ([]types.Product, error) {
	return []types.Product{}, nil
}
func (m *mockProductsStore) GetProductByID(id int) (*types.Product, error) { return nil, nil }
func (m *mockProductsStore) CreateProduct(types.Product) (int64, error) {
	return 0, nil
}
func (m *mockProductsStore) DeleteProductByID(id int) (int64, error) { return 0, nil }
func (m *mockProductsStore) DeleteProduct(types.Product) (int64, error) {
	return 0, nil
}
func (m *mockProductsStore) UpdateProduct(types.Product) (int64, error) {
	return 0, nil
}

type mockTokenStore struct{}

func (m *mockTokenStore) GetBlacklistedTokens() ([]types.Token, error) { return nil, nil }
func (m *mockTokenStore) CreateBlacklistTokens(types.Token) error {
	return nil
}
func (m *mockTokenStore) GetBlacklistTokenByString(string) (types.Token, error) {
	return types.Token{}, nil
}

type mockUserStore struct{}

func (m *mockUserStore) GetUsers() ([]types.User, error)                   { return nil, nil }
func (m *mockUserStore) GetUsersByIDs(userIDS []int) ([]types.User, error) { return nil, nil }
func (m *mockUserStore) UpdateVerifiedUserByEmail(string) error            { return nil }
func (m *mockUserStore) GetUserByEmail(string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}
func (m *mockUserStore) GetUserByID(id int) (*types.User, error)   { return nil, nil }
func (m *mockUserStore) DeleteUserByID(id int) (int64, error)      { return 0, nil }
func (m *mockUserStore) DeleteUser(user types.User) (int64, error) { return 0, nil }
func (m *mockUserStore) UpdateUser(user types.User) (int64, error) { return 0, nil }
func (m *mockUserStore) CreateUser(types.User) error               { return nil }
