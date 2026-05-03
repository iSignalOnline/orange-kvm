//go:build orangepi

package kvm

import "fmt"

// setProcTitle is a no-op on Orange Pi Zero targets because the gspt package
// requires CGO (C-based setproctitle) which is not needed on this platform.
func setProcTitle(status string) {
	if status != "" {
		status = " " + status
	}
	_ = fmt.Sprintf("%s%s", procPrefix, status)
}
