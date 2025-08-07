package sensors

import (
	"periph.io/x/conn/v3/gpio"
	"time"
)

func GetDistance(trig, echo gpio.PinIO) float64 {
	trig.Out(gpio.Low)
	time.Sleep(60 * time.Millisecond) // Short delay

	trig.Out(gpio.High)
	time.Sleep(10 * time.Microsecond)
	trig.Out(gpio.Low)

	// Wait for echo to go HIGH
	timeout := time.After(1 * time.Second)
	var pulseStart time.Time
	for echo.Read() == gpio.Low {
		select {
		case <-timeout:
			return -1
		default:
			time.Sleep(10 * time.Microsecond)
		}
	}
	pulseStart = time.Now()

	// Wait for echo to go LOW
	timeout = time.After(1 * time.Second)
	var pulseEnd time.Time
	for echo.Read() == gpio.High {
		select {
		case <-timeout:
			return -1
		default:
			time.Sleep(10 * time.Microsecond)
		}
	}
	pulseEnd = time.Now()

	pulseDuration := pulseEnd.Sub(pulseStart).Seconds()
	distance := pulseDuration * 17150 // cm
	return distance
}
