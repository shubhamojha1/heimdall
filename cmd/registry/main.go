package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

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

func NewServiceRegistry() (*ServiceRegistry, error) {
	ServiceRegistryPort := os.Getenv("SERVICE_REGISTRY_PORT")
	if ServiceRegistryPort == "" {
		return &ServiceRegistry{}, fmt.Errorf("SERVICE_REGISTRY_PORT environment variable not defined")
	}

	sr := &ServiceRegistry{
		Backends:     make([]*config.Backend, 0),
		HealthChecks: make([]*config.HealthCheck, 0),
	}

	srMux := http.NewServeMux()
	srMux.HandleFunc("/", sr.HandleHTTP)
	srMux.HandleFunc("/register", sr.HandleRegister)
	srMux.HandleFunc("/heartbeat", sr.HandleHeartbeat)
	srMux.HandleFunc("/remove", sr.HandleRemove)

	// grpcServer := grpc.NewServer()
	// pb.RegisterServiceRegistryServer(grpcServer, sr)

	go func() {
		log.Printf("Starting HTTP Service Registry on port %s", ServiceRegistryPort)
		if err := http.ListenAndServe(":"+ServiceRegistryPort, srMux); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// start gRPC server concurrently
	// go func() {
	// 	log.Printf("Starting gRPC Service Registry on port %s", ServiceRegistryPort)
	// 	if err := http.ListenAndServe(":"+ServiceRegistryPort, srMux); err != nil {
	// 		log.Fatalf("Failed to start HTTP server: %v", err)
	// 	}
	// }()
	fmt.Println("Service Registry started: ", time.Now().String())

	return sr, nil

}

func (sr *ServiceRegistry) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[Registry] Handling HTTP request...")

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

func (sr *ServiceRegistry) HandleRemove(w http.ResponseWriter, r *http.Request) {
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
	for i, b := range sr.Backends {
		if b.URL == backend.URL {
			sr.Backends = append(sr.Backends[:i], sr.Backends[i+1:]...)
			break
		}
	}
	sr.mu.Unlock()

	log.Printf("Removed backend: %v", backend)
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

//	func (sr *ServiceRegistry) AddBackend(backend *config.Backend) {
//		sr.Backends = append(sr.Backends, backend)
//	}
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	_, err := NewServiceRegistry()
	if err != nil {
		log.Fatalf("Failed to start service registry: %v", err)
	}

	select {}
	// block main goroutine to keep the servers running
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
