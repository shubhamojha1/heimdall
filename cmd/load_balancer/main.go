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

	"github.com/shubhamojha1/heimdall/internal/server"
	"google.golang.org/grpc"
)

// refer to server metrics needed by the lb. important

// move these structs to a common folder (/internal/server)
// then import here and in algorithms.go
// type ServiceRegistry struct {
// 	Backends       []*config.Backend
// 	HealthChecks   []*config.HealthCheck
// 	UpdateInterval time.Duration
// }

func NewLoadBalancer(configuration *config.Config) (*server.LoadBalancer, error) {

	lb := &server.LoadBalancer{
		Configuration: configuration,
		// ServiceRegistry: NewServiceRegistry,
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

	fmt.Println("Load Balancer started: ", time.Now().String())

	return lb, nil
}

func SendHelloToServiceRegistry()

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

	// notify service registry about new load balancer
	// grpc
	resp, err := grpc.NewClient()

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
