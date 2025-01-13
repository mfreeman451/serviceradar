package agent

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type ProcessChecker struct {
	ProcessName string
}

func (p *ProcessChecker) Check(ctx context.Context) (bool, string) {
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", p.ProcessName)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Sprintf("Process %s is not running", p.ProcessName)
	}
	status := strings.TrimSpace(string(output))
	return status == "active", fmt.Sprintf("Process status: %s", status)
}

type PortChecker struct {
	Host string
	Port int
}

func (p *PortChecker) Check(ctx context.Context) (bool, string) {
	addr := fmt.Sprintf("%s:%d", p.Host, p.Port)
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, fmt.Sprintf("Port %d is not accessible: %v", p.Port, err)
	}
	conn.Close()
	return true, fmt.Sprintf("Port %d is accessible", p.Port)
}
