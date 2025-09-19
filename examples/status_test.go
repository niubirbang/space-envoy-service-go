package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestStatus(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	status, err := m.Status()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(status)
}
