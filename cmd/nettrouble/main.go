package main

import (
	"flag"
	"fmt"
	"os"

	"nettrouble/internal/bandwidth"
	"nettrouble/internal/dns"
	"nettrouble/internal/interfaces"
	"nettrouble/internal/mitm"
	"nettrouble/internal/ping"
	"nettrouble/internal/report"
	"nettrouble/internal/traceroute"
)

func main() {
	skipSpeed := flag.Bool("skip-speedtest", false, "skip the bandwidth speed test (it's slow)")
	privileged := flag.Bool("privileged", false, "use privileged (raw socket) ping — requires root")
	flag.Parse()

	dnsTargets := []string{"google.com", "cloudflare.com", "github.com"}
	pingTargets := []string{"8.8.8.8", "1.1.1.1"}
	tracerouteTarget := "8.8.8.8"
	mitmTarget := "google.com"

	fmt.Println("Checking network interfaces...")
	ifaces, err := interfaces.CheckInterfaces()
	if err != nil {
		fmt.Fprintf(os.Stderr, "interfaces error: %v\n", err)
	}

	fmt.Println("Checking DNS...")
	dnsResults := dns.CheckDNS(dnsTargets)

	fmt.Println("Running ping tests...")
	var pingResults []ping.PingResult
	for _, target := range pingTargets {
		r, err := ping.RunPing(target, 4, *privileged)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ping error for %s: %v\n", target, err)
			continue
		}
		pingResults = append(pingResults, *r)
	}

	fmt.Println("Running traceroute...")
	traceOutput, _ := traceroute.RunTraceroute(tracerouteTarget)

	fmt.Println("Checking for MITM/interception...")
	mitmResult, _ := mitm.CheckMITM(mitmTarget)

	var bwResult *bandwidth.SpeedTestResult
	if !*skipSpeed {
		fmt.Println("Running speed test (use -skip-speedtest to skip)...")
		bwResult, err = bandwidth.RunSpeedTest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "speed test error: %v\n", err)
		}
	}

	report.PrintReport(&report.Report{
		Interfaces: ifaces,
		DNS:        dnsResults,
		Ping:       pingResults,
		Traceroute: traceOutput,
		Bandwidth:  bwResult,
		MITM:       mitmResult,
	})
}
