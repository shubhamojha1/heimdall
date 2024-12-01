package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/shubhamojha1/heimdall/internal/config"
)

// refer to server metrics needed by the lb. important
type ServiceRegistry struct {
	Backends       []*config.Backend
	HealthChecks   []*config.HealthCheck
	UpdateInterval time.Duration
}

type LoadBalancer struct {
	mu              sync.RWMutex
	Configuration   *config.Config
	LayerConfig     interface{} `json"-"`
	ServiceRegistry ServiceRegistry
	stopChan        chan struct{}
}

func (lb *LoadBalancer) Start() error {
	log.Println("Starting Load Balancer...")

	// go lb.runHealthChecks()

	return nil
}

func (lb *LoadBalancer) Stop() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	close(lb.stopChan)
	log.Println("Load Balancer stopped")
}

func NewLoadBalancer(configuration *config.Config) (*LoadBalancer, error) {
	layer := configuration.Layer
	algo := configuration.Algorithm
	// layerConfig := configuration.LayerConfig
	lg := configuration.LayerConfig
	fmt.Printf("%s %s \n\n %s", layer, lg, algo)
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
	return nil, nil
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

	var currentLoadBalancer *LoadBalancer
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
