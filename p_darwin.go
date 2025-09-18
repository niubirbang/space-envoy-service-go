//go:build darwin

package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"path"
	"strings"
	"time"
)

func (m *Manager) initClient() {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", path.Join("/tmp", fmt.Sprintf("%s.sock", m.serviceName)))
		},
	}
	m.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func (m *Manager) install() error {
	quotedPath := fmt.Sprintf(`"%s"`, m.serviceFile)
	shell := strings.Join(
		[]string{
			fmt.Sprintf(`chmod +x %s`, quotedPath),
			fmt.Sprintf(`%s install`, quotedPath),
			fmt.Sprintf(`%s start`, quotedPath),
		},
		"\n",
	)
	script := fmt.Sprintf(
		`do shell script "%s" with prompt "Kernel %s requires authorization to use" with administrator privileges`,
		strings.ReplaceAll(shell, `"`, `\"`),
		m.serviceName,
	)
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install: %v\n%s", err, string(output))
	}
	var ok bool
	for i := 0; i < 60; i++ {
		time.Sleep(500 * time.Millisecond)
		if err := m.checkService(); err == nil {
			ok = true
			break
		}
	}
	if !ok {
		return errors.New("socket failed")
	}
	return nil
}

func (m *Manager) uninstall() error {
	quotedPath := fmt.Sprintf(`"%s"`, m.serviceName)
	shell := strings.Join(
		[]string{
			fmt.Sprintf(`chmod +x %s`, quotedPath),
			fmt.Sprintf(`%s uninstall`, quotedPath),
		},
		"\n",
	)
	script := fmt.Sprintf(
		`do shell script "%s" with prompt "Kernel %s requires authorization to use" with administrator privileges`,
		strings.ReplaceAll(shell, `"`, `\"`),
		m.serviceName,
	)
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install: %v\n%s", err, string(output))
	}
	return nil
}
