//go:build linux

package ses

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path"
	"time"
)

func buildClient() *http.Client {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", path.Join("/tmp", fmt.Sprintf("%s.sock", name)))
		},
	}
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func Install() error {

	return nil
}

func Uninstall() error {
	return nil
}
