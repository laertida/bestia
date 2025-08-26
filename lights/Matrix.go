package lights

import (
	"time"
)

type Matrix struct {
	left  Lights
	right Lights
}

func NewMatrix(left Lights, right Lights) Matrix {
	return Matrix{left: left, right: right}
}

func (matrix *Matrix) RowOn(number int) {
	matrix.left.On(number)
	matrix.right.On(number)
}

func (matrix *Matrix) RowOff(number int) {
	matrix.left.Off(number)
	matrix.right.Off(number)
}

func (matrix *Matrix) RowFlash(number int, during time.Duration) {
	matrix.RowOn(number)
	time.Sleep(during)
	matrix.RowOff(number)
	time.Sleep(during)
}

func (matrix *Matrix) RowToggle(number int) {
	matrix.left.Toggle(number)
	matrix.right.Toggle(number)
}

func (matrix *Matrix) AllOn() {
	matrix.left.AllOn()
	matrix.right.AllOn()
}

func (matrix *Matrix) AllOff() {
	matrix.left.AllOff()
	matrix.right.AllOff()
}

func (matrix *Matrix) Steps() {
	matrix.left.Step()
	matrix.right.Step()
}

func (matrix *Matrix) RowStep() {
	matrix.left.Step()
	matrix.right.Step()
}

func (matrix *Matrix) RowRandom() {
	matrix.left.RandomOn()
	matrix.right.RandomOn()
}
