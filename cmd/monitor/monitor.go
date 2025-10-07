package monitor

import (
	"fmt"
	"sync"
)

type Result struct {
	Host     string
	Duration string
	Error    error
}

type Monitor struct {
	Pinger Pinger
	Mutex  sync.Mutex
	Results map[string]string
}

func NewMonitor(p Pinger) *Monitor {
	return &Monitor{
		Pinger: p,
		Results: make(map[string]string),
	}
}

func (m *Monitor) CheckHosts(hosts []string) []Result {
	resultsChan := make(chan Result)
	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			duration, err := m.Pinger.Ping(h)
			if err != nil {
				resultsChan <- Result{Host: h, Error: err}
				return
			}
			resultsChan <- Result{Host: h, Duration: duration.String()}
		}(host)
	}

	// Close channel once all goroutines complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []Result
	for r := range resultsChan {
		m.Mutex.Lock()
		if r.Error != nil {
			m.Results[r.Host] = "unreachable"
		} else {
			m.Results[r.Host] = r.Duration
		}
		m.Mutex.Unlock()
		results = append(results, r)
	}

	return results
}

func (m *Monitor) PrintResults() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	fmt.Println("----- Ping Results -----")
	for host, latency := range m.Results {
		fmt.Printf("%s: %s\n", host, latency)
	}
}
