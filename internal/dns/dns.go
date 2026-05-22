package dns

import (
	"context"
	"net"
	"time"
)

// DNSResult stores the resolution results for a domain
type DNSResult struct {
	Domain string
	IPv4   []string
	IPv6   []string
	Error  error
}

// CheckDNS performs lookups for the given domains
func CheckDNS(domains []string) []DNSResult {
	var results []DNSResult
	resolver := &net.Resolver{}

	for _, domain := range domains {
		res := DNSResult{Domain: domain}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		ips, err := resolver.LookupIP(ctx, "ip", domain)
		cancel()
		
		if err != nil {
			res.Error = err
		} else {
			for _, ip := range ips {
				if ip.To4() != nil {
					res.IPv4 = append(res.IPv4, ip.String())
				} else {
					res.IPv6 = append(res.IPv6, ip.String())
				}
			}
		}
		results = append(results, res)
	}
	return results
}
