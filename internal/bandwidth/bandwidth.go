package bandwidth

import (
	"fmt"
	"github.com/showwin/speedtest-go/speedtest"
)

// SpeedTestResult contains bandwidth test results
type SpeedTestResult struct {
	Download float64 // Mbps
	Upload   float64 // Mbps
	Latency  float64 // ms
	Error    error
}

// RunSpeedTest performs a speed test using speedtest.net servers
func RunSpeedTest() (*SpeedTestResult, error) {
	var speedtestClient = speedtest.New()

	serverList, err := speedtestClient.FetchServers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch servers: %w", err)
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil {
		return nil, fmt.Errorf("failed to find suitable server: %w", err)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no speedtest servers found")
	}

	server := targets[0]

	if err := server.PingTest(nil); err != nil {
		return nil, fmt.Errorf("ping test failed: %w", err)
	}

	if err := server.DownloadTest(); err != nil {
		return nil, fmt.Errorf("download test failed: %w", err)
	}

	if err := server.UploadTest(); err != nil {
		return nil, fmt.Errorf("upload test failed: %w", err)
	}

	return &SpeedTestResult{
		Download: float64(server.DLSpeed) / 1000000,
		Upload:   float64(server.ULSpeed) / 1000000,
		Latency:  float64(server.Latency.Milliseconds()),
	}, nil
}
