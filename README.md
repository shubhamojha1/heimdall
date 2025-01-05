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

## Steps to start:
1. Start load balancer (also has service registry)
2. Start server manager
3. 

## Server Manager commands:-
1. Add new server: curl -X POST http://localhost:{MANAGER_PORT}/servers/add
2. List all servers running: curl http://localhost:{MANAGER_PORT}/servers/list
3. Get info about a specific server: curl http://localhost:{MANAGER_PORT}/servers/get?port={BASE_PORT}
4. Remove a server: curl -X DELETE http://localhost:{MANAGER_PORT}/servers/remove?port={BASE_PORT}

Note: 
- {MANAGER_PORT} is the port at which the server manager is running.
- {BASE_PORT} is the starting point from which new servers will be allocated ports (the closest one that i available)
- {LOAD_BALANCER_PORT}
