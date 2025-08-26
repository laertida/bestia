package sensors

import (
	"fmt"
	"os"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"time"
)

type Ultra struct {
	trig    gpio.PinIO
	echo    gpio.PinIO
	timeLow time.Duration
	timeUp  time.Duration
	timeout time.Duration
}

func NewUltra(trigPort, echoPort string, low, up, tout int) Ultra {
	trig := gpioreg.ByName(trigPort)
	echo := gpioreg.ByName(echoPort)
	if trig == nil || echo == nil {
		fmt.Println("GPIO pins not found")
		os.Exit(1)
	}

	trig.Out(gpio.Low)

	timeLow := time.Duration(low) * time.Millisecond
	timeUp := time.Duration(up) * time.Millisecond
	timeout := time.Duration(tout) * time.Second

	return Ultra{trig: trig, echo: echo, timeLow: timeLow, timeUp: timeUp, timeout: timeout}
}

func (ultra *Ultra) GetDistance() float64 {
	ultra.trig.Out(gpio.Low)
	time.Sleep(ultra.timeLow) // Short delay

	ultra.trig.Out(gpio.High)
	time.Sleep(ultra.timeUp)
	ultra.trig.Out(gpio.Low)

	// Wait for echo to go HIGH
	timeout := time.After(ultra.timeout)
	var pulseStart time.Time
	for ultra.echo.Read() == gpio.Low {
		select {
		case <-timeout:
			return -1
		default:
			time.Sleep(ultra.timeUp)
		}
	}
	pulseStart = time.Now()

	// Wait for echo to go LOW
	timeout = time.After(ultra.timeout)
	var pulseEnd time.Time
	for ultra.echo.Read() == gpio.High {
		select {
		case <-timeout:
			return -1
		default:
			time.Sleep(ultra.timeUp)
		}
	}
	pulseEnd = time.Now()

	pulseDuration := pulseEnd.Sub(pulseStart).Seconds()
	distance := pulseDuration * 17150 // cm
	return distance
}
