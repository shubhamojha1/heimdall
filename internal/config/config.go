package config

import (
	"encoding/json"
	"os"
)

type Mode string

const (
	ModeStatic  Mode = "static"
	ModeDynamic Mode = "dynamic"
)

type Algorithm string

const (
	AlgorithmRoundRobin         Algorithm = "round_robin"
	AlgorithmStickyRoundRobin   Algorithm = "sticky_round_robin"
	AlgorithmWeightedRoundRobin Algorithm = "weighted_round_robin"
	AlgorithmIPHash             Algorithm = "ip_hash"
	AlgorithmURLHash            Algorithm = "url_hash"
	AlgorithmLeastConnections   Algorithm = "least_connections"
	AlgorithmLeastTime          Algorithm = "least_time"
)

type Config struct {
	Mode      Mode      `json:"mode"`
	Algorithm Algorithm `json:"algorithm"`

	Listen struct {
		Port    int `json:"port"`
		Address int `json:"address"`
	} `json:"listen"`

	Backends []Backend `json:"backends"`

	HealthCheck struct {
		Interval time.Duration `json:"interval"`
		Timeout  time.Duration `json:"timeout"`
		Path     string        `json:"path"`
	} `json:"healthcheck"`

	Metrics struct {
		Enabled bool `json:"enabled"`
		Port    int  `json:"port"`
	} `json:"metrics"`
}

type Backend struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Weight   int    `json:"weight"`
	MaxConns int    `json:"maxConns"`
	// also add queue full boolean
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
