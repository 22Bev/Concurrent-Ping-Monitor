package main

import (
	"fmt"
	"time"

	"github.com/22Bev/concurrentpingmonitor/cmd/monitor"
)

func main() {
	hosts := []string{
		"google.com",
		"github.com",
		"openai.com",
		"invalid.example",
	}

	m := monitor.NewMonitor(monitor.HTTPPinger{})

	for {
		m.CheckHosts(hosts)
		m.PrintResults()
		fmt.Println()
		time.Sleep(10 * time.Second)
	}
}
