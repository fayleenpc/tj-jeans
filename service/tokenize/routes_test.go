package tokenize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/gorilla/mux"
)

func TestTokenizeServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	store := &mockTokenStore{}
	handler := NewHandler(store, userStore)

	t.Run("should fail if the token payload is invalid", func(t *testing.T) {
		payload := types.Token{
			Token: "test",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodGet, "/refresh", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/refresh", handler.handleRefresh)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
	t.Run("should correctly refresh the token", func(t *testing.T) {
		payload := types.Token{
			Token: "secret/access token",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodGet, "/refresh", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/refresh", handler.handleRefresh)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockTokenStore struct{}

func (m *mockTokenStore) CreateBlacklistTokens(types.Token) error {
	return nil
}
func (m *mockTokenStore) GetBlacklistTokenByString(string) (types.Token, error) {
	return types.Token{}, nil
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}
func (m *mockUserStore) GetUserByID(id int) (*types.User, error) { return nil, nil }
func (m *mockUserStore) CreateUser(types.User) error             { return nil }
func (m *mockUserStore) UpdateVerifiedUserByEmail(string) error  { return nil }
