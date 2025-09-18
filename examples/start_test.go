package examples

import (
	"testing"

	service "github.com/niubirbang/space-envoy-service-go"
)

func TestStart(t *testing.T) {
	m := service.NewManager(serviceName, serviceFile)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	if err := m.Start(
		"/Users/mac/go/src/github.com/MetaCubeX/mihomo/tmp",
		"/Users/mac/go/src/github.com/MetaCubeX/mihomo/tmp/config_feimiao.yaml",
	); err != nil {
		t.Fatal(err)
	}
}
