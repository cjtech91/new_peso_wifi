# Deployment Guide – new_peso_wifi

This document describes how to deploy the `piso-wifi` binary on supported hardware platforms (Orange Pi, NanoPi, Raspberry Pi) and on a generic x86_64 machine for development.

> Note: At this stage the project focuses on hardware detection and GPIO control. PostgreSQL and Caddy integration will be added later; this guide will then be extended with database and HTTP/captive-portal deployment steps.

## Quick Start

1. Install dependencies:
   - On Debian/Ubuntu/Armbian/Raspberry Pi OS:
     ```bash
     sudo apt update
     sudo apt install -y git build-essential golang
     ```
2. Clone the repository:
   ```bash
   cd /opt
   sudo git clone https://github.com/cjtech91/new_peso_wifi.git
   sudo chown -R "$USER":"$USER" new_peso_wifi
   cd new_peso_wifi
   ```
3. Build the binary:
   ```bash
   go build ./cmd/piso-wifi
   sudo mv piso-wifi /usr/local/bin/
   ```
4. Run on SBC hardware (Orange Pi / NanoPi / Raspberry Pi):
   ```bash
   sudo piso-wifi
   ```
5. Run on a generic x86_64 PC (simulation mode, no GPIO):
   ```bash
   ./piso-wifi
   ```

The sections below provide more detailed information and configuration options.

## 1. Supported Platforms

The hardware detection logic uses the Linux device-tree `compatible` string to select the correct coin/bill/relay GPIO pins.

Supported boards include:

- Orange Pi (H2+/H3/H5 family and related):
  - Orange Pi One
  - Orange Pi One – OP0100
  - Orange Pi PC
  - Orange Pi PC – OP0600
  - Orange Pi PC Plus
  - Orange Pi Plus 2E
  - Orange Pi Zero (older H2+/H3 layout)
  - Orange Pi 3 / Orange Pi 3 – OP0300
  - Orange Pi 4
  - Orange Pi 5 / 5B / 5 Plus / 5 Ultra
- Orange Pi Zero 3 (H616/H618 family):
  - Orange Pi Zero 3
  - OrangePi Zero3
- Raspberry Pi family:
  - Raspberry Pi Zero W
  - Raspberry Pi Zero 2 W
  - Raspberry Pi 3B / 3B+
  - Raspberry Pi 4B
  - Raspberry Pi 5
- NanoPi boards:
  - NanoPi NEO
  - NanoPi NEO2
  - NanoPi M1
- Generic x86_64:
  - Treated as `Generic x86_64` with `gpio_disabled = true` (no real GPIO writes).

The relevant mapping is defined in `internal/hardware/board.go`.

## 2. Prerequisites

### 2.1. Operating System

- Linux distribution with:
  - `/proc/device-tree/compatible` available (typical on Armbian, Debian-based SBC images, Raspberry Pi OS).
  - Root access (via SSH or local console).

### 2.2. System Packages

On Debian/Ubuntu/Armbian/Raspberry Pi OS:

```bash
sudo apt update
sudo apt install -y git build-essential
```

### 2.3. Go Toolchain

Install Go either from the distribution or from the official tarball.

**Option A: distro packages (simpler, but version may be older)**

```bash
sudo apt install -y golang
go version
```

**Option B: official Go tarball (recommended if distro Go is outdated)**

```bash
cd /usr/local
sudo wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz   # adjust arch/version as needed
sudo tar -xzf go1.22.0.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
source /etc/profile
go version
```

## 3. Getting the Source Code

Choose a directory for the source, e.g. `/opt`:

```bash
cd /opt
sudo git clone https://github.com/cjtech91/new_peso_wifi.git
sudo chown -R "$USER":"$USER" new_peso_wifi
cd new_peso_wifi
```

If your actual repository URL or path is different, adjust accordingly.

## 4. Building the Binary

From the project root (where `go.mod` is located):

```bash
cd /opt/new_peso_wifi
go build ./cmd/piso-wifi
```

This produces an executable named `piso-wifi` in the project root.

For convenience, move it into a directory in your `PATH`, for example:

```bash
sudo mv piso-wifi /usr/local/bin/
```

Verify:

```bash
piso-wifi -h 2>/dev/null || echo "binary available"
```

## 5. Running on Target Hardware

### 5.1. Standard SBC (Orange Pi / NanoPi / Raspberry Pi)

Run the binary as root so it can access `/sys/class/gpio` and bind to the HTTP port:

```bash
sudo /usr/local/bin/piso-wifi
```

Expected behavior:

- The application detects the board based on `/proc/device-tree/compatible`.
- It starts an HTTP server (default on `:8080`, configurable with `PISO_HTTP_ADDR`).
- If the board has GPIO support, GPIO control is available to the process.

Client and admin portals:

- Client portal (voucher/coin UI):
  - URL: `http://DEVICE_IP:8080/`
  - Shows the board name and a simple form to submit a voucher code.
- Admin portal:
  - URL: `http://DEVICE_IP:8080/admin`
  - Basic page with board information and a link to JSON status.
- Admin JSON status:
  - URL: `http://DEVICE_IP:8080/admin/status`
  - Returns a JSON document with the detected board configuration.

### 5.2. Generic x86_64 (development / simulation)

On a regular PC or server with `GOARCH=amd64` or `GOARCH=386`, the board is treated as:

- `Generic x86_64`
- `HasGPIO = false`
- Coin/bill/relay pins set to `-1`

Run:

```bash
./piso-wifi
```

The HTTP server still runs (default `:8080`), but GPIO is disabled in this mode.

## 6. Forcing a Specific Board (Override)

You can override board detection for testing purposes by setting the `PISO_BOARD_COMPATIBLE` environment variable to one of the supported `compatible` strings.

Examples:

```bash
export PISO_BOARD_COMPATIBLE=raspberrypi,4-model-b
sudo /usr/local/bin/piso-wifi
```

Other valid values include (non-exhaustive):

- `xunlong,orangepi-one`
- `xunlong,orangepi-pc`
- `xunlong,orangepi-pc-plus`
- `xunlong,orangepi-plus2e`
- `xunlong,orangepi-zero`
- `xunlong,orangepi-zero2`
- `xunlong,orangepi-zero3`
- `xunlong,orangepi-3`
- `xunlong,orangepi-4`
- `xunlong,orangepi-5`
- `xunlong,orangepi-5b`
- `xunlong,orangepi-5-plus`
- `xunlong,orangepi-5-ultra`
- `raspberrypi,model-zero-w`
- `raspberrypi,model-zero-2-w`
- `raspberrypi,3-model-b`
- `raspberrypi,3-model-b-plus`
- `raspberrypi,4-model-b`
- `raspberrypi,5-model-b`
- `friendlyarm,nanopi-neo`
- `friendlyarm,nanopi-neo2`
- `friendlyarm,nanopi-m1`

If the override does not match any known key, detection falls back to the normal device-tree-based logic.

## 7. Running as a Systemd Service

To run `piso-wifi` automatically at boot, create a systemd unit.

1. Ensure the binary is at `/usr/local/bin/piso-wifi`.

2. Create the service file:

```bash
sudo tee /etc/systemd/system/piso-wifi.service >/dev/null <<'EOF'
[Unit]
Description=Piso WiFi Vendo Hardware Service
After=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/opt/new_peso_wifi
ExecStart=/usr/local/bin/piso-wifi
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
```

3. Reload systemd and enable the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable piso-wifi
sudo systemctl start piso-wifi
sudo systemctl status piso-wifi
```

If the service fails to start, check logs:

```bash
journalctl -u piso-wifi -e
```

## 8. Future: Database (PostgreSQL) and Caddy Integration

The current codebase focuses on hardware detection and GPIO control only. The planned architecture will also include:

- A Go HTTP server exposing:
  - Captive portal endpoints.
  - Status and control APIs.
- PostgreSQL integration for:
  - Storing credits, sessions, transactions, and configuration.
- Caddy configuration:
  - Acting as a reverse proxy in front of the Go HTTP server.
  - Handling TLS and automatic certificates when public-facing.

Once those components are implemented, this `deployment.md` will be extended with:

- PostgreSQL installation and schema migration steps.
- Configuration of environment variables (e.g. database DSN).
- Caddyfile examples for the captive portal and admin interface.

For now, you can deploy and run the `piso-wifi` binary using the steps above to validate GPIO behavior and hardware mapping on your target boards.

