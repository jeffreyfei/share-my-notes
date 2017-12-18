package load_balancer

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffreyfei/share-my-notes/server/lib/router"
)

type provider struct {
	url     string
	idleInd int
}

type LoadBalancer struct {
	ClientRouter   *mux.Router
	ProviderRouter *mux.Router
	Providers      []provider
	nextProvider   int
	providerClient *http.Client
}

func NewLoadBalancer() *LoadBalancer {
	lb := new(LoadBalancer)
	lb.ClientRouter = router.BuildRouter(lb.buildClientRoutes())
	lb.ProviderRouter = router.BuildRouter(lb.buildProviderRoutes())
	lb.nextProvider = 0
	lb.providerClient = &http.Client{}
	lb.Providers = []provider{}
	return lb
}

func (lb *LoadBalancer) buildClientRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"GET",
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Home"))
			},
		},
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

func (lb *LoadBalancer) getNextProvider() (string, error) {
	if len(lb.Providers) == 0 {
		return "", errors.New("no available provider")
	}
	providerURL := lb.Providers[lb.nextProvider].url
	if lb.Providers[lb.nextProvider].idleInd != 0 {
		lb.Providers[lb.nextProvider].idleInd--
	} else {
		if lb.nextProvider+1 == len(lb.Providers) {
			lb.nextProvider = 0
		} else {
			lb.nextProvider++
		}
	}
	return providerURL, nil
}
