package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/shubhamojha1/heimdall/internal/config"

	"github.com/joho/godotenv"
)

// refer to server metrics needed by the lb. important
func main() {
	err := godotenv.Load(filepath.Join("../..", ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	_, err = config.LoadConfig(os.Getenv("CONFIG_JSON_PATH"))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger := log.New(os.Stdout, "[LoadBalancer]", log.LstdFlags|log.Lshortfile)
	// fmt.Printf(string(cfg.Algorithm.AlgorithmRoundRobin))
	// lb, err := lb.NewLoadBalancer(cfg, logger)
	// if err != nil {
	// 	log.Fatalf("Failed to create load balancer: %v", err)
	// }

	// go lb.StartMetricsServer()

	// go lb.StartHealthCheck()

	// go func() {
	// 	if err := lb.Start(); err != nil {
	// 		logger.Fatalf("Load balancer failed: %v", err)
	// 	}
	// }()

	// needs to be implemented separately
	// adminServer := admin.NewAdminServer(lb)
	go func() {
		logger.Printf("Admin server started on :8080")
		// if err := http.ListenAndServe(":8080", adminServer); err != nil {
		// 	logger.Fatalf("Admin server failed: %v", err)
		// }
	}()

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, s)
}
