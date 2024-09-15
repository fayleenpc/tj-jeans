package landing

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/config"
	"github.com/fayleenpc/tj-jeans/types"
	"github.com/fayleenpc/tj-jeans/views"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Home().Render(r.Context(), w)
	})

	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		var payload []types.Product
		resp, err := http.Get(config.Envs.PublicHost + ":" + config.Envs.Port + "/api/v1/products")
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			log.Fatal(err)
		}
		views.Products(payload).Render(r.Context(), w)
	})
	router.HandleFunc("/gallery", func(w http.ResponseWriter, r *http.Request) {
		views.Gallery().Render(r.Context(), w)
	})
	router.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		views.Login_Register().Render(r.Context(), w)
	})
	// router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
	// 	// authentication for only 1 user in db , role is unique id and password as admin itself

	// })
}
