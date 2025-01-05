package registry

import (
	"log"
	"net/http"

	"github.com/shubhamojha1/heimdall/internal/config"
)

type ServiceRegistry struct {
	Backends     []*config.Backend
	HealthChecks []*config.HealthCheck
	listener     http.Server // for getting info from backends
	// UpdateInterval time.Duration
}

func (sr *ServiceRegistry) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling HTTP request...")

	w.Write([]byte("Hello from service registry!"))
}

func (sr *ServiceRegistry) AddBackend(backend *config.Backend) {
	sr.Backends = append(sr.Backends, backend)
}

// start load balancer first
// load the service registry along with load balancer
// start server manager
// add servers
// as servers are added, send server inital info to service registry
// to register the service
// add multiple servers as such
// when a new client comes in, it will talk only to the load balancer
// so load balancer is also a service registry
