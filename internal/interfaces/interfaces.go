package interfaces

import (
	"net"
)

// InterfaceStatus represents the status of a network interface
type InterfaceStatus struct {
	Name string
	Up   bool
	IPv4 []string
	IPv6 []string
}

// CheckInterfaces enumerates all network interfaces and their statuses
func CheckInterfaces() ([]InterfaceStatus, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var statuses []InterfaceStatus
	for _, iface := range ifaces {
		status := InterfaceStatus{
			Name: iface.Name,
			Up:   (iface.Flags & net.FlagUp) != 0,
		}

		addrs, err := iface.Addrs()
		if err != nil {
			statuses = append(statuses, status)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			if ip.To4() != nil {
				status.IPv4 = append(status.IPv4, ip.String())
			} else {
				status.IPv6 = append(status.IPv6, ip.String())
			}
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}
