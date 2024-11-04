import (
	"fmt"
	"log"

	"github.com/shubhamojha1/heimdall/config"
)

// refer to server metrics needed by the lb. important
func main() {
	cfg, err := config.LoadConfig("../../config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}