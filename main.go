package main

import (
	"AQM1248/LcdProc"
	"fmt"
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
	var cnt uint8
	var str string

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
			cnt += 1
			break
		}
		led.Low()
		time.Sleep(100 * time.Millisecond)
		led.High()
		time.Sleep(100 * time.Millisecond)

		disp.LcdPrint(0, 0, "漢字だよ")
		disp.LcdPrint(0, 16, "Test Program")
		str = fmt.Sprintf("Cnt=%d", cnt)
		disp.LcdPrint(0, 32, str)
	}

}
