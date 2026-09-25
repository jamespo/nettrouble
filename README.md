# nettrouble

`nettrouble` is a network troubleshooting CLI tool written in Go. It provides a comprehensive report on your network status, including DNS resolution, ping latency, traceroute, bandwidth performance, and security checks.

## Features

- **Network Interfaces**: Lists active interfaces and their IP addresses (IPv4 and IPv6).
- **DNS Check**: Verifies DNS resolution for common domains using OS settings.
- **Ping Latency**: Measures round-trip time for both IPv4 and IPv6 (if available).
- **Traceroute**: Performs a traceroute to identify network hops.
- **Bandwidth Test**: Conducts an internal speedtest to measure download and upload speeds.
- **MITM Detection**: Checks if HTTPS requests are being intercepted (Man-in-the-Middle).
- **Cross-Platform**: Supports macOS, Linux, and Windows.

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed.

```bash
git clone https://github.com/jamespo/nettrouble.git
cd nettrouble
go build -o nettrouble ./cmd/nettrouble/
./nettrouble
```

## Usage

```bash
./nettrouble                    # full diagnostics (requires root for ping)
sudo ./nettrouble               # privileged raw-socket ping
./nettrouble -skip-speedtest    # skip the bandwidth test (it's slow)
```

## Development

This project was created using **Gemini**, an interactive AI CLI agent.

## License

[MIT](LICENSE)
