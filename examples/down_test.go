package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestDown(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	if err := m.Down(); err != nil {
		t.Fatal(err)
	}
}
