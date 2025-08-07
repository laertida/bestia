package main

import (
	"fmt"
	"lmnl/bestia/lights"
	"lmnl/bestia/sensors"
	"os"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	// "periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"time"
)

const (
	TRIG_PORT = "23"
	ECHO_PORT = "24"
	LEFT_ADDR = 0x20
	RIGH_ADDR = 0x23
	IODIR     = 0x00
	GPIO      = 0x09
	ADDR      = 0x02
)

// L 0x40 0x20 0x10 0x08 0x04 0x02
// R 0x02 0x04 0x08 0x10 0x20 0x40

func main() {
	// Initialize periph.io
	L_MAP := map[int]int{1: 0x40, 2: 0x20, 3: 0x10, 4: 0x08, 5: 0x04, 6: 0x02}
	R_MAP := map[int]int{1: 0x02, 2: 0x04, 3: 0x08, 4: 0x10, 5: 0x20, 6: 0x40}

	if _, err := host.Init(); err != nil {
		fmt.Printf("Failed to initialize periph: %v \n", err)
		os.Exit(1)
		return
	}
	bus, err := i2creg.Open("")
	if err != nil {
		fmt.Printf("failed to open I2C bus: %v \n", err)
		os.Exit(1)
	}
	defer bus.Close()

	left := lights.NewLights("left side", LEFT_ADDR, L_MAP, bus)
	right := lights.NewLights("right side", RIGH_ADDR, R_MAP, bus)
	matrix := lights.NewLigthsMatrix(left, right)
	matrix.RowToggle(1)
	time.Sleep(1 * time.Second)
	matrix.RowToggle(1)
	trig := gpioreg.ByName(TRIG_PORT) // GPIO23
	echo := gpioreg.ByName(ECHO_PORT) // GPIO24

	//testMode := gpioreg.ByName("17")

	if trig == nil || echo == nil {
		fmt.Println("GPIO pins not found")
		os.Exit(1)
		return
	}

	trig.Out(gpio.Low)

	// Channel for distances
	distCh := make(chan float64)

	// Goroutine to read distances periodically
	go func() {
		for {
			dist := sensors.GetDistance(trig, echo)
			distCh <- dist
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Main goroutine: consume distance and do something
	for dist := range distCh {
		if dist >= 0 {
			fmt.Printf("distance is %.2f cm\n", dist)

		} else {
			fmt.Println("Timeout reading distance")
		}
	}
}
