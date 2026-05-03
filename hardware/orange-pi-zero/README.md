# Orange Pi Zero — KVM Hardware Guide

This document describes how to build an **Orange KVM** node using an
**Orange Pi Zero 2** or **Orange Pi Zero 3** as the main compute board.
The resulting device provides the same KVM-over-IP features as JetKVM
(keyboard, mouse, video capture, virtual storage, Wake-on-LAN, OTA) but
runs on low-cost, widely available Allwinner hardware.

---

## Supported boards

| Board | SoC | RAM | USB OTG | Recommended |
|---|---|---|---|---|
| Orange Pi Zero 2 | Allwinner H616 (ARM Cortex-A53, 64-bit) | 512 MB / 1 GB | USB-C | ✓ |
| Orange Pi Zero 3 | Allwinner H618 (ARM Cortex-A53, 1.5 GHz) | 1.5 / 2 / 4 GB | USB-C | ✓ (preferred) |
| Orange Pi Zero | Allwinner H2+ (ARM Cortex-A7, 32-bit) | 256 / 512 MB | micro-USB | Minimum spec |

> **Note**: The Go binary is compiled for `GOARCH=arm GOARM=7` (ARMv7) and runs on
> all three boards. The H616/H618 64-bit SoCs run the 32-bit binary without issues.

---

## Bill of Materials (BOM)

| Qty | Item | Purpose | Example part |
|-----|------|---------|--------------|
| 1 | Orange Pi Zero 3 (2 GB) | Main compute | Orange Pi Zero 3 |
| 1 | USB HDMI capture dongle | Capture video from target machine | MacroSilicon MS2109 (e.g. Elgato Cam Link compatible clone) |
| 1 | USB-C to USB-A cable (30 cm) | HID gadget connection to target machine | Any USB 2.0 cable |
| 1 | USB hub (2-port minimum) | Attach HDMI capture + spare port | Any USB 2.0 hub |
| 1 | USB-C power supply (5 V / 3 A) | Power the board | Any PD-capable charger |
| 1 | microSD card (8 GB+, Class 10) | OS + application | Samsung EVO or equivalent |
| 1 | Ethernet cable or USB WiFi dongle | Network connectivity | — |
| 1 | 3D-printed or off-the-shelf enclosure | Housing | See `hardware/orange-pi-zero/enclosure/` |

**Optional / recommended extras**

| Qty | Item | Purpose |
|-----|------|---------|
| 1 | ATX power control cable (2.54 mm header) | Remote power/reset buttons on target |
| 1 | USB serial adapter (3.3 V) | Debug console access to Orange Pi |
| 1 | 470 Ω resistor + 3 mm LED | Power-on status indicator |

---

## Schematic

See [`schematic.md`](schematic.md) for the full wiring diagram.

---

## OS setup

1. Flash **Armbian** (recommended) or the official Orange Pi OS image for your
   board to the microSD card.
2. Boot, configure networking, and install dependencies:

   ```bash
   sudo apt-get update && sudo apt-get install -y \
       libusb-1.0-0 v4l-utils
   ```

3. Enable USB gadget mode (OTG) in `/boot/armbianEnv.txt`:

   ```
   overlays=usbhost0 usbhost1 usbhost2 usbhost3
   # For USB gadget (HID) on the USB-C OTG port:
   overlays=usbhost0 usbhost1 usbhost2 usbhost3 otg_role=peripheral
   ```

4. Copy the `orange_kvm_app` binary to `/usr/local/bin/orange_kvm_app` and
   create a systemd service (see `hardware/orange-pi-zero/orange-kvm.service`).

---

## Build the binary

```bash
# Build locally (requires Go 1.25+ and CGO_ENABLED=0)
make _build_dev_orangepi_inner VERSION_DEV=0.5.8-opi-dev

# Or via Docker (no local Go needed)
make build_image_orangepi       # build the cross-compile Docker image (once)
make build_dev_orangepi         # produce bin/orange_kvm_app
```

The resulting `bin/orange_kvm_app` is a statically-linked ARMv7 binary.
The `orangepi` build tag swaps out the Rockchip CGO layer for no-op stubs;
the only remaining CGO dependency is `gspt` (process-title setter), which
links against standard libc and needs no Rockchip-specific libraries.

---

## Differences from JetKVM hardware

| Feature | JetKVM (RV1106) | Orange KVM (OPi Zero) |
|---|---|---|
| SoC | Rockchip RV1106 | Allwinner H2+/H616/H618 |
| Video capture | On-chip (Rockchip MPP/RGA hardware encoder) | USB HDMI capture dongle (V4L2/UVC) |
| LVGL display | On-device 240×240 LCD | Not included (headless) |
| Native daemon | `jetkvm_native` (CGO, libjknative) | Not used — `EmptyNativeInterface` |
| USB HID | USB gadget (built into SoC) | USB gadget via OTG port |
| Storage | eMMC | microSD |
| Power | 5 V / 2 A | 5 V / 3 A (USB-C) |
