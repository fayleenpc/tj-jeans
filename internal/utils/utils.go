package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func CraftJSON(method string, url string, m []byte, r *http.Request) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(m))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	// req = req.WithContext(r.Context())
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("Authorization-X", r.Header.Get("Authorization-X"))

	defer req.Body.Close()
	log.Printf("who's url [%v]\n", url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {

		return nil, http.StatusBadRequest, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	log.Printf("response body : %s\n", resBody)
	return resBody, res.StatusCode, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

// func AuthMiddlewareChain(middlewares ...http.HandlerFunc) {

// }
