package main

import (
	"flag"
	"fmt"
	"lmnl/bestia/helpers"
	"lmnl/bestia/lights"
	"lmnl/bestia/sensors"
	"log/slog"
	"os"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"sync/atomic"
	"time"
)

const (
	TRIG_PORT = "23"
	ECHO_PORT = "24"
	TEST_PORT = "17"
	LEFT_ADDR = 0x23
	RIGH_ADDR = 0x27
	IODIR     = 0x00
	GPIO      = 0x09
	ADDR      = 0x02
)

func main() {
	var timeLow = flag.Int("timeLow", 60, "time in milliseconds where trigger is low.")
	var timeUp = flag.Int("timeUp", 10, "time in microseconds where trigger is low.")
	var timeOut = flag.Int("timeout", 1000, "time in milliseconds to consider timeout for the ultrasonic sensor.")
	var lightsOn = flag.Bool("lightsOn", false, "This flags allows to configure if lights will be on or off.")
	var logLevel = flag.String("logLevel", "error", "This flag determinates which is the log level, valid values are debug, info, warn and error ")
	flag.Parse()

	level := helpers.GetLogLevel(*logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	slog.Info("Bestia started with this initial values", "timeLow", *timeLow, "timeUp", *timeUp, "timeout", *timeOut, "lightsOn", *lightsOn)

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

	left := lights.NewLights("left side", LEFT_ADDR, L_MAP, bus, *logger)
	right := lights.NewLights("right side", RIGH_ADDR, R_MAP, bus, *logger)
	matrix := lights.NewMatrix(left, right)
	ultra := sensors.NewUltra(TRIG_PORT, ECHO_PORT, *timeLow, *timeUp, *timeOut, *logger)

	testMode := gpioreg.ByName("17")
	test := false
	if testMode.Read() == gpio.High {
		test = true
	}

	var latestDist atomic.Value
	// Goroutine to read distances periodically
	go func() {
		for {
			dist := ultra.GetDistance()
			latestDist.Store(dist)
			time.Sleep(1 * time.Second)
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
				logger.Error("Distance lecture is not ok or got timeout")
				time.Sleep(step)
			}
			var interval time.Duration
			interval = 1 * time.Second

			if dist > 400 || (dist < 99 && dist > 0) {
				if test {
					matrix.AllOn()
					logger.Info("-  on  - Lights due test mode", "dist", dist)
				} else {
					matrix.AllOff()
					logger.Info("- off - Lights off due: [dist > 400 || (dist < 99 && dist > 0)]", "dist", dist)
				}
			} else if dist < 400 && dist > 100 {

				if *lightsOn == true {
					matrix.RowStep()
					logger.Info("-normal- Ligths on normal pattern due: [dist < 400 && dist > 100]", "dist", dist)
				}
			} else {
				interval = time.Duration(100 * time.Millisecond)

				if *lightsOn == true {
					matrix.RowRandom()
					logger.Info("-random- Lights on due: dist in rage", "dist", dist)
				}
			}

			enlapsed := time.Duration(0)

			for enlapsed < interval {
				time.Sleep(step)
				enlapsed += step
				newDistVal := latestDist.Load()
				newDist, ok := newDistVal.(float64)

				if ok && ((dist > 100.0 && newDist <= 100.0) || (dist <= 100.0 && newDist > 100.0) || (dist-newDist) > 10.0 || (newDist-dist) > 10.0) {
					break
				}
			}
		}

	}()

	select {}
}
