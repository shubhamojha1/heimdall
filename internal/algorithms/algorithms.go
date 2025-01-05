package algorithms

// look up server list from service registry
// do round robin of requests
// use a global counter/pointer to keep track of which server to hit
// do modulus for going from last to first and so on
import (
	"net"

	"github.com/shubhamojha1/heimdall/internal/registry"
)

// const (
// 	Loadbalancer = load_balancer.LoadBalancer
// )

// func main() {

// 	// get list of backends from service registry
// 	backends := load_balancer.LoadBalancer.ServiceRegistry.Backends

// }

func handleConnectionRoundRobin(lb *registry.ServiceRegistry, clientConn net.Conn) {
	backends := registry.Backends
	// backends :=
}
