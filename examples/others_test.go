package examples

import (
	"testing"

	ses "github.com/niubirbang/space-envoy-service-go"
)

func TestInstall(t *testing.T) {
	ses.SetServiceFile("/Users/mac/go/src/github.com/niubirbang/space-envoy-service/space-envoy")
	if err := ses.Install(); err != nil {
		t.Fatal(err)
	}
}

func TestUninstall(t *testing.T) {
	ses.SetServiceFile("/Users/mac/go/src/github.com/niubirbang/space-envoy-service/space-envoy")
	if err := ses.Uninstall(); err != nil {
		t.Fatal(err)
	}
}

func TestCheck(t *testing.T) {
	ses.SetServiceFile("/Users/mac/go/src/github.com/niubirbang/space-envoy-service/space-envoy")
	if err := ses.Check(); err != nil {
		t.Fatal(err)
	}
}
