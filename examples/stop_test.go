package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestStop(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	if err := m.Stop(); err != nil {
		t.Fatal(err)
	}
}
