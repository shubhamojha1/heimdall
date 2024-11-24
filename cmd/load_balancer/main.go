package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	Configuration   *config.Config
	LayerConfig     interface{} `json"-"`
	ServiceRegistry ServiceRegistry
}

func NewLoadBalancer(configuration *config.Config) (*LoadBalancer, error) {
	layer := configuration.Layer
	// algo := configuration.Algorithm
	// layerConfig := configuration.LayerConfig
	lg := configuration.LayerConfig
	fmt.Printf("%s %s", layer, lg)
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg, err := config.LoadConfig(os.Getenv("CONFIG_JSON_PATH"))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// logger := log.New(os.Stdout, "[LoadBalancer]", log.LstdFlags|log.Lshortfile)

	// lb, err := lb.NewLoadBalancer(cfg, logger)
	lb, err := NewLoadBalancer(cfg)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	fmt.Println(lb, "~ignore~")

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
