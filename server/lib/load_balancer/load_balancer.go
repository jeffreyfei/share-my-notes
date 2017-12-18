package load_balancer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jeffreyfei/share-my-notes/server/lib/router"
	log "github.com/sirupsen/logrus"
)

type provider struct {
	url     string
	idleInd int
}

type LoadBalancer struct {
	ClientRouter        *mux.Router
	ProviderRouter      *mux.Router
	Providers           []provider
	nextProvider        int
	healthCheckInterval int
	providerClient      *http.Client
}

func NewLoadBalancer(heatlCheckInterval int) *LoadBalancer {
	lb := new(LoadBalancer)
	lb.ClientRouter = router.BuildRouter(lb.buildClientRoutes())
	lb.ProviderRouter = router.BuildRouter(lb.buildProviderRoutes())
	lb.nextProvider = 0
	lb.providerClient = &http.Client{}
	lb.Providers = []provider{}
	lb.healthCheckInterval = heatlCheckInterval
	return lb
}

func (lb *LoadBalancer) StartHealthCheck() {
	go func() {
		for {
			log.Info("Performing health check.")
			lb.healthCheck()
			log.Info("Health check completed.")
			time.Sleep(time.Duration(lb.healthCheckInterval) * time.Millisecond)
		}
	}()
}

// Check if each provider is online, if offline remove provider from providerl ist
// Fetch current job count on each server. The values are normalized and computed to idle indexes.
func (lb *LoadBalancer) healthCheck() {
	maxJobCount := 0
	providerJobs := make(map[string]int)
	for i, provider := range lb.Providers {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/report/status", provider.url), nil)
		if err != nil {
			log.WithField("err", err).Error("Failed to create health check request.")
			continue
		}
		rec, err := lb.providerClient.Do(req)
		if err != nil {
			log.WithField("err", err).Errorf("Failed to contact provider. Removing provider %s from list", provider.url)
			lb.Providers = append(lb.Providers[:i], lb.Providers[i+1:]...)
			continue
		}
		body, err := ioutil.ReadAll(rec.Body)
		if err != nil {
			log.WithField("err", err).Error("Failed to read provider response body.")
			continue
		}
		jobCount, err := strconv.Atoi(string(body))
		if err != nil {
			log.WithField("err", err).Error("Invalid body content.")
			continue
		}
		providerJobs[provider.url] = jobCount
		if jobCount > maxJobCount {
			maxJobCount = jobCount
		}
	}
	lb.computeIdleIndexes(maxJobCount, providerJobs)
}

func (lb *LoadBalancer) computeIdleIndexes(maxJobCount int, providerJobs map[string]int) {
	for i, provider := range lb.Providers {
		lb.Providers[i].idleInd = maxJobCount - providerJobs[provider.url]
	}
}

func (lb *LoadBalancer) hasProvider(url string) bool {
	for _, provider := range lb.Providers {
		if provider.url == url {
			return true
		}
	}
	return false
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
