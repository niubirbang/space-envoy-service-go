//go:build windows

package ses

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/natefinch/npipe"
)

func buildClient() *http.Client {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return npipe.DialTimeout(fmt.Sprintf(`\\.\pipe\%s`, name), 30*time.Second)
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
