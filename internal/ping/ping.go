package ping

import (
	"time"

	"github.com/go-ping/ping"
)

// PingResult contains statistics for a ping test
type PingResult struct {
	Target   string
	Sent     int
	Received int
	Loss     float64
	Min      time.Duration
	Avg      time.Duration
	Max      time.Duration
	Error    error
}

// RunPing performs a ping test to the specified target
func RunPing(target string, count int, privileged bool) (*PingResult, error) {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return &PingResult{Target: target, Error: err}, nil
	}

	pinger.Count = count
	pinger.Timeout = time.Second * 5
	pinger.SetPrivileged(privileged)

	err = pinger.Run()
	if err != nil {
		return &PingResult{Target: target, Error: err}, nil
	}

	stats := pinger.Statistics()
	return &PingResult{
		Target:   target,
		Sent:     stats.PacketsSent,
		Received: stats.PacketsRecv,
		Loss:     stats.PacketLoss,
		Min:      stats.MinRtt,
		Avg:      stats.AvgRtt,
		Max:      stats.MaxRtt,
	}, nil
}
