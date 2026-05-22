package traceroute

import (
	"os/exec"
	"runtime"
)

// RunTraceroute executes the system's traceroute command
func RunTraceroute(target string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tracert", "-d", target)
	} else {
		// Use -n to avoid DNS lookups for hops, making it faster
		cmd = exec.Command("traceroute", "-n", "-w", "1", "-q", "1", target)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
