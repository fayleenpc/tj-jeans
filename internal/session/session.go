package session

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	SessionStoreClient = NewSessionStoreClient("init")
	SessionStoreServer = NewSessionStoreServer()
)

func SetJWTAccessToken(w http.ResponseWriter, v string) {
	c := &http.Cookie{
		Name:     "access_token",
		Value:    v,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
}

func GetJWTAccessToken(r *http.Request) (string, error) {
	c, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func SetJWTSecretToken(w http.ResponseWriter, v string) {
	c := &http.Cookie{
		Name:     "secret_token",
		Value:    v,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
}

func GetJWTSecretToken(r *http.Request) (string, error) {
	c, err := r.Cookie("secret_token")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

type Sss struct {
	server chan map[string]string
	mu     sync.RWMutex
}

type SssServerStore interface {
	Set(string)
	Get() (string, bool)
}

func NewSessionStoreServer() SssServerStore {
	return &Sss{
		server: make(chan map[string]string, 1), // Use buffered channel
	}
}

func (s *Sss) Set(v string) {
	s.mu.Lock()
	defer s.mu.Unlock()                                   // Ensure the mutex is unlocked at the end
	s.server <- map[string]string{"secretTokenCookie": v} // Store the secret token
}

func (s *Sss) Get() (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock() // Ensure the mutex is unlocked at the end

	select {
	case v := <-s.server:
		if v["secretTokenCookie"] != "" {
			return v["secretTokenCookie"], true
		} else {
			return "", false
		}

	default:
		return "", false // No value available
	}
}

type Ssc struct {
	client       chan map[string]bool
	currentToken string
	mu           sync.RWMutex
}

type SscClientStore interface {
	Set(string)
	Refresh(string)
	Blacklist() string
	Get() string
	String() string
}

func NewSessionStoreClient(token string) SscClientStore {
	return &Ssc{
		client:       make(chan map[string]bool, 1), // Use buffered channel
		currentToken: token,
	}
}

func (c *Ssc) Set(token string) {
	c.mu.Lock()
	defer c.mu.Unlock() // Ensure the mutex is unlocked at the end
	c.clean()
	c.client <- map[string]bool{token: true}
	c.currentToken = token

}

func (c *Ssc) Refresh(token string) {
	c.Set(token)
}

func (c *Ssc) Blacklist() string {
	c.mu.Lock()
	c.clean()
	defer c.mu.Unlock() // Ensure the mutex is unlocked at the end

	select {
	case v := <-c.client:
		v[c.currentToken] = false
	default:
	}

	return c.currentToken
}

func (c *Ssc) Get() string {
	c.mu.Lock()
	c.clean()
	defer c.mu.Unlock() // Ensure the mutex is unlocked at the end

	select {
	case v := <-c.client:
		if !v[c.currentToken] {

			return ""
		}
	default:
		return c.currentToken
	}
	return c.currentToken
}

func (c *Ssc) clean() {

	select {
	case take := <-c.client:
		for k, v := range take {
			if !v {
				delete(take, k)
			}
		}
		// If needed, you can reassign to client here to clear it out
		c.client <- take
	default:
	}
}

func (c *Ssc) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock() // Ensure the mutex is unlocked at the end
	return fmt.Sprintf("\n%v\n %+v\n", c.currentToken, c.client)
}
