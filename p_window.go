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

func (m *Manager) install() error {
	quotedPath := fmt.Sprintf(`"%s"`, m.serviceFile)
	for _, shell := range []string{
		fmt.Sprintf(`%s install`, quotedPath),
		fmt.Sprintf(`%s start`, quotedPath),
	} {
		script := fmt.Sprintf(
			`Start-Process "cmd.exe" -ArgumentList '/c %s' -Verb RunAs -WindowStyle Hidden`,
			shell,
		)
		fmt.Println(script)
		cmd := exec.Command("powershell", "-Command", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install: %v\n%s", err, string(output))
		}
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
		return errors.New("service failed to start")
	}
	return nil
}

func (m *Manager) uninstall() error {
	quotedPath := fmt.Sprintf(`"%s"`, m.serviceFile)
	for _, shell := range []string{
		fmt.Sprintf(`%s stop`, quotedPath),
		fmt.Sprintf(`%s uninstall`, quotedPath),
	} {
		script := fmt.Sprintf(
			`Start-Process "cmd.exe" -ArgumentList '/c %s' -Verb RunAs -WindowStyle Hidden`,
			shell,
		)
		fmt.Println(script)
		cmd := exec.Command("powershell", "-Command", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install: %v\n%s", err, string(output))
		}
	}
	return nil
}
