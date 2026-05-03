//go:build linux && orangepi

// On Orange Pi Zero targets the Rockchip in-process native daemon is not
// available. RunNativeProcess is a no-op so that cmd/main.go can call it
// safely; it will never be invoked because the NativeProxy subprocess launch
// is disabled via JETKVM_NATIVE_DISABLE=1 in the systemd service.
// setProcTitle is a no-op for the same reason (gspt requires CGO).

package native

// RunNativeProcess is a no-op on Orange Pi Zero targets.
func RunNativeProcess(_ string) {}

func setProcTitle(_ string) {}

func updateProcessTitle(_ *VideoState) {}
