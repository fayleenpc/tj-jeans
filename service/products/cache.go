package products

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/types"
	"github.com/fayleenpc/tj-jeans/utils"
	"github.com/nats-io/nats.go"
)

func PublishForProducts(ps []types.Product) {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Example slice of products
	products := ps

	// Serialize products to JSON
	productData, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
	}

	// Publish product data to the "products" subject
	err = nc.Publish("products", productData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Product data published")
}

func SubscribeForProduct(w http.ResponseWriter, r *http.Request) {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	defer nc.Close()

	// Subscribe to the "products" subject
	nc.Subscribe("products", func(m *nats.Msg) {
		var products []types.Product
		err := json.Unmarshal(m.Data, &products)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
		}

		// Cache the product data (e.g., store it in memory or a database)
		log.Printf("Received product data: %+v", products)
		utils.WriteJSON(w, http.StatusOK, products)
	})

	// Keep the connection open
	select {}
}
