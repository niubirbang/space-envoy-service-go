package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestInit(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
}
