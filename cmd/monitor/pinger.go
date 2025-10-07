package monitor

import (
	"net/http"
	"time"
)

// Pinger defines an interface for different ping strategies
type Pinger interface {
	Ping(host string) (time.Duration, error)
}

// HTTPPinger is a simple implementation using HTTP GET requests
type HTTPPinger struct{}

func (h HTTPPinger) Ping(host string) (time.Duration, error) {
	start := time.Now()
	resp, err := http.Get("https://" + host)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return time.Since(start), nil
}