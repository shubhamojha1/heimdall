package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/shubhamojha1/heimdall/internal/metrics"
)

type Layer string

const (
	LayerFour  Layer = "l4"
	LayerSeven Layer = "l7"
)

type Algorithm string

const (
	// L4 algorithms
	AlgorithmRoundRobin         Algorithm = "round_robin"
	AlgorithmWeightedRoundRobin Algorithm = "weighted_round_robin"
	AlgorithmLeastConnections   Algorithm = "least_connections"
	AlgorithmStickyRoundRobin   Algorithm = "sticky_round_robin"

	// L7 algorithms
	AlgorithmURLHash      Algorithm = "url_hash"
	AlgorithmCookieBased  Algorithm = "cookie_based"
	AlgorithmContentBased Algorithm = "content_based"

	// Both layers
	AlgorithmIPHash    Algorithm = "ip_hash"
	AlgorithmLeastTime Algorithm = "least_time"
)

func (a Algorithm) IsValidForLayer(l Layer) bool {
	switch l {
	case LayerFour:
		switch a {
		case AlgorithmRoundRobin, AlgorithmWeightedRoundRobin, AlgorithmStickyRoundRobin, AlgorithmLeastConnections, AlgorithmIPHash, AlgorithmLeastTime:
			return true
		}

	case LayerSeven:
		switch a {
		case AlgorithmURLHash, AlgorithmCookieBased, AlgorithmContentBased, AlgorithmIPHash, AlgorithmLeastTime:
			return true
		}
	}
	return false
}

type Backend struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Weight int    `json:"weight"`

	// Connection limits
	MaxConns          int           `json:"max_connections"`
	MaxQueueSize      int           `json:"max_queue_size"`
	QueueTimeout      time.Duration `json:"queue_timeout"`
	ConnectionTimeout time.Duration `json:"connection_timeout"`

	// Server capacity settings
	MaxCPUUsage     float64       `json:"max_cpu_usage"`    // percentage threshold
	MaxMemoryUsage  float64       `json:"max_memory_usage"` // percentage threshold
	MaxResponseTime time.Duration `json:"max_response_time"`

	// Current state
	QueueFull       bool `json:"queue_full"`
	Enabled         bool `json:"enabled"`
	MaintenanceMode bool `json:"maintenance_mode"`

	SSL struct {
		Enabled  bool   `json:"enabled"`
		Verify   bool   `json:"verify"`
		CertPath string `json:"cert_path,omitempty"`
		KeyPath  string `json:"key_path,omitempty"`
	} `json:"ssl,omitempty"`

	// Current metrics
	Metrics metrics.ServerMetrics `json:"metrics"`
}

// main configuration structure for the load balancer
type Config struct {
	Layer     Layer     `json:"layer"`
	Algorithm Algorithm `json:"algorithm"`

	Listen struct {
		Port    int    `json:"port"`
		Address string `json:"address"`
	} `json:"listen"`

	// remove backends[] as it will be managed and stored dynamically.
	// Backends []Backend `json:"backends"`

	// L4Settings *L4Settings `json:"l4_settings,omitempty"`
	// L7Settings *L7Settings `json:"l7_settings,omitempty"`
	LayerConfig interface{} `json:"-"`

	HealthCheck HealthCheck `json:"healthcheck"`

	Metrics struct {
		Enabled bool `json:"enabled"`
		Port    int  `json:"port"`
	} `json:"metrics"`
}

type HealthCheck struct {
	Enabled  bool          `json:"enabled"`
	Protocol string        `json:"protocol"` // "tcp" for L4, "http" for L7
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
	Path     string        `json:"path"`
	Port     int           `json:"port"`
	Expected string        `json:"expected_status,omitempty"`
}

type L4Settings struct {
	TCP struct {
		KeepAlive      bool          `json:"keepalive"`
		KeepAliveTime  time.Duration `json:"keepalive_time"`
		MaxConnections int           `json:"max_connections"`
		// ConnectionTimeout	time.Duration	`json:connection_timeout,omitempty`
		// IdleTimeout		time.Duration		`json:idle_timeout,omitempty`
	} `json:"tcp"`
}

type L7Settings struct {
	HTTP struct {
		MaxHeaderSize   int               `json:"max_header_size"`
		IdleTimeout     time.Duration     `json:"idle_timeout"`
		WriteTimeout    time.Duration     `json:"write_timeout"`
		EnableHTTPS     bool              `json:"enable_https"`
		SSLCertPath     string            `json:"ssl_cert_path,omitempty"`
		SSLKeyPath      string            `json:"ssl_key_path,omitempty"`
		RequestHeaders  map[string]string `json:"request_headers"`
		ResponseHeaders map[string]string `json:"response_headers"`
	} `json:"http"`

	Sticky struct {
		Enabled    bool   `json:"enabled"`
		CookieName string `json:"cookie_name,omitempty"`
		HashMethod string `json:"hash_method,omitempty"`
		HashKey    string `json:"hash_key,omitempty"`
	} `json:"sticky"`

	// Monitoring and metrics configuration
	Monitoring struct {
		Enabled          bool          `json:"enabled"`
		UpdateInterval   time.Duration `json:"update_interval"`
		MetricsRetention time.Duration `json:"metrics_retention"`

		// Thresholds for automatic backend disabling
		Thresholds struct {
			MaxResponseTime time.Duration `json:"max_response_time"`
			MaxFailureRate  float64       `json:"max_failure_rate"`
			MaxCPUUsage     float64       `json:"max_cpu_usage"`
			MaxMemoryUsage  float64       `json:"max_memory_usage"`
			MaxQueueLength  int           `json:"max_queue_length"`
		} `json:"thresholds"`

		// Alerting configuration
		Alerts struct {
			Enabled     bool     `json:"enabled"`
			Endpoints   []string `json:"endpoints"`
			MinSeverity string   `json:"min_severity"`
		} `json:"alerts"`

		// Metrics export
		Export struct {
			Prometheus bool   `json:"prometheus"`
			Port       int    `json:"port"`
			Path       string `json:"path"`
		} `json:"export"`
	} `json:"monitoring"`

	// Dynamic reconfiguration settings
	DynamicConfig struct {
		Enabled            bool          `json:"enabled"`
		UpdateInterval     time.Duration `json:"update_interval"`
		MinBackends        int           `json:"min_backends"`
		MaxBackends        int           `json:"max_backends"`
		ScaleUpThreshold   float64       `json:"scale_up_threshold"`
		ScaleDownThreshold float64       `json:"scale_down_threshold"`
	} `json:"dynamic_config"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// var config Config
	// if err := json.NewDecoder(file).Decode(&config); err != nil {
	// 	return nil, err
	// }

	// initialize load balancer according to the layer
	// var layerConfig interface{}
	// if config.Layer == LayerFour {
	// 	var l4_settings L4Settings
	// 	if err := json.NewDecoder(file).Decode(&config); err != nil {
	// 		return nil, err
	// 	}
	// 	layerConfig = l4_settings
	// } else if config.Layer == LayerSeven {
	// 	var l7_settings L7Settings
	// 	if err := json.NewDecoder(file).Decode(&config); err != nil {
	// 		return nil, err
	// 	}
	// 	layerConfig = l7_settings
	// } else {
	// 	return nil, fmt.Errorf("unsupported layer: %s", config.Layer)
	// }

	// if !config.Algorithm.IsValidForLayer(config.Layer) {
	// 	return nil, fmt.Errorf("algorithm %s is not valid for layer %s", config.Algorithm, config.Layer)
	// }

	// return &Config{
	// 	LayerConfig: layerConfig,
	// }, nil
	var temp struct {
		Layer Layer `json:"layer"`
	}

	if err := json.NewDecoder(file).Decode(&temp); err != nil {
		return nil, fmt.Errorf("failed to decode layer: %w", err)
	}

	file.Seek(0, 0)
	decoder := json.NewDecoder(file)

	var config Config
	switch temp.Layer {
	case LayerFour:
		var l4Config struct {
			Config
			L4Settings L4Settings `json:"l4_settings"`
		}
		if err := decoder.Decode(&l4Config); err != nil {
			return nil, err
		}
		config = l4Config.Config
		config.LayerConfig = l4Config.L4Settings

	case LayerSeven:
		var l7Config struct {
			Config
			L7Settings L7Settings `json:"l7_settings"`
		}
		if err := decoder.Decode(&l7Config); err != nil {
			return nil, err
		}
		config = l7Config.Config
		config.LayerConfig = l7Config.L7Settings
	default:
		return nil, fmt.Errorf("unsupported layer type: %s", temp.Layer)
	}

	if !config.Algorithm.IsValidForLayer(config.Layer) {
		return nil, fmt.Errorf("algorithm %s is not valid for layer %s", config.Algorithm, config.Layer)
	}

	return &config, nil
}

func (b *Backend) IsHealthy() bool {
	return b.Enabled &&
		!b.MaintenanceMode &&
		b.Metrics.HealthCheckStatus &&
		b.Metrics.ActiveConnections < int64(b.MaxConns) &&
		b.Metrics.CPUUsage < b.MaxCPUUsage &&
		b.Metrics.MemoryUsage < b.MaxMemoryUsage &&
		b.Metrics.ResponseTime < b.MaxResponseTime
}

func (b *Backend) UpdateMetrics(metrics metrics.ServerMetrics) {
	b.Metrics = metrics
	b.QueueFull = b.Metrics.QueuedConnections >= int64(b.MaxQueueSize)
}
