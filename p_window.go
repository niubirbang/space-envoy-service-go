//go:build windows

package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/natefinch/npipe"
)

func (m *Manager) initClient() {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return npipe.DialTimeout(fmt.Sprintf(`\\.\pipe\%s`, m.serviceName), 30*time.Second)
		},
	}
	m.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func (m *Manager) install() error {
	return nil
}

func (m *Manager) uninstall() error {
	return nil
}
