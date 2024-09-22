package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	servers []*url.URL
	index   int
	mu      sync.Mutex
}

func NewLoadBalancer(servers []string) (*LoadBalancer, error) {
	var urls []*url.URL
	for _, server := range servers {
		url, err := url.Parse(server)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return &LoadBalancer{
		servers: urls,
	}, nil
}

func (lb *LoadBalancer) getNextServer() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers)
	return server
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.getNextServer()
	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.ServeHTTP(w, r)
}

func Load() {
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		// Add more backend servers as needed
	}

	lb, err := NewLoadBalancer(servers)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8080", lb)
}

// backend 1
func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from server 1!")
}

func service1() {
	http.HandleFunc("/", handler1)
	http.ListenAndServe(":8081", nil)
}

// backend 2
func handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from server 2!")
}

func service2() {
	http.HandleFunc("/", handler2)
	http.ListenAndServe(":8082", nil)
}
