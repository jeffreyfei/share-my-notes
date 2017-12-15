package load_balancer

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffreyfei/share-my-notes/server/lib/router"
)

type LoadBalancer struct {
	ClientRouter   *mux.Router
	ProviderRouter *mux.Router
	Providers      map[string]int
	nextProvider   int
	providerClient *http.Client
}

func NewLoadBalancer() *LoadBalancer {
	lb := new(LoadBalancer)
	lb.ClientRouter = router.BuildRouter(lb.buildClientRoutes())
	lb.ClientRouter = router.BuildRouter(lb.buildProviderRoutes())
	lb.nextProvider = 0
	lb.providerClient = &http.Client{}
	return lb
}

func (l *LoadBalancer) buildClientRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"GET",
			"/md/{action}/{id}",
			l.mdClientHandler,
		},
		router.Route{
			"POST",
			"/md/{action}/{id}",
			l.mdClientHandler,
		},
		router.Route{
			"POST",
			"/md/{action}",
			l.mdClientHandler,
		},
		router.Route{
			"GET",
			"/auth/google/{action}",
			l.googleAuthHandler,
		},
	}
}

func (l *LoadBalancer) buildProviderRoutes() router.Routes {
	return router.Routes{}
}

func (l *LoadBalancer) getNextProvider() string {
	return ""
}
