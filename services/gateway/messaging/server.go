package messaging

import (
	"github.com/nats-io/nats.go"
)

// PublishTokenEvent publishes a message to NATS
func PublishTokenEvent(nc *nats.Conn, token string) error {
	return nc.Publish("secretTokenCookie", []byte(token))
}

func Connect() (*nats.Conn, error) {
	// NATS integration
	nc, err := nats.Connect(nats.DefaultURL)
	return nc, err
}
