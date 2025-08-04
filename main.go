package main

import (
	"fmt"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"time"
)

const (
	TRIG_PORT = "14"
	ECHO_PORT = "15"
	LEFT_ADDR = 0x20
	RIGH_ADDR = 0x21
	IODIR     = 0x00
	GPIO      = 0x09
)

func getDistance(trig, echo gpio.PinIO) float64 {
	// Ensure TRIG is LOW
	trig.Out(gpio.Low)
	time.Sleep(100 * time.Millisecond)

	// Send 10us pulse to TRIG
	trig.Out(gpio.High)
	time.Sleep(10 * time.Microsecond)
	trig.Out(gpio.Low)

	// Wait for ECHO to go HIGH
	start := time.Now()
	for echo.Read() == gpio.Low {
		if time.Since(start) > time.Second {
			return -1 // timeout
		}
	}
	pulseStart := time.Now()

	// Wait for ECHO to go LOW
	for echo.Read() == gpio.High {
		if time.Since(pulseStart) > time.Second {
			return -1 // timeout
		}
	}
	pulseEnd := time.Now()

	pulseDuration := pulseEnd.Sub(pulseStart).Seconds()
	distance := pulseDuration * 17150 // in cm
	return distance
}

func main() {
	// Initialize periph.io

	if _, err := host.Init(); err != nil {
		fmt.Println("Failed to initialize periph:", err)
		return
	}

	bus, err := i2creg.Open("")

	if err != nil {
		fmt.Println("failed to open I2C bus: %v", err)
	}

	dev := i2c.Dev{Addr: LEFT_ADDR, Bus: bus}
	dev2 := i2c.Dev{Addr: RIGH_ADDR, Bus: bus}
	//var _ conn.Conn = &dev

	_, err = dev.Write([]byte{IODIR, 0x00})
	_, err = dev2.Write([]byte{IODIR, 0x00})

	if err != nil {
		fmt.Println("failed to set IODIR: %v", err)
	}

	// Helper function to write to GPIO
	setLEDs := func(val byte) {
		_, err := dev.Write([]byte{GPIO, val})
		if err != nil {
			fmt.Println("failed to write GPIO: %v", err)
		}
		_, err = dev2.Write([]byte{GPIO, val})
		if err != nil {
			fmt.Println("failed to write GPIO: %v", err)
		}

	}

	fmt.Println("Turning all LEDs ON")
	setLEDs(0xFF) // "11111111" "10000000" -> hex
	time.Sleep(2 * time.Second)

	fmt.Println("Turning all LEDs OFF")
	setLEDs(0x00)
	time.Sleep(2 * time.Second)

	fmt.Println("Turning all LEDs OFF")
	setLEDs(0x80)
	time.Sleep(2 * time.Second)

	fmt.Println("Alternating LEDs")
	setLEDs(0xAA)
	time.Sleep(2 * time.Second)

	setLEDs(0x55)
	time.Sleep(2 * time.Second)
	fmt.Println("Turning all LEDs OFF")
	setLEDs(0x00)
	time.Sleep(2 * time.Second)

	fmt.Println("Done")

	trig := gpioreg.ByName(TRIG_PORT) // GPIO14
	echo := gpioreg.ByName(ECHO_PORT) // GPIO15

	if trig == nil || echo == nil {
		fmt.Println("GPIO pins not found")
		return
	}

	trig.Out(gpio.Low)

	// Channel for distances
	distCh := make(chan float64)

	// Goroutine to read distances periodically
	go func() {
		for {
			dist := getDistance(trig, echo)
			distCh <- dist
			time.Sleep(1 * time.Second)
		}
	}()

	// Main goroutine: consume distance and do something
	for dist := range distCh {
		if dist >= 0 {
			fmt.Printf("Distance: %.2f cm\n", dist)
		} else {
			fmt.Println("Timeout reading distance")
		}
		// You can add more logic here, e.g. trigger alerts, update UI, etc.
	}
}
