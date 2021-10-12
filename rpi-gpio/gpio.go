package rpigpio

import (
	"fmt"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

type GpioBoard struct {
	N1 rpio.Pin
	N2 rpio.Pin

	N3 rpio.Pin
	N4 rpio.Pin
}

func Setup() *GpioBoard {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rpio.Close()
	g := GpioBoard{}
	fmt.Printf("===>%T\n", rpio.Pin(4))
	g.N1 = rpio.Pin(4) //GPIO7
	g.N1.Output()

	g.N2 = rpio.Pin(17) //GPIO0
	g.N2.Output()

	g.N3 = rpio.Pin(27) //GPIO2
	g.N3.Output()

	g.N4 = rpio.Pin(22) //GPIO3
	g.N4.Output()

	g.Clear()
	return &g
}

func (g *GpioBoard) moveW() {
	g.N1.High()
	g.N2.Low()
}

func (g *GpioBoard) Clear() {
	g.N1.Low()
	g.N2.Low()
	g.N3.Low()
	g.N4.Low()
	fmt.Print("clear pin to LOW")
}
