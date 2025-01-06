package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/shubhamojha1/heimdall/internal/config"
	"github.com/shubhamojha1/heimdall/internal/registry"

	"github.com/shubhamojha1/heimdall/internal/server"
)

// refer to server metrics needed by the lb. important

// move these structs to a common folder (/internal/server)
// then import here and in algorithms.go
// type ServiceRegistry struct {
// 	Backends       []*config.Backend
// 	HealthChecks   []*config.HealthCheck
// 	UpdateInterval time.Duration
// }

// type LoadBalancer struct {
// 	mu              sync.RWMutex
// 	Configuration   *config.Config
// 	LayerConfig     interface{} `json"-"`
// 	ServiceRegistry ServiceRegistry
// 	stopChan        chan struct{}
// 	listener        net.Listener
// }

// clear concepts wrt each L4 algorithm first.
// then start implementation
// func (lb *LoadBalancer) selectBackend() string{
// 	index :=
// }

// func (lb *LoadBalancer) handleConnectionRoundRobin

// func (lb *server.LoadBalancer) handleConnection(clientConn net.Conn) {
// 	defer clientConn.Close()

// 	if lb.Configuration.Algorithm == config.AlgorithmRoundRobin {
// 		backendAddrs := lb.selectBackend()
// 		algorithms.handleConnectionRoundRobin(lb, clientConn)
// 	}
// }

// func (lb *server.LoadBalancer) Start() error {
// 	log.Println("Starting Load Balancer...")
// 	log.Println("PORT: ", os.Getenv("LOAD_BALANCER_PORT"))
// 	// listen for new client requests

// 	for {
// 		clientConn, err := lb.listener.Accept()
// 		if err != nil {
// 			log.Printf("Error accepting connection: %v", err)
// 			continue
// 		}

// 		go lb.handleConnection(clientConn)
// 	}

// 	// go lb.runHealthChecks()

// 	return nil
// }

// func (lb *LoadBalancer) Stop() {
// 	lb.mu.Lock()
// 	defer lb.mu.Unlock()

// 	close(lb.stopChan)
// 	log.Println("Load Balancer stopped")
// }

// start load balancer first
// load the service registry along with load balancer
// start server manager
// add servers
// as servers are added, send server inital info to service registry
// to register the service
// add multiple servers as such
// when a new client comes in, it will talk only to the load balancer
// so load balancer is also a service registry
func NewServiceRegistry() (*registry.ServiceRegistry, error) {
	ServiceRegistryPort := os.Getenv("SERVICE_REGISTRY_PORT")
	if ServiceRegistryPort == "" {
		return nil, fmt.Errorf("SERVICE_REGISTRY_PORT environment variable not defined")
	}

	sr := &registry.ServiceRegistry{
		Backends:     make([]*config.Backend, 0),
		HealthChecks: make([]*config.HealthCheck, 0),
	}

	srMux := http.NewServeMux()
	srMux.HandleFunc("/", sr.HandleHTTP)
	srMux.HandleFunc("/register", sr.HandleRegister)
	srMux.HandleFunc("/heartbeat", sr.HandleHeartbeat)

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

	return sr, nil

}

func NewLoadBalancer(configuration *config.Config) (*server.LoadBalancer, error) {
	// layer := configuration.Layer
	// algo := configuration.Algorithm
	// // layerConfig := configuration.LayerConfig
	// lg := configuration.LayerConfig
	// fmt.Printf("%s %s \n\n %s", layer, lg, algo)

	// if algo == config.AlgorithmRoundRobin {
	// 	algorithms.round_robin.
	// }

	// algorithm := configuration.Algorithm
	// port := configuration.Listen.Port
	// address := configuration.Listen.Address

	// L4Settings := configuration.L4Settings
	// TCPKeepAlive := L4Settings.TCP.KeepAlive
	// TCPKeepAliveTime := L4Settings.TCP.KeepAliveTime * time.Second
	// MaxConnections := L4Settings.TCP.MaxConnections

	// var layerConfig interface{}
	// if layer == config.LayerFour {
	// 	layerConfig = config.L4Settings
	// }

	// switch settings := configuration.LayerConfig.(type) {
	// case config.L4Settings:
	// 	fmt.Printf("L4 Max Connections: %d\n", settings.MaxConnections)
	// case config.L7Settings:
	// 	fmt.Printf("L7 Max Header Size: %d\n", settings.MaxHeaderSize)
	// default:
	// 	fmt.Println("Unknown layer configuration")
	// }
	// load service registry first

	// ServiceRegistryPort := os.Getenv("SERVICE_REGISTRY_PORT")

	NewServiceRegistry, err := NewServiceRegistry()
	if err != nil {
		return nil, fmt.Errorf("Failed to create service registry: %v", err)
	}
	// listener, err := ser

	lb := &server.LoadBalancer{
		Configuration:   configuration,
		ServiceRegistry: NewServiceRegistry,
	}

	LoadBalancerPort := os.Getenv("LOAD_BALANCER_PORT")
	if LoadBalancerPort == "" {
		return nil, fmt.Errorf("LOAD_BALANCER_PORT environment variable not defined")
	}

	lbMux := http.NewServeMux()
	lbMux.HandleFunc("/", lb.HandleHTTP)
	go func() {
		log.Printf("Starting HTTP Load Balancer on port %s", LoadBalancerPort)
		if err := http.ListenAndServe(":"+LoadBalancerPort, lbMux); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	return lb, nil
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	configPath := os.Getenv("CONFIG_JSON_PATH")
	if configPath == "" {
		log.Fatalf("Config path environment variable not defined")
	}

	// file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}

	defer watcher.Close()

	var currentLoadBalancer *server.LoadBalancer
	applyConfig := func() error {
		// Stop existing load balancer if it exists
		if currentLoadBalancer != nil {
			currentLoadBalancer.Stop()
		}

		// Load new configuration
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		// Create new load balancer
		newLoadBalancer, err := NewLoadBalancer(cfg)
		if err != nil {
			return fmt.Errorf("failed to create load balancer: %v", err)
		}

		// Start the new load balancer
		if err := newLoadBalancer.Start(); err != nil {
			return fmt.Errorf("failed to start load balancer: %v", err)
		}

		// Update current load balancer
		currentLoadBalancer = newLoadBalancer
		return nil
	}

	if err := applyConfig(); err != nil {
		log.Fatalf("Initial configuration load failed: %v", err)
	}

	// add config file to watcher
	if err := watcher.Add(configPath); err != nil {
		log.Fatalf("Failed to watch config file: %v", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// only react to write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Config file changed, reloading...")

					time.Sleep(1 * time.Second)
					// apply new config
					if err := applyConfig(); err != nil {
						log.Printf("Failed to apply new configuration: %v", err)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Wacther error: %v", err)
			}
		}
	}()

	select {}

	// logger := log.New(os.Stdout, "[LoadBalancer]", log.LstdFlags|log.Lshortfile)

	// initial load balancer
	// loadBalancer, err := NewLoadBalancer(cfg)
	// if err != nil {
	// 	log.Fatalf("Failed to create load balancer: %v", err)
	// }
	// println(loadBalancer, "~~~~~~~~~")

	// err = watcher.Add(os.Getenv("CONFIG_JSON_PATH"))
	// if err != nil {
	// 	log.Fatalf("Failed to add file to watcher: %v", err)
	// }

	// event loop to 'watch' the file for changes
	// go func() {
	// 	for {
	// 		select {
	// 		case event := <-watcher.Events:
	// 			if event.Op&fsnotify.Write == fsnotify.Write {
	// 				log.Println("Config file changed, reloading...")

	// 				newCfg, err := config.LoadConfig(os.Getenv("CONFIG_JSON_PATH"))
	// 				if err != nil {
	// 					log.Fatalf("Failed to load config: %v", err)
	// 					continue
	// 				}

	// 				newLoadBalancer, err := NewLoadBalancer(newCfg)
	// 				if err != nil {
	// 					log.Fatalf("Failed to create load balancer: %v", err)
	// 					continue
	// 				}
	// 				// fmt.Println(lb, "~ignore~")
	// 				loadBalancer = newLoadBalancer
	// 			}
	// 		case err := <-watcher.Errors:
	// 			log.Printf("Watcher error: %v", err)
	// 		}
	// 	}
	// }()

	// go lb.StartMetricsServer()

	// go lb.StartHealthCheck()

	// go func() {
	// 	if err := lb.Start(); err != nil {
	// 		logger.Fatalf("Load balancer failed: %v", err)
	// 	}
	// }()

	// needs to be implemented separately
	// adminServer := admin.NewAdminServer(lb)
	// go func() {
	// 	logger.Printf("Admin server started on :8080")
	// 	// if err := http.ListenAndServe(":8080", adminServer); err != nil {
	// 	// 	logger.Fatalf("Admin server failed: %v", err)
	// 	// }
	// }()

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, s)
}
