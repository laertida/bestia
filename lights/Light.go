package lights

import (
	"fmt"
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
}

func NewLights(name string, addr uint16, ports map[int]int, bus i2c.Bus) Lights {

	i2c := &i2c.Dev{Addr: addr, Bus: bus}
	// Set all GPIO as outputs
	_, err := i2c.Write([]byte{IODIR, 0x00})
	_, err = i2c.Write([]byte{GPIO, byte(INIT)})

	if err != nil {
		fmt.Println("failed to set IODIR: %v", err)
		os.Exit(1)
	}

	return Lights{addr: addr, i2c: i2c, lights: byte(0x00), name: name, ports: ports}
}

func (lights *Lights) On(number int) {
	port := lights.ports[number]
	lights.lights = lights.lights | byte(port)
	lights.i2c.Write([]byte{GPIO, byte(lights.lights)})
}

func (lights *Lights) Off(number int) {
	port := lights.ports[number]
	lights.lights = lights.lights &^ byte(port)
	lights.i2c.Write([]byte{GPIO, byte(lights.lights)})
}

func (lights *Lights) Toggle(number int) {
	port := lights.ports[number]
	lights.lights = lights.lights ^ byte(port)
	lights.i2c.Write([]byte{GPIO, byte(lights.lights)})

}

func (lights *Lights) Flash(number int, during time.Duration) {
	lights.On(number)
	time.Sleep(during)
	lights.Off(number)
	time.Sleep(during)
}

type LightsMatrix struct {
	left  Lights
	right Lights
}

func NewLigthsMatrix(left Lights, right Lights) LightsMatrix {
	return LightsMatrix{left: left, right: right}
}

func (matrix *LightsMatrix) RowOn(number int) {
	matrix.left.On(1)
	matrix.right.On(1)
}

func (matrix *LightsMatrix) RowOff(number int) {
	matrix.left.Off(1)
	matrix.right.Off(1)
}

func (matrix *LightsMatrix) RowFlash(number int, during time.Duration) {
	matrix.RowOn(1)
	time.Sleep(during)
	matrix.RowOff(1)
	time.Sleep(during)
}

func (matrix *LightsMatrix) RowToggle(number int) {
	matrix.left.Toggle(number)
	matrix.right.Toggle(number)
}

func (matrix *LightsMatrix) TestMode() {
	for i := range 6 {
		matrix.RowToggle(i)
	}
}
