//go:build !orangepi

package main

import (
	"fmt"
	"os"

	"github.com/erikdubbelboer/gspt"
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

func setProcTitle(status string) {
	if status != "" {
		status = " " + status
	}
	title := fmt.Sprintf("jetkvm: [supervisor]%s", status)
	gspt.SetProcTitle(title)
}
