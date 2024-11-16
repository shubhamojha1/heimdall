package loadbalancer

import (
	"log"

	"github.com/shubhamojha1/heimdall/internal/config"
)

type LoadBalancer struct {
	config      *config.Config
	algorithm   *config.Algorithm
	backends    []*config.Backend
	healthcheck *config.HealthCheck
	logger      *log.Logger
}
