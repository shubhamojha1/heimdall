package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/shubhamojha1/heimdall/internal/config"
	"github.com/shubhamojha1/heimdall/internal/metrics"
)

type ServiceRegistry struct {
	Backends     []*config.Backend
	HealthChecks []*config.HealthCheck
	listener     http.Server // for getting info from backends
	mu           sync.Mutex
	// UpdateInterval time.Duration
}

func (sr *ServiceRegistry) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling HTTP request...")

	w.Write([]byte("Hello from service registry!"))
}

func (sr *ServiceRegistry) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var backend config.Backend
	if err := json.NewDecoder(r.Body).Decode(&backend); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	sr.mu.Lock()
	sr.Backends = append(sr.Backends, &backend)
	sr.mu.Unlock()

	log.Printf("Registered backend: %v", backend)
	w.WriteHeader(http.StatusOK)
}

func (sr *ServiceRegistry) HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var heartbeat struct {
		Port      int                    `json:"port"`
		URL       string                 `json:"url"`
		Status    string                 `json:"status"`
		StartedAt time.Time              `json:"started_at"`
		Metrics   *metrics.ServerMetrics `json:"metrics"`
	}
	if err := json.NewDecoder(r.Body).Decode(&heartbeat); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	sr.mu.Lock()
	for _, backend := range sr.Backends {
		if backend.URL == heartbeat.URL {
			// backend.Status = heartbeat.Status
			// Update other fields as needed
			break
		}
	}
	sr.mu.Unlock()

	log.Printf("Received heartbeat from backend: %v", heartbeat)
	w.WriteHeader(http.StatusOK)
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
