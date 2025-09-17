package examples

import (
	"testing"

	ses "github.com/niubirbang/space-envoy-service-go"
)

func TestStop(t *testing.T) {
	ses.SetServiceFile("/Users/mac/go/src/github.com/niubirbang/space-envoy-service/space-envoy")
	if err := ses.Check(); err != nil {
		t.Fatal(err)
	}
	if err := ses.Stop(); err != nil {
		t.Fatal(err)
	}
}
