package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jeffreyfei/share-my-notes/server/lib/load_balancer"

	log "github.com/sirupsen/logrus"
)

// Bootstraps and runs the load balancer

var (
	loadBalancer *load_balancer.LoadBalancer
)

// Initializes the load balancer instance
func initLoadBalancer() {
	loadBalancer = load_balancer.NewLoadBalancer(30000)
	loadBalancer.StartHealthCheck()
}

// Bootstrapping
func main() {
	initLoadBalancer()
	client_port := fmt.Sprintf(":%s", os.Getenv("CLIENT_PORT"))
	provider_port := fmt.Sprintf(":%s", os.Getenv("PROVIDER_PORT"))
	go http.ListenAndServe(client_port, loadBalancer.ClientRouter)
	go http.ListenAndServe(provider_port, loadBalancer.ProviderRouter)
	log.Infof("Load balancer running on port: %s (client), %s (provider)", client_port, provider_port)
	fmt.Scanln()
}
