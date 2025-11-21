//go:build darwin

package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
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
	cmd := exec.Command("launchctl", "print", fmt.Sprintf("system/%s", m.serverName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	re := regexp.MustCompile(`pid = (\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return false
	}
	return matches[1] != "0"
}

func (m *Manager) install() error {
	fmt.Println("[space-envoy] installing")

	quotedPath := fmt.Sprintf(`"%s"`, m.serverFile)
	shell := strings.Join(
		[]string{
			fmt.Sprintf(`chmod +x %s`, quotedPath),
			fmt.Sprintf(`%s install`, quotedPath),
		},
		"\n",
	)
	script := fmt.Sprintf(
		`do shell script "%s" with prompt "Kernel %s requires authorization to use" with administrator privileges`,
		strings.ReplaceAll(shell, `"`, `\"`),
		m.serverName,
	)
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install: %v\n%s", err, string(output))
	}
	return m.installServerAfterCheck()
}

func (m *Manager) uninstall() error {
	fmt.Println("[space-envoy] uninstalling")

	quotedPath := fmt.Sprintf(`"%s"`, m.serverName)
	shell := strings.Join(
		[]string{
			fmt.Sprintf(`%s uninstall`, quotedPath),
		},
		"\n",
	)
	script := fmt.Sprintf(
		`do shell script "%s" with prompt "Kernel %s requires authorization to uninstall" with administrator privileges`,
		strings.ReplaceAll(shell, `"`, `\"`),
		m.serverName,
	)
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to uninstall: %v\n%s", err, string(output))
	}
	return nil
}

func (m *Manager) log() (string, error) {
	body, err := os.ReadFile(fmt.Sprintf("/var/log/%s.out.log", m.serverName))
	if err != nil {
		return "", err
	}
	return string(body), nil
}
