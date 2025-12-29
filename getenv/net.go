package getenv

import (
	"fmt"
	"os"
	"strconv"
)

func NetworkPort() (int, error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return 0, fmt.Errorf("PORT environment variable is not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PORT value: %v", err)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port number out of range: %d", port)
	}

	return port, nil
}
