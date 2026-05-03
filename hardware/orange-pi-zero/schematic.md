# Orange Pi Zero KVM — Hardware Schematic

This document describes the electrical connections for an Orange KVM device
built around an **Orange Pi Zero 3** (Allwinner H618). The same wiring applies
to the Orange Pi Zero 2 and the original Orange Pi Zero with minor connector
substitutions noted in parentheses.

---

## Block diagram

```
                                         ┌─────────────────────────────────┐
  Target machine                         │        Orange Pi Zero 3         │
  ┌──────────────┐                       │                                 │
  │              │ HDMI out              │  ┌─────────────┐                │
  │  Computer /  ├──────────────────────►│  │ USB HDMI    │                │
  │  Server /    │                       │  │ Capture     │ USB-A port     │
  │  Workstation │                       │  │ (MS2109)    │                │
  │              │                       │  └─────────────┘                │
  │              │ USB-A (HID)           │                       Ethernet  │
  │              │◄──────────────────────┤ USB-C (OTG/gadget) ◄──────────►│ LAN
  │              │  keyboard + mouse     │                                  │
  └──────────────┘                       │  microSD slot                   │
                                         │  ┌──────────┐                  │
  Power supply                           │  │  OS +    │                  │
  ┌──────────────┐                       │  │  App     │                  │
  │  5 V / 3 A   ├──────────────────────►│  └──────────┘                  │
  │  USB-C PD    │                       │                                 │
  └──────────────┘                       └─────────────────────────────────┘
                                                        │
                                              Remote user (browser / API)
```

---

## Detailed pin connections

### USB HDMI capture dongle → Orange Pi Zero 3

The capture dongle plugs into any of the Orange Pi Zero 3's USB-A host ports.
No additional wiring is required — the dongle is a standard UVC (USB Video
Class) device and is detected automatically by the Linux kernel.

| Signal | Dongle connector | OPi Zero 3 port |
|--------|-----------------|-----------------|
| HDMI input from target | HDMI female | — (dongle handles this) |
| USB data + power | USB-A male plug | USB-A host port 1 |

### Target machine USB-A → Orange Pi Zero 3 USB-C (HID gadget)

The orange-kvm app configures the USB-C OTG port as a **USB composite gadget**
(HID keyboard + HID mouse + mass storage). The target machine sees it as a
standard USB keyboard and mouse.

| Signal | OPi Zero 3 (USB-C OTG) | Target machine |
|--------|------------------------|----------------|
| VBUS (power sense) | USB-C pin A9 / B9 | USB-A pin 1 (+5 V) |
| D− | USB-C pin A7 / B7 | USB-A pin 2 |
| D+ | USB-C pin A6 / B6 | USB-A pin 3 |
| GND | USB-C pin A1/B1/A12/B12 | USB-A pin 4 |

> **For Orange Pi Zero (original):** the OTG port is micro-USB. Use a
> micro-USB to USB-A cable and connect the USB-A end to the target machine.

### Power supply

| Signal | Source | OPi Zero 3 connector |
|--------|--------|---------------------|
| +5 V | USB-C PD charger (3 A minimum) | USB-C power port |
| GND | Charger GND | USB-C GND |

> The USB-C OTG (gadget) port and the USB-C power port **are the same
> physical connector** on the Orange Pi Zero 3. The board determines
> whether the port acts as host or device via the `dr_mode` devicetree
> overlay. Use a USB-C splitter/hub that provides both power delivery and
> OTG data if you want to power the board and use OTG simultaneously, or
> power the board via the 40-pin header 5 V pins (pins 2 and 4) instead.

### Optional: ATX power-button control

To allow the orange-kvm app to remotely power-cycle the target machine,
wire the Orange Pi's GPIO to the target motherboard's front-panel header:

| Signal | OPi Zero 3 GPIO | Target motherboard header |
|--------|-----------------|--------------------------|
| PWR_BTN | GPIO PA6 (pin 7, 3.3 V logic) | PWR_BTN+ (via 10 kΩ pull-up) |
| GND | Pin 6 (GND) | PWR_BTN− |
| RESET | GPIO PA7 (pin 11, 3.3 V logic) | RST_BTN+ (via 10 kΩ pull-up) |
| GND | Pin 9 (GND) | RST_BTN− |

> Use a 3.3 V → open-drain level-shift (e.g. a BSS138 MOSFET) if the
> motherboard front-panel header is 5 V tolerant but not 3.3 V compatible.

---

## 40-pin header reference (Orange Pi Zero 3)

```
         3.3V  [ 1][ 2]  5V
    SDA1/PA12  [ 3][ 4]  5V
    SCL1/PA11  [ 5][ 6]  GND
         PA6   [ 7][ 8]  TXD/PH0
          GND  [ 9][10]  RXD/PH1
         PA7   [11][12]  PA8
         PA9   [13][14]  GND
        PA10   [15][16]  PA13
        3.3V   [17][18]  PA14
  SPI/MOSI/PH7 [19][20]  GND
  SPI/MISO/PH8 [21][22]  PA15
  SPI/CLK/PH9  [23][24]  SPI/CE0/PH10
          GND  [25][26]  SPI/CE1/PA16
         PA17  [27][28]  PA18
         PA19  [29][30]  GND
         PA20  [31][32]  PA21
         PA22  [33][34]  GND
         PA23  [35][36]  PA24
         PA25  [37][38]  PA26
          GND  [39][40]  PA27
```

Recommended GPIO for ATX control:
- **PWR_BTN** → pin 7 (PA6)
- **RESET** → pin 11 (PA7)
- **GND reference** → pin 9

---

## Network connectivity

The Orange Pi Zero 3 has a built-in 100 Mbps Ethernet port. Connect directly
to a switch or router for wired access. For wireless, attach a USB WiFi
dongle to one of the USB-A host ports (any MT7601U or RTL8188EUS based
dongle is supported by mainline Armbian kernels out of the box).

---

## Enclosure notes

The Orange Pi Zero 3 PCB is **65 × 30 mm**. A typical enclosure should
expose:
- 1× HDMI-A port (pass-through from the capture dongle)
- 1× USB-C port (power + OTG gadget cable to target)
- 1× RJ-45 Ethernet port
- 1× micro-SD slot (access for reflashing)
- Optional: USB-A port for WiFi dongle
