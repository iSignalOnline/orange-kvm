//go:build orangepi

package kvm

// setProcTitle is a no-op on Orange Pi Zero targets because the gspt package
// requires CGO (C-based setproctitle) which is not needed on this platform.
func setProcTitle(_ string) {}
