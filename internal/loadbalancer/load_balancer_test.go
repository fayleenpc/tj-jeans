package loadbalancer

import (
	"testing"
)

func TestLoadBalancer(t *testing.T) {
	lb := NewLoadBalancer()
	t.Log(lb.GetBackend())
}
