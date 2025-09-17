package examples

import (
	"testing"

	ses "github.com/niubirbang/space-envoy-service-go"
)

func TestStart(t *testing.T) {
	ses.SetServiceFile("/Users/mac/go/src/github.com/niubirbang/space-envoy-service/space-envoy")
	if err := ses.Check(); err != nil {
		t.Fatal(err)
	}
	if err := ses.Start(
		"/Users/mac/go/src/github.com/MetaCubeX/mihomo/tmp",
		"/Users/mac/go/src/github.com/MetaCubeX/mihomo/tmp/config_feimiao.yaml",
	); err != nil {
		t.Fatal(err)
	}
}
