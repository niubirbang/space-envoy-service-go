//go:build linux

package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path"
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

func (m *Manager) isRunning() bool {
	return false
}

func (m *Manager) install() error {
	return nil
}

func (m *Manager) uninstall() error {
	return nil
}
