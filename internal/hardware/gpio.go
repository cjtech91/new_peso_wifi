package hardware

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type GPIOPin struct {
	Number      int
	ActiveLevel ActiveLevel
}

func (p GPIOPin) path() string {
	return filepath.Join("/sys/class/gpio", fmt.Sprintf("gpio%d", p.Number))
}

func (p GPIOPin) export() error {
	if p.Number < 0 {
		return errors.New("invalid GPIO pin")
	}
	if _, err := os.Stat(p.path()); err == nil {
		return nil
	}
	return os.WriteFile("/sys/class/gpio/export", []byte(strconv.Itoa(p.Number)), 0o644)
}

func (p GPIOPin) setDirection(direction string) error {
	if err := p.export(); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(p.path(), "direction"), []byte(direction), 0o644)
}

func (p GPIOPin) Write(active bool) error {
	if p.Number < 0 {
		return nil
	}
	if err := p.setDirection("out"); err != nil {
		return err
	}
	level := "0"
	if active && p.ActiveLevel == ActiveHigh {
		level = "1"
	}
	if !active && p.ActiveLevel == ActiveLow {
		level = "1"
	}
	return os.WriteFile(filepath.Join(p.path(), "value"), []byte(level), 0o644)
}

