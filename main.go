package main

import (
	"AQM1248/LcdProc"
	"machine"
	"time"
)

func timer1ms(ch chan<- bool) {
	for {
		time.Sleep(10 * time.Millisecond)
		ch <- true
	}
}
func main() {
	var err error

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	spi := machine.SPI0
	err = spi.Configure(machine.SPIConfig{})
	if err != nil {
		panic(err)
	}

	dc := machine.GP14
	csn1 := machine.GP15
	csn2 := machine.GP17

	disp := LcdProc.New(spi, &csn1, &dc, &csn2)
	disp.Configure()

	ch := make(chan bool)
	go timer1ms(ch)

	for {
		select {
		case <-ch:
			break
		}
		led.Low()
		time.Sleep(100 * time.Millisecond)
		led.High()
		time.Sleep(100 * time.Millisecond)

		disp.LcdPrint(0, 0, "漢字だよ")
		disp.LcdPrint(0, 16, "本行　圭介")
		disp.LcdPrint(0, 32, "本行　ゆき")
	}

}
