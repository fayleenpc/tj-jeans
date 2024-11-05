package loadbalancer

import (
	"math/rand"
	"sync"
)

type Backend struct {
	Address string
	Weight  int
}

type LoadBalancer struct {
	backends    []Backend
	totalWeight int
	mu          sync.Mutex
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		backends: []Backend{
			{Address: "http://localhost:8081", Weight: 1},
			{Address: "http://localhost:8082", Weight: 2},
			// {Address: "http://localhost:8083", Weight: 1},
		},
	}
}

func (lb *LoadBalancer) GetBackend() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Calculate total weight
	if lb.totalWeight == 0 {
		for _, backend := range lb.backends {
			lb.totalWeight += backend.Weight
		}
	}

	r := rand.Intn(lb.totalWeight)
	for _, backend := range lb.backends {
		if r < backend.Weight {
			return backend.Address
		}
		r -= backend.Weight
	}

	return "" // Should never reach here
}
