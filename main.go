package main

import (
	"flag"
	"fmt"
	"lmnl/bestia/lights"
	"lmnl/bestia/sensors"
	"os"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"sync/atomic"
	"time"
)

const (
	TRIG_PORT = "23"
	ECHO_PORT = "24"
	LEFT_ADDR = 0x22
	RIGH_ADDR = 0x27
	IODIR     = 0x00
	GPIO      = 0x09
	ADDR      = 0x02
)

func main() {

	var timeLow = flag.Int("timeLow", 60, "time in milliseconds where trigger is low.")
	var timeUp = flag.Int("timeUp", 10, "time in milliseconds where trigger is low.")
	var timeOut = flag.Int("timeout", 1000, "time in milliseconds to consider timeout for the ultrasonic sensor.")
	var lightsOn = flag.Bool("ligths", false, "This flags allows to configure if lights will be on or off.")
	var log = flag.Bool("log", true, "This flag allows to enable, disable log.")

	// L 0x40 0x20 0x10 0x08 0x04 0x02
	// R 0x02 0x04 0x08 0x10 0x20 0x40
	L_MAP := map[int]int{1: 0x40, 2: 0x20, 3: 0x10, 4: 0x08, 5: 0x04, 6: 0x02}
	R_MAP := map[int]int{1: 0x02, 2: 0x04, 3: 0x08, 4: 0x10, 5: 0x20, 6: 0x40}

	// Initialize periph.io
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
	matrix := lights.NewMatrix(left, right)
	ultra := sensors.NewUltra(TRIG_PORT, ECHO_PORT, *timeLow, *timeUp, *timeOut)

	var latestDist atomic.Value
	// Goroutine to read distances periodically
	go func() {
		for {
			dist := ultra.GetDistance()
			latestDist.Store(dist)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// Main goroutine: consume distance and do something
	go func() {

		step := 5 * time.Millisecond
		for {
			distVal := latestDist.Load()
			dist, ok := distVal.(float64)

			if !ok || dist < 0 {
				matrix.AllOff()
				time.Sleep(step)
			}
			var interval time.Duration
			interval = 1 * time.Second
			if *log {
				fmt.Printf("dist: %.2f\n", dist)
			}
			if dist > 400 {
				matrix.AllOff()
			} else if dist < 400 && dist > 100 {
				if *lightsOn {
					matrix.RowStep()
				}
			} else {
				interval = time.Duration(100 * time.Millisecond)
				if *lightsOn {
					matrix.RowRandom()
				}
			}

			enlapsed := time.Duration(0)

			for enlapsed < interval {
				time.Sleep(step)
				enlapsed += step
				newDistVal := latestDist.Load()
				newDist, ok := newDistVal.(float64)

				if ok && (dist > 100 && newDist < 100) {
					break
				}
			}
		}

	}()

	select {}
}
