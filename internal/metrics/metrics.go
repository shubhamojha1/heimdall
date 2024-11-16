package metrics

import (
	"time"
)

type ServerMetrics struct {
	// Connection metrics
	ActiveConnections  int64   `json:"active_connections"`
	TotalConnections   int64   `json:"total_connections"`
	ConnectionRate     float64 `json:"connection_rate"` // connections per second
	QueuedConnections  int64   `json:"queued_connections"`
	DroppedConnections int64   `json:"dropped_connections"`

	// Performance metrics
	ResponseTime     time.Duration `json:"response_time"`      // average response time
	LastResponseTime time.Duration `json:"last_response_time"` // last request response time
	ProcessingTime   time.Duration `json:"processing_time"`    // time spent processing request
	QueueTime        time.Duration `json:"queue_time"`         // time spent in queue

	// Resource metrics
	CPUUsage         float64 `json:"cpu_usage"`         // percentage
	MemoryUsage      float64 `json:"memory_usage"`      // percentage
	DiskIOUsage      float64 `json:"disk_io_usage"`     // IO operations per second
	NetworkBandwidth float64 `json:"network_bandwidth"` // bytes per second

	// Health metrics
	LastHealthCheck   time.Time `json:"last_health_check"`
	HealthCheckStatus bool      `json:"health_check_status"`
	FailureCount      int       `json:"failure_count"`
	SuccessCount      int       `json:"success_count"`

	// Layer 7 specific metrics
	RequestCount       int64 `json:"request_count,omitempty"`
	SuccessfulRequests int64 `json:"successful_requests,omitempty"`
	FailedRequests     int64 `json:"failed_requests,omitempty"`
	Status5xx          int64 `json:"status_5xx,omitempty"`
	Status4xx          int64 `json:"status_4xx,omitempty"`
	BytesSent          int64 `json:"bytes_sent,omitempty"`
	BytesReceived      int64 `json:"bytes_received,omitempty"`
}

// LoadBalancerMetrics represents global load balancer metrics
type LoadBalancerMetrics struct {
	TotalRequests       int64         `json:"total_requests"`
	ActiveConnections   int64         `json:"active_connections"`
	RequestsPerSecond   float64       `json:"requests_per_second"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	ErrorRate           float64       `json:"error_rate"`
	BackendsAvailable   int           `json:"backends_available"`
	BackendsTotal       int           `json:"backends_total"`
	LastUpdated         time.Time     `json:"last_updated"`
}

// MetricsStore interface for different metrics storage implementations
type MetricsStore interface {
	StoreMetrics(serverName string, metrics ServerMetrics) error
	GetMetrics(serverName string) (ServerMetrics, error)
	GetHistoricalMetrics(serverName string, duration time.Duration) ([]ServerMetrics, error)
}
