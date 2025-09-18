package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestVersion(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	version, err := m.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(version)
}
