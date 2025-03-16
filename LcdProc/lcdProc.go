package LcdProc

import (
	"AQM1248/AQM1248"
	"AQM1248/GT20L16J1Y"
	"machine"
)

// Display wraps
type Display struct {
	lcd  *aqm1248.Device
	Font *GT20L16J1Y.Device
	XPos uint16
	YPos uint16
}

func New(spi *machine.SPI, lcdCsn *machine.Pin, dc *machine.Pin, fntCsn *machine.Pin) *Display {
	l := aqm1248.New(spi, lcdCsn, dc)
	f := GT20L16J1Y.New(spi, fntCsn)

	return &Display{
		lcd:  l,
		Font: f,
	}
}

func (d *Display) Configure() {
	d.lcd.Configure()
	d.Font.Initialize()
}

func (d *Display) SetPixel(x uint16, y uint16) {
	d.lcd.SetPixel(x, y, aqm1248.DotOn)
}

func (d *Display) UnSetPixel(x uint16, y uint16) {
	d.lcd.SetPixel(x, y, aqm1248.DotOff)
}

func (d *Display) Clear() {
	d.lcd.LcdClear()
}

func (d *Display) ShowPicture() {
	d.lcd.ShowPicture()
}

func (d *Display) LcdPrint(x uint16, y uint16, str string) {
	d.XPos = x // set position X
	d.YPos = y // set position Y
	d.printText(str)
}

func (d *Display) printText(str string) {
	var f GT20L16J1Y.Fonts
	tmp := d.XPos
	f = d.Font.ReadFonts(str)
	for i := 0; i < len(f); i++ {
		// Font Data Output
		d.printChar(f[i])
		d.XPos += uint16(f[i].FontWidth)
	}
	d.XPos = tmp
	d.ShowPicture()
}

func (d *Display) printChar(font GT20L16J1Y.Font) {
	var x, y uint16
	for y = 0; y < font.FontHeight; y++ {
		for x = 0; x < font.FontWidth; x++ {
			if font.FontData[x]&(0x01<<y) != 0x00 {
				d.SetPixel(x+d.XPos, y+d.YPos)
			} else {
				d.UnSetPixel(x+d.XPos, y+d.YPos)
			}
		}
	}
}
