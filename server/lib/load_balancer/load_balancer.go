package load_balancer

import (
	"github.com/gorilla/mux"
	"github.com/jeffreyfei/share-my-notes/server/lib/router"
)

type LoadBalancer struct {
	ClientRouter   *mux.Router
	ProviderRouter *mux.Router
	Providers      map[string]int
}

func NewLoadBalancer() *LoadBalancer {
	lb := new(LoadBalancer)
	lb.ClientRouter = router.BuildRouter(lb.buildClientRoutes())
	lb.ClientRouter = router.BuildRouter(lb.buildProviderRoutes())
	return lb
}

func (l *LoadBalancer) buildClientRoutes() router.Routes {
	return router.Routes{}
}

func (l *LoadBalancer) buildProviderRoutes() router.Routes {
	return router.Routes{}
}
