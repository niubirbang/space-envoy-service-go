package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestLog(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	log, err := m.Log()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(log)
}
