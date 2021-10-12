package rpigpio

import (
	"fmt"
	"time"

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
	g.N1 = rpio.Pin(10) //GPIO23
	g.N1.Output()

	g.N2 = rpio.Pin(17) //GPIO17
	g.N2.Output()

	g.N3 = rpio.Pin(27) //GPIO27
	g.N3.Output()

	g.N4 = rpio.Pin(22) //GPIO22
	g.N4.Output()

	g.Clear()
	return &g
}

func (g *GpioBoard) MoveW() {
	fmt.Println("MoveW")

	g.Clear()

	g.N1.High()
	g.N2.Low()

	g.N3.High()
	g.N4.Low()
}

func (g *GpioBoard) MoveS() {
	fmt.Println("MoveS")

	g.Clear()

	g.N1.Low()
	g.N2.High()

	g.N3.Low()
	g.N4.High()
}

func (g *GpioBoard) MoveA() {
	fmt.Println("MoveA")

	// g.Clear()

	g.N1.Toggle()
	// g.N2.Low()

}

func (g *GpioBoard) MoveD() {
	fmt.Println("MoveD")

	g.Clear()

	g.N3.High()
	g.N4.Low()

}

func (g *GpioBoard) Clear() {
	g.N1.Low()
	g.N2.Low()
	g.N3.Low()
	g.N4.Low()
	fmt.Println("clear pin to LOW")
	time.Sleep(time.Second / 5)
}
