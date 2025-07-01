package display

import (
	"fmt"
	"log"
	"strings"

	font "github.com/Nondzu/ssd1306_font"
	"machine"
	"tinygo.org/x/drivers/ssd1306"

	"pico_co2/internal/types"
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
	font       *font.Display
	clear      func()
	waitsCount int
}

func (f *FontDisplay) DisplayBasic(r *types.Readings) {
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

func (f *FontDisplay) DisplayError(msg string) {
	if f == nil {
		return
	}
	f.clear()
	f.font.Configure(font.Config{FontType: font.FONT_6x8})

	var (
		maxLines       = 4
		maxCharPerLine = displayWidth / 6
	)

	if len(msg) < maxCharPerLine {
		f.font.XPos = 0
		f.font.YPos = 0
		f.font.PrintText(msg)

		return
	}

	msg = strings.ReplaceAll(msg, "\n", " ")

	log.Println("splitting error message into lines")

	lines := splitStringToLines(msg, maxCharPerLine)
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	log.Printf("Displaying error message with %d lines", len(lines))
	for i, line := range lines {
		log.Printf("Display error, line %d: %s", i, line)
		f.font.XPos = 0
		f.font.YPos = int16(i * 8)
		f.font.PrintText(line)
	}
}

// Split string into multiple lines if it exceeds the display width
func splitStringToLines(s string, maxCharPerLine int) []string {
	lines := make([]string, 0)
	for i := 0; i < len(s); i += maxCharPerLine {
		end := i + maxCharPerLine
		if end > len(s) {
			end = len(s)
		}
		lines = append(lines, s[i:end])
	}

	return lines
}

func (f *FontDisplay) DisplayReadings(r *types.Readings) {
	if f == nil {
		return
	}
	f.clear()

	// Big numbers for eCO2 and temperature
	f.font.Configure(font.Config{FontType: font.FONT_16x26})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.CO2))
	tempStr := fmt.Sprintf("%.0f", r.Temperature)
	f.font.XPos = int16(128 - (len(tempStr) * 16))
	f.font.YPos = 0
	f.font.PrintText(tempStr)

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
	f.waitsCount = 0
}
