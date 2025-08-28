package sensors

import (
	"log/slog"
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
	logger  slog.Logger
}

func NewUltra(trigPort, echoPort string, low, up, tout int, logger slog.Logger) Ultra {
	trig := gpioreg.ByName(trigPort)
	echo := gpioreg.ByName(echoPort)
	logger.Debug("ULTRA: Intial values", "trigPort", trigPort, "echoPort", echoPort, "timeLow", low, "timeUp", up, "timeout", tout)
	if trig == nil || echo == nil {
		logger.Error("GPIO pins not found")
		os.Exit(1)
	}

	trig.Out(gpio.Low)

	timeLow := time.Duration(low) * time.Millisecond
	timeUp := time.Duration(up) * time.Microsecond
	timeout := time.Duration(tout) * time.Millisecond

	return Ultra{trig: trig, echo: echo, timeLow: timeLow, timeUp: timeUp, timeout: timeout, logger: logger}
}

func (ultra *Ultra) GetDistance() float64 {
	ultra.logger.Debug("ULTRA: pin low", "port", ultra.trig, "ms", ultra.timeLow.Milliseconds())
	ultra.trig.Out(gpio.Low)
	time.Sleep(ultra.timeLow) // Short delay
	ultra.logger.Debug("ULTRA: pin up", "port", ultra.trig, "μs", ultra.timeUp.Microseconds())
	ultra.trig.Out(gpio.High)
	time.Sleep(ultra.timeUp)
	ultra.trig.Out(gpio.Low)

	// Wait for echo to go HIGH
	timeout := time.After(ultra.timeout)
	var pulseStart time.Time

	for ultra.echo.Read() == gpio.Low {
		select {
		case <-timeout:
			ultra.logger.Debug("ULTRA: timeout reached", "timeout", timeout)
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
			ultra.logger.Debug("ULTRA: timeout reached", "timeout", timeout)
			return -1
		default:
			time.Sleep(ultra.timeUp)
		}
	}
	pulseEnd = time.Now()

	pulseDuration := pulseEnd.Sub(pulseStart).Seconds()
	ultra.logger.Debug("ULTRA: pulse end received", "after", pulseDuration)
	distance := pulseDuration * 17150 // cm
	ultra.logger.Debug("ULTRA: calculated distance", "cm", distance)
	return distance
}
