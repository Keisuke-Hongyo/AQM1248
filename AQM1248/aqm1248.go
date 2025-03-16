package aqm1248

import (
	"AQM1248/GT20L16J1Y"
	"machine"
	"time"
)

const (
	LCD_WIDTH  = 128
	LCD_HEIGHT = 48
	PAGE       = 8

	DotOn  = 0x01
	DotOff = 0x00
)

type Device struct {
	spi  *machine.SPI
	csn  *machine.Pin
	rs   *machine.Pin
	font *GT20L16J1Y.Device
}

var gram [PAGE][LCD_WIDTH]byte

func New(spi *machine.SPI, csn *machine.Pin, rs *machine.Pin) *Device {
	return &Device{
		spi: spi,
		csn: csn,
		rs:  rs,
	}
}

func (aqm1248 *Device) Configure() {

	aqm1248.csn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	aqm1248.rs.Configure(machine.PinConfig{Mode: machine.PinOutput})
	aqm1248.rs.High()
	aqm1248.csn.High()

	// 初期化
	aqm1248.lcdCmd(0xe2) // RESET
	aqm1248.lcdCmd(0xae) // DC select: normal
	aqm1248.lcdCmd(0xa0) // Display OFF
	aqm1248.lcdCmd(0xc8) // Common output mode: reverse
	aqm1248.lcdCmd(0xa3) // LCD bias: 1/7

	// 内部レギュレータを順番にON
	aqm1248.lcdCmd(0x2c) // Power control 1 booster: ON
	time.Sleep(2 * time.Millisecond)
	aqm1248.lcdCmd(0x2e) // Power control 2 Voltage regulator circuit: ON
	time.Sleep(2 * time.Millisecond)
	aqm1248.lcdCmd(0x2f) // Power control 3 Voltage follower circuit: ON

	aqm1248.lcdCmd(0x23) // // V0 Voltage Regulator Internal Resistor Ratio Set
	aqm1248.lcdCmd(0x81) // The Electronic Volume Mode Set
	aqm1248.lcdCmd(0x1c) // Electronic Volume Register Set

	// 表示設定
	aqm1248.lcdCmd(0xa4) // Display all points: OFF
	aqm1248.lcdCmd(0x40) // Display start address: 0x00
	aqm1248.lcdCmd(0xa6) // Display: normal
	aqm1248.lcdCmd(0xaf) // Display: ON

	aqm1248.LcdClear()
}

func (aqm1248 *Device) LcdClear() {
	var x, p uint8

	for p = 0; p < PAGE; p++ {
		aqm1248.lcdCmd(0xb0 + p)
		aqm1248.lcdCmd(0x10)
		aqm1248.lcdCmd(0x00)
		for x = 0; x < LCD_WIDTH; x++ {
			aqm1248.lcdData(0x00)
		}
	}
}

func (aqm1248 *Device) SetPixel(x uint16, y uint16, dot uint8) {

	p := y / PAGE
	pDiv := y % PAGE

	if dot == DotOn {
		gram[p][x] |= DotOn << pDiv
	} else {
		gram[p][x] &= ^(DotOn << pDiv)
	}
}

func (aqm1248 *Device) ShowPicture() {
	var p, col uint8
	for p = 0; p < PAGE; p++ {
		aqm1248.lcdCmd(0xb0 + p)
		aqm1248.lcdCmd(0x10)
		aqm1248.lcdCmd(0x00)
		for col = 0; col < LCD_WIDTH; col++ {
			aqm1248.lcdData(gram[p][col])
		}
	}
}

// コマンド書込み
func (aqm1248 *Device) lcdCmd(cmd uint8) {
	aqm1248.csn.Low()
	aqm1248.rs.Low()
	_, _ = aqm1248.spi.Transfer(cmd)
	aqm1248.csn.High()
}

// データ書込み
func (aqm1248 *Device) lcdData(data uint8) {
	aqm1248.csn.Low()
	aqm1248.rs.High()
	_, _ = aqm1248.spi.Transfer(data)
	aqm1248.csn.High()
}
