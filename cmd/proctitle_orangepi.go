//go:build orangepi

package main

import (
	"fmt"
	"os"

	"github.com/jetkvm/kvm"
	"github.com/jetkvm/kvm/internal/native"
	"github.com/jetkvm/kvm/internal/supervisor"
)

func program() {
	subcomponentOverride := os.Getenv(supervisor.EnvSubcomponent)
	if subcomponentOverride != "" {
		subcomponent = subcomponentOverride
	}
	switch subcomponent {
	case "native":
		native.RunNativeProcess(os.Args[0])
	default:
		kvm.Main()
	}
}

// setProcTitle is a no-op on Orange Pi Zero targets. The gspt package requires
// CGO-based setproctitle which is not needed on this platform.
func setProcTitle(status string) {
	if status != "" {
		status = " " + status
	}
	_ = fmt.Sprintf("jetkvm: [supervisor]%s", status)
}
