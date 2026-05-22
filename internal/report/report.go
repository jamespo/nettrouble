package report

import (
	"fmt"
	"strings"

	"nettrouble/internal/bandwidth"
	"nettrouble/internal/dns"
	"nettrouble/internal/interfaces"
	"nettrouble/internal/mitm"
	"nettrouble/internal/ping"
)

// Report aggregates all diagnostic results
type Report struct {
	Interfaces []interfaces.InterfaceStatus
	DNS        []dns.DNSResult
	Ping       []ping.PingResult
	Traceroute string
	Bandwidth  *bandwidth.SpeedTestResult
	MITM       *mitm.MITMResult
}

// PrintReport displays the diagnostics in a readable format
func PrintReport(r *Report) {
	fmt.Println("\n==========================================")
	fmt.Println("       NETWORK TROUBLESHOOTING REPORT      ")
	fmt.Println("==========================================\n")

	fmt.Println("--- Network Interfaces ---")
	for _, iface := range r.Interfaces {
		status := "DOWN"
		if iface.Up {
			status = "UP"
		}
		fmt.Printf("[%s] %s\n", status, iface.Name)
		if len(iface.IPv4) > 0 {
			fmt.Printf("  IPv4: %s\n", strings.Join(iface.IPv4, ", "))
		}
		if len(iface.IPv6) > 0 {
			fmt.Printf("  IPv6: %s\n", strings.Join(iface.IPv6, ", "))
		}
	}
	fmt.Println()

	fmt.Println("--- DNS Resolution ---")
	for _, d := range r.DNS {
		if d.Error != nil {
			fmt.Printf("%s: ERROR (%v)\n", d.Domain, d.Error)
		} else {
			fmt.Printf("%s: OK (v4: %d, v6: %d)\n", d.Domain, len(d.IPv4), len(d.IPv6))
		}
	}
	fmt.Println()

	fmt.Println("--- Ping Latency ---")
	for _, p := range r.Ping {
		if p.Error != nil {
			fmt.Printf("%s: ERROR (%v)\n", p.Target, p.Error)
		} else {
			fmt.Printf("%s: %d/%d received, Avg RTT: %v\n", p.Target, p.Received, p.Sent, p.Avg)
		}
	}
	fmt.Println()

	fmt.Println("--- Traceroute ---")
	if r.Traceroute != "" {
		fmt.Println(r.Traceroute)
	} else {
		fmt.Println("No traceroute data available.")
	}
	fmt.Println()

	fmt.Println("--- Bandwidth / Speedtest ---")
	if r.Bandwidth != nil {
		fmt.Printf("Download: %.2f Mbps\n", r.Bandwidth.Download)
		fmt.Printf("Upload:   %.2f Mbps\n", r.Bandwidth.Upload)
		fmt.Printf("Latency:  %.2f ms\n", r.Bandwidth.Latency)
	} else {
		fmt.Println("Bandwidth test failed or was skipped.")
	}
	fmt.Println()

	fmt.Println("--- MITM / Security Check ---")
	if r.MITM != nil {
		if r.MITM.Error != nil && r.MITM.Issuer == "" {
			fmt.Printf("Security check failed: %v\n", r.MITM.Error)
		} else {
			fmt.Printf("Target: %s\n", r.MITM.Target)
			fmt.Printf("Issuer: %s\n", r.MITM.Issuer)
			status := "CLEAN"
			if r.MITM.Intercepted {
				status = "SUSPICIOUS (Possible Interception)"
			}
			if !r.MITM.IsTrusted {
				status = "DANGER (Untrusted Certificate)"
			}
			fmt.Printf("Status: %s\n", status)
		}
	}
	fmt.Println("\n==========================================")
}
