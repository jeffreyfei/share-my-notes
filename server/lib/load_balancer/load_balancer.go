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
	lb.ProviderRouter = router.BuildRouter(lb.buildProviderRoutes())
	lb.nextProvider = 0
	lb.providerClient = &http.Client{}
	lb.Providers = make(map[string]int)
	return lb
}

func (lb *LoadBalancer) buildClientRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"GET",
			"/md/{action}/{id}",
			lb.mdClientHandler,
		},
		router.Route{
			"POST",
			"/md/{action}/{id}",
			lb.mdClientHandler,
		},
		router.Route{
			"POST",
			"/md/{action}",
			lb.mdClientHandler,
		},
		router.Route{
			"GET",
			"/auth/google/{action}",
			lb.googleAuthHandler,
		},
	}
}

func (lb *LoadBalancer) buildProviderRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"POST",
			"/provider/register",
			lb.providerRegisterHandler,
		},
	}
}

func (lb *LoadBalancer) getNextProvider() string {
	return ""
}
