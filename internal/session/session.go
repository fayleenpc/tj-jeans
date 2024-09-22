package session

import (
	"fmt"
	"log"
	"sync"
)

var (
	SessionStoreClient = NewSessionStoreClient("init")
	SessionStoreServer = NewSessionStoreServer()
	SessionStore       = NewClientServer()
)

type ClientServer struct {
	client *ssc
	server *sss
}

func NewClientServer() *ClientServer {
	return &ClientServer{
		client: NewSessionStoreClient("init"),
		server: NewSessionStoreServer(),
	}
}

type sss struct {
	server map[string]string
	mu     sync.RWMutex
}

func NewSessionStoreServer() *sss {
	v := &sss{
		server: make(map[string]string),
	}
	return v
}

func (s *sss) Set(v string) {

	s.mu.Lock() // Lock the mutex

	s.server["secretTokenCookie"] = v // Store the secret token
	s.mu.Unlock()                     // Ensure the mutex is unlocked
	log.Println("store secret token success!")

}

func (s *sss) Get() (string, bool) {

	s.mu.RLock() // Lock the mutex

	log.Println("retrieve secret token success!")
	v, ok := s.server["secretTokenCookie"]
	s.mu.RUnlock() // Ensure the mutex is unlocked
	return v, ok
}

type ssc struct {
	client       map[string]bool
	currentToken string
	mu           sync.RWMutex
}

type SessionClient interface {
	Set(string)
	Refresh(string)
	Blacklist() string
	Get() string
	String() string
}

func NewSessionStoreClient(token string) *ssc {
	v := &ssc{
		client:       make(map[string]bool),
		currentToken: token,
	}
	v.mu.Lock()
	v.client[token] = false
	v.mu.Unlock()

	return v
}

func (c *ssc) Set(token string) {

	c.mu.Lock()

	c.client[token] = true
	c.currentToken = token
	c.mu.Unlock()
	c.clean()

}

func (c *ssc) Refresh(token string) {
	c.Set(token)
}

func (c *ssc) Blacklist() string {
	c.mu.Lock()
	c.client[c.cursor()] = false
	// c.currentToken = ""
	c.mu.Unlock()
	c.clean()
	return c.cursor()
}

func (c *ssc) Get() string {

	c.mu.RLock()
	if !c.client[c.currentToken] {
		c.clean()
		return ""
	}
	c.mu.RUnlock()
	return c.currentToken
}

func (c *ssc) cursor() string {
	return c.currentToken
}

func (c *ssc) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return fmt.Sprintf("\n%v\n %+v\n", c.currentToken, c.client)
}

func (c *ssc) clean() {
	for k, _ := range c.client {

		if !c.client[k] {
			c.mu.Lock()
			delete(c.client, k)
			c.mu.Unlock()
		}

	}
}
