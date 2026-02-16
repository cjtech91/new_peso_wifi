package main

import (
	"fmt"
	"log"

	"github.com/cjtech-nads/new_peso_wifi/internal/hardware"
)

func main() {
	board, err := hardware.DetectBoard()
	if err != nil {
		log.Fatalf("detect board: %v", err)
	}

	fmt.Printf("Detected board: %s (%s)\n", board.Name, board.ID)
	if !board.HasGPIO {
		fmt.Println("GPIO not available on this platform, running in simulation mode")
		return
	}

	relay := hardware.GPIOPin{
		Number:      board.Pins.RelayPin,
		ActiveLevel: board.Pins.RelayActive,
	}

	if err := relay.Write(true); err != nil {
		log.Fatalf("enable relay: %v", err)
	}
}

