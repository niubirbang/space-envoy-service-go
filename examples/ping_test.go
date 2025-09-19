package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestPing(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	delay, err := m.Ping("115.236.100.149", 44542, 2000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(delay)
}
