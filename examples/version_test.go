package examples

import (
	"testing"

	ses "github.com/niubirbang/space-envoy-service-go"
)

func TestVersion(t *testing.T) {
	ses.SetServiceFile("/Users/mac/go/src/github.com/niubirbang/space-envoy-service/space-envoy")
	if err := ses.Check(); err != nil {
		t.Fatal(err)
	}
	version, err := ses.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(version)
}
