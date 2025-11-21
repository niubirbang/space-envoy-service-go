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
	"path"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
)

func (m *Manager) initClient() {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			f := fmt.Sprintf(`\\.\pipe\%s`, m.serverName)
			if _, err := os.Stat(f); err != nil {
				return nil, errors.New("pipe not found")
			}
			timeout := 5 * time.Second
			return winio.DialPipe(f, &timeout)
		},
	}
	m.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func (m *Manager) getServerIsRunningByServer() bool {
	cmd := exec.Command("sc", "query", m.serverName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "running")
}

func (m *Manager) install() error {
	fmt.Println("[space-envoy] installing")

	ps1 := path.Join(os.TempDir(), "space_service_install.ps1")
	if err := os.WriteFile(ps1, []byte(fmt.Sprintf(
		"Start-Process -FilePath \"%s\" -ArgumentList \"install\" -Verb RunAs -Wait -WindowStyle Hidden",
		m.serverFile,
	)), os.ModePerm); err != nil {
		return err
	}
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", fmt.Sprintf("\"%s\"", ps1))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install: %v\n%s", err, string(output))
	}
	return m.installServerAfterCheck()
}

func (m *Manager) uninstall() error {
	fmt.Println("[space-envoy] uninstalling")

	ps1 := path.Join(os.TempDir(), "space_service_uninstall.ps1")
	if err := os.WriteFile(ps1, []byte(fmt.Sprintf(
		"Start-Process -FilePath \"%s\" -ArgumentList \"uninstall\" -Verb RunAs -Wait -WindowStyle Hidden",
		m.serverFile,
	)), os.ModePerm); err != nil {
		return err
	}
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", fmt.Sprintf("\"%s\"", ps1))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to uninstall: %v\n%s", err, string(output))
	}
	return nil
}

func (m *Manager) log() (string, error) {
	script := fmt.Sprintf(`Get-EventLog -LogName Application -Source %s -Newest 1000`, m.serverName)
	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %v\n%s", err, string(output))
	}
	return string(output), nil
}
