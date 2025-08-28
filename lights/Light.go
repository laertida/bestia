package lights

import (
	"log/slog"
	"math/rand"
	"os"
	"periph.io/x/conn/v3/i2c"
	"time"
)

const (
	IODIR = 0x00
	GPIO  = 0x09
	INIT  = 0x00
)

type Lights struct {
	addr   uint16
	i2c    *i2c.Dev
	lights byte
	name   string
	ports  map[int]int
	on     int
	logger slog.Logger
}

func NewLights(name string, addr uint16, ports map[int]int, bus i2c.Bus, logger slog.Logger) Lights {

	i2c := &i2c.Dev{Addr: addr, Bus: bus}
	// Set all GPIO as outputs
	_, err := i2c.Write([]byte{IODIR, 0x00})
	_, err = i2c.Write([]byte{GPIO, byte(INIT)})

	if err != nil {
		logger.Error("There was a problem trying to init I2C protocol")
		os.Exit(1)
	}

	NewLights := Lights{addr: addr, i2c: i2c, lights: byte(0x00), name: name, ports: ports, on: 1, logger: logger}
	NewLights.Write() // off all lights on instanciation
	return NewLights
}

func (lights *Lights) Write() {
	lights.logger.Debug("Lighths on", "name", lights.name, "lights", byte(lights.lights))
	lights.i2c.Write([]byte{GPIO, byte(lights.lights)})
}

func (lights *Lights) On(number int) {
	port := lights.ports[number]
	lights.lights = lights.lights | byte(port)
	lights.Write()
}

func (lights *Lights) RandomOn() {
	rand := rand.Intn(len(lights.ports) + 1)
	if rand == 0 {
		rand = 1
	}
	port := lights.ports[rand]
	lights.lights = byte(port)
	lights.Write()
}

func (lights *Lights) Off(number int) {
	port := lights.ports[number]
	lights.lights = lights.lights &^ byte(port)
	lights.Write()
}

func (lights *Lights) Step() {
	if lights.on == len(lights.ports) {
		lights.on = 1
	} else {
		lights.on += 1
	}
	lights.lights = byte(lights.ports[lights.on])
	lights.Write()
}

func (lights *Lights) Toggle(number int) {
	port := lights.ports[number]
	lights.lights = lights.lights ^ byte(port)
	lights.Write()
}

func (lights *Lights) AllOn() {
	lights.lights = 0xFF
	lights.Write()
}

func (lights *Lights) AllOff() {
	lights.lights = 0x00
	lights.Write()
}

func (lights *Lights) Flash(number int, during time.Duration) {
	lights.On(number)
	time.Sleep(during)
	lights.Off(number)
	time.Sleep(during)
}
