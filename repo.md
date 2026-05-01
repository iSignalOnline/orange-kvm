# JetKVM — Repository Overview

## What is this project?

**JetKVM** is an open-source, high-performance **KVM over IP** (Keyboard, Video, Mouse) solution that lets you remotely manage computers, servers, and workstations at the BIOS/hardware level — even when the target machine has no OS running. It runs on dedicated JetKVM hardware powered by a Rockchip RV1106 SoC (ARMv7).

Capabilities include:
- 1080p@60 FPS H.264 video capture with 30–60 ms latency
- Full keyboard and mouse emulation (USB HID gadget)
- Virtual mass-storage (ISO/image mounting)
- Optional cloud remote access via WebRTC (JetKVM Cloud)
- Optional Tailscale / Headscale VPN integration
- mDNS device discovery
- MQTT support with Home Assistant auto-discovery
- OTA (Over-the-Air) firmware updates
- Wake-on-LAN
- Prometheus metrics endpoint
- TLS support (self-signed or user-defined certificate)

---

## Repository Structure

```
.
├── Dockerfile.build          # Docker builder image for cross-compilation
├── Makefile                  # Top-level build, test, release targets
├── go.mod / go.sum           # Go module definition
├── main.go                   # Application entry point (calls kvm.Main())
├── cmd/main.go               # CLI wrapper
├── config.go                 # Configuration loading/saving (/userdata/kvm_config.json)
├── web.go / web_tls.go       # HTTP(S) web server (Gin)
├── webrtc.go                 # WebRTC peer connection (cloud streaming)
├── video.go                  # Video capture and H.264 encoding
├── usb.go / usb_mass_storage.go  # USB gadget and virtual storage
├── serial.go / cdc_acm_console.go  # Serial/CDC-ACM console
├── network.go                # Network interface management
├── tailscale.go              # Tailscale integration
├── mdns.go                   # mDNS advertisement
├── mqtt.go / mqtt_*.go       # MQTT client and Home Assistant discovery
├── ota.go                    # OTA update orchestration
├── jiggler.go                # Mouse jiggler (anti-idle)
├── timesync.go               # NTP time synchronization
├── prometheus.go             # Prometheus metrics
├── display.go                # On-device display (LVGL)
├── failsafe.go               # Failsafe / watchdog mode
├── wol.go                    # Wake-on-LAN
├── cloud.go                  # JetKVM Cloud WebSocket client
├── hidrpc.go                 # HID RPC (keyboard/mouse commands)
├── jsonrpc.go                # JSON-RPC over WebSocket
│
├── internal/
│   ├── native/               # CGo bridge to native C/LVGL library (libjknative.a)
│   ├── ota/                  # OTA update logic
│   ├── usbgadget/            # USB gadget configuration
│   ├── network/              # Network types and helpers
│   ├── tailscale/            # Tailscale state machine
│   ├── timesync/             # Time synchronization
│   ├── supervisor/           # Process supervisor
│   ├── logging/              # Structured logging (zerolog)
│   └── ...
│
├── ui/                       # React + TypeScript frontend (Vite)
│   ├── src/                  # UI source code
│   └── e2e/                  # Playwright end-to-end tests
│
├── scripts/
│   ├── build_utils.sh        # Docker build context helper (generates entrypoint.sh)
│   ├── build_cgo.sh          # CMake cross-compile for native C library
│   ├── ci_helper.sh          # CI pipeline helper
│   └── dev_deploy.sh         # One-shot deploy to a connected device
│
└── .devcontainer/
    └── install-deps.sh       # APT packages + jetkvm-native-buildkit toolchain installer
```

---

## How It Runs

### On the Device

The compiled binary (`jetkvm_app`) runs as a Linux daemon on the JetKVM hardware (ARMv7, uClibc). At startup (`kvm.Main()`), it:

1. Checks for a **failsafe reason** (watchdog tripped, bad config, etc.)
2. Loads configuration from `/userdata/kvm_config.json`
3. Initialises the **USB gadget** (keyboard, mouse, mass-storage)
4. Initialises the **native C library** (LVGL display, video capture via Rockchip MPP/RGA)
5. Starts **NTP time sync**
6. Starts **mDNS** advertising (`_jetkvm._tcp`)
7. Starts **Prometheus** metrics
8. Starts **MQTT** client (if configured)
9. Starts the **HTTP web server** (Gin, port 80 by default) and optionally an **HTTPS server**
10. Connects to **JetKVM Cloud** via WebSocket / WebRTC (if a cloud token is configured)
11. Starts a background goroutine that checks for **OTA updates** every hour (after an initial 15-minute delay)
12. Waits for `SIGINT` / `SIGTERM` to shut down gracefully

Configuration is a JSON file stored at `/userdata/kvm_config.json` on the device. Defaults are embedded in `config.go`.

### Frontend

The UI is a **React + TypeScript** SPA built with Vite. It has three build targets:

| Target | Command | Purpose |
|---|---|---|
| `device` | `npm run build:device` | Bundled into the Go binary (via `statigz`) and served by the on-device HTTP server |
| `development` | `npm run dev` | Local hot-reload development against the cloud API |
| `production` | `npm run build` | Cloud-hosted version (app.jetkvm.com) |

---

## When Things Run (Triggers)

| Event | What happens |
|---|---|
| Device boot | `jetkvm_app` starts automatically via the init system |
| Push to `dev` or `main` branch | CI workflow (`build.yml`) cross-compiles the app and runs Go unit tests |
| Pull request | Same CI workflow as above |
| `make build_dev` | Builds the ARM binary (locally with buildkit toolchain or in Docker) |
| `make frontend` | Builds the React UI and gzip-compresses assets into `static/` |
| `make test` | Runs all Go unit tests (`go test ./...`) |
| `make test_e2e DEVICE_IP=...` | Runs Playwright E2E tests against a physical JetKVM device |
| `make dev_release` | Publishes a dev build to the R2 CDN bucket |

---

## Build System

Cross-compilation targets **ARMv7 (arm-rockchip830-linux-uclibcgnueabihf)**.

### Local build (with toolchain installed)
```bash
# Install toolchain (once)
.devcontainer/install-deps.sh

# Build frontend, then binary
make frontend
make build_dev
```

### Docker build (no local toolchain needed)
```bash
make build_dev        # auto-detects missing toolchain, runs in Docker
```

The Docker image is built from `Dockerfile.build`. The CI workflow assembles the build context with `scripts/ci_helper.sh prepare`, which calls `prepare_docker_build_context()` in `scripts/build_utils.sh`.

---

## About `Dockerfile.build` — Missing Files Explained

`Dockerfile.build` references two files that **do not exist at the repository root**:

### 1. `install-deps.sh`
- **Actual location:** `.devcontainer/install-deps.sh`
- **How it gets there:** `scripts/build_utils.sh → prepare_docker_build_context()` copies it from `.devcontainer/install-deps.sh` into a **temporary Docker build context directory** before the image is built.
- **What it does:** Installs APT build dependencies and downloads the `jetkvm-native-buildkit` cross-compilation toolchain from GitHub releases into `/opt/jetkvm-native-buildkit/`.

### 2. `entrypoint.sh` ⚠️
- **Actual location:** **Generated at runtime** — it does not exist as a file in the repository.
- **How it gets there:** `scripts/build_utils.sh → prepare_docker_build_context()` writes it inline into the temporary build context:
  ```bash
  cat > "${DOCKER_BUILD_CONTEXT_DIR}/entrypoint.sh" << 'EOF'
  #!/bin/bash
  git config --global --add safe.directory /build
  exec $@
  EOF
  chmod +x "${DOCKER_BUILD_CONTEXT_DIR}/entrypoint.sh"
  ```
- **What it does:** Marks `/build` as a safe Git directory (required in newer Git versions when the directory is owned by a different UID inside the container), then `exec`s the passed command.

> **Summary:** Neither file is missing from the project — they are intentionally absent from the repo root and are provided to the Docker build context programmatically by `scripts/build_utils.sh`. If you try to `docker build -f Dockerfile.build .` directly from the repo root without running the prepare step first, the build will fail because both files won't be present.

### What IS actually missing

The `Dockerfile.build` will fail if run directly without first calling `scripts/build_utils.sh`'s `prepare_docker_build_context` (or `scripts/ci_helper.sh prepare`). There is no top-level `Makefile` target or README note that warns about this. Adding a dedicated `make build_image` target that calls `prepare_docker_build_context` before running `docker build` would make this more self-documenting.

---

## Key Dependencies

| Library | Purpose |
|---|---|
| `gin-gonic/gin` | HTTP web framework |
| `pion/webrtc` | WebRTC (cloud video streaming) |
| `rs/zerolog` | Structured logging |
| `prometheus/client_golang` | Metrics |
| `eclipse/paho.mqtt.golang` | MQTT |
| `vishvananda/netlink` | Network interface management |
| `insomniacslk/dhcp` | DHCP client |
| `Masterminds/semver` | Semantic version comparison (OTA) |
| `go-co-op/gocron` | Scheduler (jiggler, etc.) |
| `beevik/ntp` | NTP sync |
| `creack/pty` | PTY for serial console |
| `pojntfx/go-nbd` | NBD for virtual mass-storage |
| LVGL (C, via CGo) | On-device display rendering |
| Rockchip MPP/RGA (C) | Hardware video encode/decode |
