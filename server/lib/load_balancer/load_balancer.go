package load_balancer

import (
	"github.com/gorilla/mux"
)

type LoadBalancer struct {
	ClientRouter   *mux.Router
	ProviderRouter *mux.Router
	Providers      []string
}

func NewLoadBalancer() *LoadBalancer {
	lb := new(LoadBalancer)
	lb.ClientRouter = lb.buildRouter(lb.buildClientRoutes())
	lb.ClientRouter = lb.buildRouter(lb.buildProviderRoutes())
	return lb
}
