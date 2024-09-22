package landing

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/session"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/platform/web/views"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.TokenStore
}

func NewHandler(store types.TokenStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	tokenChan := make(chan string)
	// need user is login or not
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			if session.SessionStoreClient.Get() != "" {
				tokenChan <- session.SessionStoreClient.Get()
			}
		}()
		select {
		case t := <-tokenChan:

			token, err := h.store.GetBlacklistTokenByString(t)
			if err == nil {
				log.Printf("blacklisted token %v, path : %v", token, r.URL.Path)
				views.Home("").Render(r.Context(), w)
			} else {
				log.Printf("got token %v, path : %v", token, r.URL.Path)
				views.Home(auth.GetUserNameFromSession(token.Token)).Render(r.Context(), w)

			}

		default:
			log.Printf("no token, path %v", r.URL.Path)
			views.Home("").Render(r.Context(), w)
		}
	})

	// need user is login or not
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

		go func() {
			if session.SessionStoreClient.Get() != "" {
				tokenChan <- session.SessionStoreClient.Get()
			}
		}()
		select {
		case t := <-tokenChan:
			token, err := h.store.GetBlacklistTokenByString(t)
			if err == nil {
				log.Printf("blacklisted token %v, path : %v", token, r.URL.Path)
				views.Products(payload, "").Render(r.Context(), w)

			} else {
				log.Printf("got token %v, path : %v", token, r.URL.Path)
				views.Products(payload, auth.GetUserNameFromSession(token.Token)).Render(r.Context(), w)

			}

		default:
			log.Printf("no token, path %v", r.URL.Path)
			views.Products(payload, "").Render(r.Context(), w)
		}

	})

	// need user is login or not
	router.HandleFunc("/gallery", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			if session.SessionStoreClient.Get() != "" {
				tokenChan <- session.SessionStoreClient.Get()
			}
		}()
		select {
		case t := <-tokenChan:
			token, err := h.store.GetBlacklistTokenByString(t)
			if err == nil {
				log.Printf("blacklisted token %v, path : %v", token, r.URL.Path)
				views.Gallery("").Render(r.Context(), w)
			} else {
				log.Printf("got token %v, path : %v", token, r.URL.Path)
				views.Gallery(auth.GetUserNameFromSession(token.Token)).Render(r.Context(), w)
			}

		default:
			log.Printf("no token, path %v", r.URL.Path)
			views.Gallery("").Render(r.Context(), w)
		}

	})

	// need user is login or not
	router.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			if session.SessionStoreClient.Get() != "" {
				tokenChan <- session.SessionStoreClient.Get()
			}
		}()
		select {
		case t := <-tokenChan:
			token, err := h.store.GetBlacklistTokenByString(t)
			if err == nil {
				log.Printf("blacklisted token %v, path : %v", token, r.URL.Path)
				views.Login_Register("").Render(r.Context(), w)
			} else {
				log.Printf("got token %v, path : %v", token, r.URL.Path)
				views.Login_Register(auth.GetUserNameFromSession(token.Token)).Render(r.Context(), w)
			}

		default:
			log.Printf("no token, path %v", r.URL.Path)
			views.Login_Register("").Render(r.Context(), w)
		}

	})

	// router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
	// 	// authentication for only 1 user in db , role is unique id and password as admin itself

	// })
}
