package hardware

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

type Edge string

const (
	EdgeRising  Edge = "rising"
	EdgeFalling Edge = "falling"
)

type ActiveLevel string

const (
	ActiveHigh ActiveLevel = "HIGH"
	ActiveLow  ActiveLevel = "LOW"
)

type BoardID string

const (
	BoardGenericX86 BoardID = "Generic x86_64"
)

type PinConfig struct {
	CoinPin       int
	RelayPin      int
	BillPin       int
	CoinPinEdge   Edge
	BillPinEdge   Edge
	RelayActive   ActiveLevel
}

type BoardConfig struct {
	ID       BoardID
	Name     string
	Pins     PinConfig
	HasGPIO  bool
	Variant  string
}

var compatibleToBoard = map[string]BoardConfig{
	"xunlong,orangepi-one": {
		ID:      "Orange Pi One",
		Name:    "Orange Pi One",
		Variant: "OP0100",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-pc": {
		ID:      "Orange Pi PC",
		Name:    "Orange Pi PC",
		Variant: "OP0600",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-pc-plus": {
		ID:      "Orange Pi PC Plus",
		Name:    "Orange Pi PC Plus",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-plus2e": {
		ID:      "Orange Pi Plus 2E",
		Name:    "Orange Pi Plus 2E",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-zero": {
		ID:      "Orange Pi Zero",
		Name:    "Orange Pi Zero",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-zero2": {
		ID:      "Orange Pi Zero 2",
		Name:    "Orange Pi Zero 2",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-zero3": {
		ID:      "Orange Pi Zero 3",
		Name:    "Orange Pi Zero 3",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     229,
			RelayPin:    228,
			BillPin:     73,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-3": {
		ID:      "Orange Pi 3",
		Name:    "Orange Pi 3",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-4": {
		ID:      "Orange Pi 4",
		Name:    "Orange Pi 4",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-5": {
		ID:      "Orange Pi 5",
		Name:    "Orange Pi 5",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-5b": {
		ID:      "Orange Pi 5B",
		Name:    "Orange Pi 5B",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-5-plus": {
		ID:      "Orange Pi 5 Plus",
		Name:    "Orange Pi 5 Plus",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"xunlong,orangepi-5-ultra": {
		ID:      "Orange Pi 5 Ultra",
		Name:    "Orange Pi 5 Ultra",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"friendlyarm,nanopi-neo": {
		ID:      "NanoPi NEO",
		Name:    "NanoPi NEO",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"friendlyarm,nanopi-neo2": {
		ID:      "NanoPi NEO2",
		Name:    "NanoPi NEO2",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"friendlyarm,nanopi-m1": {
		ID:      "NanoPi M1",
		Name:    "NanoPi M1",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     12,
			RelayPin:    11,
			BillPin:     6,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"raspberrypi,model-zero-w": {
		ID:      "Raspberry Pi Zero W",
		Name:    "Raspberry Pi Zero W",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     2,
			RelayPin:    3,
			BillPin:     4,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"raspberrypi,model-zero-2-w": {
		ID:      "Raspberry Pi Zero 2 W",
		Name:    "Raspberry Pi Zero 2 W",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     2,
			RelayPin:    3,
			BillPin:     4,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"raspberrypi,3-model-b": {
		ID:      "Raspberry Pi 3B",
		Name:    "Raspberry Pi 3B",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     2,
			RelayPin:    3,
			BillPin:     4,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"raspberrypi,3-model-b-plus": {
		ID:      "Raspberry Pi 3B+",
		Name:    "Raspberry Pi 3B+",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     2,
			RelayPin:    3,
			BillPin:     4,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"raspberrypi,4-model-b": {
		ID:      "Raspberry Pi 4B",
		Name:    "Raspberry Pi 4B",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     2,
			RelayPin:    3,
			BillPin:     4,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
	"raspberrypi,5-model-b": {
		ID:      "Raspberry Pi 5",
		Name:    "Raspberry Pi 5",
		HasGPIO: true,
		Pins: PinConfig{
			CoinPin:     2,
			RelayPin:    3,
			BillPin:     4,
			CoinPinEdge: EdgeRising,
			BillPinEdge: EdgeFalling,
			RelayActive: ActiveHigh,
		},
	},
}

func DetectBoard() (BoardConfig, error) {
	if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
		return BoardConfig{
			ID:      BoardGenericX86,
			Name:    "Generic x86_64",
			HasGPIO: false,
			Pins: PinConfig{
				CoinPin:     -1,
				RelayPin:    -1,
				BillPin:     -1,
				CoinPinEdge: EdgeRising,
				BillPinEdge: EdgeFalling,
				RelayActive: ActiveHigh,
			},
		}, nil
	}

	override := os.Getenv("PISO_BOARD_COMPATIBLE")
	if override != "" {
		if cfg, ok := compatibleToBoard[override]; ok {
			return cfg, nil
		}
	}

	data, err := os.ReadFile("/proc/device-tree/compatible")
	if err != nil {
		return BoardConfig{}, err
	}
	parts := strings.Split(string(data), "\x00")
	for _, p := range parts {
		if p == "" {
			continue
		}
		if cfg, ok := compatibleToBoard[p]; ok {
			return cfg, nil
		}
	}
	return BoardConfig{}, errors.New("unsupported board")
}
