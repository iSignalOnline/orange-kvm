//go:build linux && orangepi

// CGO-free stubs for the Orange Pi Zero target. The Rockchip-specific native
// library (libjknative / liblvgl) is not available on Allwinner hardware, so
// the in-process Native backend is replaced entirely by the EmptyNativeInterface
// (or the gRPC NativeProxy when an external native daemon is present). These
// stubs satisfy the internal function signatures called by native.go,
// video.go, and display.go so the package compiles without CGO.

package native

func setUpNativeHandlers() {}

func uiInit(_ uint16) {}

func uiTick() {}

func uiSetVar(_ string, _ string) {}

func uiGetVar(_ string) string { return "" }

func uiSwitchToScreen(_ string) {}

func uiGetCurrentScreen() string { return "" }

func uiObjAddState(_ string, _ string) (bool, error) { return false, nil }

func uiObjClearState(_ string, _ string) (bool, error) { return false, nil }

func uiObjAddFlag(_ string, _ string) (bool, error) { return false, nil }

func uiObjClearFlag(_ string, _ string) (bool, error) { return false, nil }

func uiObjHide(_ string) (bool, error) { return false, nil }

func uiObjShow(_ string) (bool, error) { return false, nil }

func uiObjSetOpacity(_ string, _ int) (bool, error) { return false, nil }

func uiObjFadeIn(_ string, _ uint32) (bool, error) { return false, nil }

func uiObjFadeOut(_ string, _ uint32) (bool, error) { return false, nil }

func uiLabelSetText(_ string, _ string) (bool, error) { return false, nil }

func uiImgSetSrc(_ string, _ string) (bool, error) { return false, nil }

func uiDispSetRotation(_ uint16) (bool, error) { return false, nil }

func uiEventCodeToName(_ int) string { return "" }

func uiGetLVGLVersion() string { return "" }

func videoInit(_ float64) error { return nil }

func videoShutdown() {}

func videoStart() {}

func videoStop() {}

func videoGetStreamingStatus() VideoStreamingStatus { return VideoStreamingStatusInactive }

func videoLogStatus() string { return "" }

func videoGetStreamQualityFactor() (float64, error) { return 0, nil }

func videoSetStreamQualityFactor(_ float64) error { return nil }

func videoSetCodecType(_ int) error { return nil }

func videoGetCodecType() (int, error) { return 0, nil }

func videoGetEDID() (string, error) { return "", nil }

func videoSetEDID(_ string) error { return nil }

func crash() {}
