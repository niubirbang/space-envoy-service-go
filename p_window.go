//go:build windows

package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/natefinch/npipe"
)

func (m *Manager) initClient() {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			f := fmt.Sprintf(`\\.\pipe\%s`, m.serviceName)
			if _, err := os.Stat(f); err != nil {
				return nil, errors.New("pipe not found")
			}
			return npipe.Dial(f)
		},
	}
	m.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func (m *Manager) isRunning() bool {
	cmd := exec.Command("sc", "query", m.serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "running")
}

func (m *Manager) install() error {
	fmt.Println("installing")
	quotedPath := fmt.Sprintf(`"%s"`, m.serviceFile)
	for _, shell := range []string{
		fmt.Sprintf(`%s install`, quotedPath),
		// fmt.Sprintf(`%s start`, quotedPath),
	} {
		script := fmt.Sprintf(
			`Start-Process "cmd.exe" -ArgumentList '/c %s' -Verb RunAs -WindowStyle Hidden`,
			shell,
		)
		cmd := exec.Command("powershell", "-Command", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install: %v\n%s", err, string(output))
		}
	}
	return m.installAfterCheck()
}

func (m *Manager) uninstall() error {
	fmt.Println("uninstalling")
	quotedPath := fmt.Sprintf(`"%s"`, m.serviceFile)
	for _, shell := range []string{
		// fmt.Sprintf(`%s stop`, quotedPath),
		fmt.Sprintf(`%s uninstall`, quotedPath),
	} {
		script := fmt.Sprintf(
			`Start-Process "cmd.exe" -ArgumentList '/c %s' -Verb RunAs -WindowStyle Hidden`,
			shell,
		)
		cmd := exec.Command("powershell", "-Command", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install: %v\n%s", err, string(output))
		}
	}
	return nil
}

func (m *Manager) log() (string, error) {
	script := fmt.Sprintf(`Get-EventLog -LogName Application -Source %s -Newest 1000`, m.serviceName)
	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %v\n%s", err, string(output))
	}
	return string(output), nil
}
