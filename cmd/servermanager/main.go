package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/shubhamojha1/heimdall/internal/config"
)

// make a struct called serverInfo with config.Backend + startedAt time, status
type ServerInfo struct {
	Port       int       `json:"port"`
	URL        string    `json:"url"`
	Status     string    `json:"status"`
	StartedAt  time.Time `json:"started_at"`
	ConfigInfo config.Config
}

type ServerManager struct {
	servers    map[int]*ServerInfo
	mutex      sync.RWMutex
	configFile config.Config // not sure if needed or not, will see
	basePort   int
}

func NewServerManager(basePort int) *ServerManager {
	return &ServerManager{
		servers:  make(map[int]*ServerInfo),
		basePort: basePort,
	}
}

func (sm *ServerManager) findFreePort() (int, error) {
	port := sm.basePort
	// find the next free port after this one

	for {
		sm.mutex.RLock()
		_, exists := sm.servers[port]
		sm.mutex.RUnlock()

		if !exists {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				listener.Close()
				return port, nil
			}
		}
		port++
		if port > sm.basePort+1000 {
			// fmt.Errorf(port)
			return 0, fmt.Errorf("no free ports found in range")
		}
	}
}

func (sm *ServerManager) AddServer() (*ServerInfo, error) {
	port, err := sm.findFreePort()
	if err != nil {
		return nil, err
	}

	ServerInfo := &ServerInfo{
		Port:      port,
		URL:       fmt.Sprintf("http://localhost:%d", port),
		Status:    "starting",
		StartedAt: time.Now(),
	}

	// start a server as a goroutine
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Server running on port %d\n", port)
		})

		sm.mutex.Lock()
		ServerInfo.Status = "running"
		sm.servers[port] = ServerInfo
		// add other info as well from config file
		sm.mutex.Unlock()

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			sm.mutex.Lock()
			ServerInfo.Status = "error"
			sm.mutex.Unlock()
		}
	}()

	// wait for server to start
	time.Sleep(100 * time.Millisecond)
	return ServerInfo, nil
}

func (sm *ServerManager) handleAddServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	server, err := sm.AddServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(server)
}

func main() {
	managerPort := os.Getenv("MANAGER_PORT")
	if managerPort == "" {
		managerPort = "7000"
	}

	basePort := os.Getenv("BASE_PORT")
	if basePort == "" {
		basePort = "8000"
	}

	basePortInt, _ := strconv.Atoi(basePort)
	manager := NewServerManager(basePortInt)

	mux := http.NewServeMux()

	mux.HandleFunc("/servers/add", manager.handleAddServer)
	// mux.HandleFunc("/servers/remove", manager.handleRemoveServer)
	// mux.HandleFunc("/servers/list", manager.handleListServers)
	// mux.HandleFunc("/servers/get", manager.handleGetServer)

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server Manager running...")
	})

	// Add basic API documentation
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, `Server Manager API:
		POST   /servers/add           - Add a new server
		DELETE /servers/remove?port=X - Remove server by port
		GET    /servers/list          - List all servers
		GET    /servers/get?port=X    - Get server info by port
		GET    /status               - Check manager status
		`)
	})

	log.Printf("Starting Server Manager on port %s", managerPort)
	log.Printf("Managing servers starting from port %s", basePort)
	log.Fatal(http.ListenAndServe(":"+managerPort, mux))

}
