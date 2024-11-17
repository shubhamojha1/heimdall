# heimdall
## A Dynamic Load Balancer

```
├── cmd/
│   └── load_balancer/
│       └── main.go                # Entry point of the application
│
├── internal/
│   ├── algorithms/                 # Load balancing algorithms implementation
│   │   ├── round_robin.go
│   │   ├── sticky_round_robin.go
│   │   ├── weighted_round_robin.go
│   │   ├── ip_url_hash.go
│   │   ├── least_connections.go
│   │   └── least_time.go
│   │
│   ├── server/                     # Server management and logic
│   │   ├── server_manager.go
│   │   └── health_check.go         # Health check logic
│   │
│   └── config/                     # Configuration management
│       └── config.go               # Loading and parsing config.json
│
├── pkg/
│   ├── logger/                     # Custom logging utilities
│   │   └── logger.go
│   │
│   ├── alerts/                     # Alert sending functionality
│   │   └── alerts.go
│   │
│   └── stress_test/                # Stress testing utilities
│       └── stress_test.go
│
├── test/                           # Integration and unit tests
│   └── load_balancer_test.go
│
├── config.json                     # Configuration file for load balancer
├── go.mod                          # Go module file
└── README.md                       # Project documentation
```