package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/shubhamojha1/heimdall/internal/config"
)

// type ServiceRegistry struct {
// 	Backends       []*config.Backend
// 	HealthChecks   []*config.HealthCheck
// 	UpdateInterval time.Duration
// }

type LoadBalancer struct {
	mu            sync.RWMutex
	Configuration *config.Config
	// LayerConfig     interface{} `json"-"`
	// ServiceRegistry *registry.ServiceRegistry (implementing as a separate process)
	stopChan chan struct{}
	listener http.Server // for clients outside the network to connect to a server via the load balancer
}

func (lb *LoadBalancer) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// if lb.Configuration.Algorithm == config.AlgorithmRoundRobin {
	// 	backendAddrs := lb.selectBackend()
	// 	algorithms.handleConnectionRoundRobin(lb, clientConn)
	// }
	// return nil
}

func (lb *LoadBalancer) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[LB Server] Handling HTTP request...")

	w.Write([]byte("Hello from load balancer!"))
}

func (lb *LoadBalancer) Start() error {
	log.Println("Starting Load Balancer...")
	log.Println("PORT: ", os.Getenv("LOAD_BALANCER_PORT"))
	// add additional checks later
	// listen for new client requests

	// for {
	// 	clientConn, err := lb.listener.Accept()
	// 	if err != nil {
	// 		log.Printf("Error accepting connection: %v", err)
	// 		continue
	// 	}

	// 	go lb.handleConnection(clientConn)
	// }

	// go lb.runHealthChecks()

	return nil
}

func (lb *LoadBalancer) Stop() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	close(lb.stopChan)
	log.Println("Load Balancer stopped")
}
