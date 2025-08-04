package lights

import (
	"fmt"
	"periph.io/x/conn/v3/i2c"
)

const (
	IODIR = 0x00
	GPIO  = 0x09
)

type Light struct {
	gpio int
	on   bool
}

func NewLight(gpio int, on bool) Light {
	return Light{gpio: gpio, on: on}
}

type Stripe struct {
	name   string
	addr   uint16
	Lights [6]Light
	i2c    i2c.Dev
}

func NewStripe(name string, addr uint16, ports []int, bus i2c.Bus) Stripe {

	var Lights [6]Light
	for index, element := range ports {
		Lights[index] = Light{gpio: element, on: false}
	}

	i2c := i2c.Dev{Addr: addr, Bus: bus}
	// Set all GPIO as outputs
	_, err := i2c.Write([]byte{IODIR, 0x00})

	if err != nil {
		fmt.Println("failed to set IODIR: %v", err)
	}

	return Stripe{name: name, addr: addr, Lights: Lights, i2c: i2c}
}

func (stripe Stripe) LightOn(number int) {

	lights := byte(0x00)
	Light := stripe.Lights[number]
	lights |= 1 << Light.gpio
	stripe.i2c.Write([]byte{0x09, lights})
}
