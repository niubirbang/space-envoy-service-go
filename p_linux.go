//go:build linux

package service

import (
	"context"
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
			return net.Dial("unix", path.Join("/tmp", fmt.Sprintf("%s.sock", m.serverName)))
		},
	}
	m.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func (m *Manager) getServerIsRunningByServer() bool {
	cmd := exec.Command("systemctl", "is-active", m.serverName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.ToLower(string(output)) == "active"
}

func (m *Manager) install() error {
	fmt.Println("[space-envoy] installing")

	quotedPath := fmt.Sprintf(`"%s"`, m.serverFile)
	shells := []string{
		fmt.Sprintf(`chmod +x %s`, quotedPath),
		fmt.Sprintf(`%s install`, quotedPath),
	}
	for _, shell := range shells {
		cmd := exec.Command("pkexec", shell)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install: %v\n%s", err, string(output))
		}
	}
	return m.installServerAfterCheck()
}

func (m *Manager) uninstall() error {
	fmt.Println("[space-envoy] uninstalling")

	quotedPath := fmt.Sprintf(`"%s"`, m.serverFile)
	shells := []string{
		fmt.Sprintf(`%s uninstall`, quotedPath),
	}
	for _, shell := range shells {
		cmd := exec.Command("pkexec", shell)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to uninstall: %v\n%s", err, string(output))
		}
	}
	return nil
}

func (m *Manager) log() (string, error) {
	cmd := exec.Command("journalctl", "-u", m.serverName, "-n", "1000")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
