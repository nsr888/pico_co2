package main

import (
	"fmt"
	"log"

	font "github.com/Nondzu/ssd1306_font"
	"machine"
	"tinygo.org/x/drivers/ssd1306"
)

// Display Configuration
const (
	displayWidth   = 128
	displayHeight  = 32
	displayAddress = ssd1306.Address_128_32
)

func NewFontDisplay(bus *machine.I2C) (*FontDisplay, error) {
	display := ssd1306.NewI2C(bus)
	display.Configure(ssd1306.Config{
		Width:   displayWidth,
		Height:  displayHeight,
		Address: displayAddress,
	})
	log.Printf("Display configured: Width=%d, Height=%d, Address=%#x", displayWidth, displayHeight, displayAddress)

	clear := func() {
		display.ClearBuffer()
		display.Display()
	}
	clear()

	fontLib := font.NewDisplay(display)
	return &FontDisplay{
		font:  &fontLib,
		clear: clear,
	}, nil
}

type FontDisplay struct {
	font  *font.Display
	clear func()
}

func (f *FontDisplay) DisplayAHT20Readings(r Readings) {
	if f == nil {
		return
	}
	f.clear()
	f.font.Configure(font.Config{FontType: font.FONT_16x26})
	tempStr := fmt.Sprintf("%.0f", r.Temperature)
	f.font.XPos = int16((128 - (len(tempStr) * 16)) / 2)
	f.font.YPos = 0
	f.font.PrintText(tempStr)

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_6x8})
	formatString := fmt.Sprintf("Temp %.1fC Hum %.1f%%", r.Temperature, r.Humidity)
	f.font.XPos = int16((128 - (len(formatString) * 6)) / 2)
	f.font.YPos = 24
	f.font.PrintText(formatString)
}

func (f *FontDisplay) DisplayFullReadings(r Readings) {
	if f == nil {
		return
	}
	f.clear()

	// Big numbers for eCO2 and AQI
	f.font.Configure(font.Config{FontType: font.FONT_16x26})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.ECO2))
	tempStr := fmt.Sprintf("%.0f", r.Temperature)
	f.font.XPos = int16(128 - (len(tempStr) * 16))
	f.font.YPos = 0
	f.font.PrintText(tempStr)

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_6x8})
	co2Str := "eCO2"
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(co2Str)
	tempTitleStr := "Temp"
	f.font.XPos = int16(128 - (len(tempTitleStr) * 6))
	f.font.YPos = 24
	f.font.PrintText(tempTitleStr)
	f.font.XPos = int16(128-(len(r.Status)*6)) / 2
	f.font.YPos = 24
	f.font.PrintText(r.Status)
}

func (f *FontDisplay) DisplayFullReadingsCO2andAQI(r Readings) {
	if f == nil {
		return
	}
	f.clear()

	// Big numbers for eCO2 and AQI
	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	f.font.XPos = 30
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.ECO2))
	f.font.XPos = 110
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.AQI))

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_6x8})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText("eCO2")
	f.font.XPos = 87
	f.font.YPos = 0
	f.font.PrintText("AQI")
	f.font.XPos = 0
	f.font.YPos = 16
	f.font.PrintText("-----------------------")
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(fmt.Sprintf("T %.0f H %.0f", r.Temperature, r.Humidity))
	f.font.XPos = int16(128 - (len(r.Status) * 6))
	f.font.YPos = 24
	f.font.PrintText(r.Status)
}
