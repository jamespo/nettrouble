package mitm

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

// MITMResult contains the findings of an HTTPS interception check
type MITMResult struct {
	Target      string
	Intercepted bool
	Issuer      string
	IsTrusted   bool
	Error       error
}

// CheckMITM connects to a host over TLS and inspects the certificate chain
func CheckMITM(host string) (*MITMResult, error) {
	address := host + ":443"
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", address, &tls.Config{
		InsecureSkipVerify: false,
	})

	if err != nil {
		// If verification fails, try to dial without verification to see who the issuer is
		connInsecure, errInsecure := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", address, &tls.Config{
			InsecureSkipVerify: true,
		})
		if errInsecure != nil {
			return &MITMResult{Target: host, Error: errInsecure}, nil
		}
		defer connInsecure.Close()
		
		certs := connInsecure.ConnectionState().PeerCertificates
		if len(certs) > 0 {
			return &MITMResult{
				Target:      host,
				Intercepted: true,
				Issuer:      certs[0].Issuer.CommonName,
				IsTrusted:   false,
				Error:       err, // Original verification error
			}, nil
		}
		return &MITMResult{Target: host, Error: err}, nil
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return &MITMResult{Target: host, Error: fmt.Errorf("no certificates found")}, nil
	}

	issuer := certs[0].Issuer.CommonName
	
	// A simple heuristic for corporate MITM: issuers that don't look like public CAs.
	// This is not exhaustive but provides good feedback.
	suspiciousIssuers := []string{"Zscaler", "Fortinet", "Cisco", "BlueCoat", "Kaspersky", "Sophos", "ESET", "Bitdefender"}
	intercepted := false
	for _, s := range suspiciousIssuers {
		if containsIgnoreCase(issuer, s) {
			intercepted = true
			break
		}
	}

	return &MITMResult{
		Target:      host,
		Intercepted: intercepted,
		Issuer:      issuer,
		IsTrusted:   true,
	}, nil
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
