package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestUninstall(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Uninstall(); err != nil {
		t.Fatal(err)
	}
}
