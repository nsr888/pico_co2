package display

import (
	"fmt"
	"log"

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
	font  *font.Display
	clear func()
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
	if f == nil || msg == "" {
		return
	}

	// Split longer messages into multiple lines
	lines := []string{msg}
	if len(msg) > 21 {
		lines = splitStringToLines(msg, 21)
	}

	f.clear()
	f.font.Configure(font.Config{FontType: font.FONT_6x8})

	for i, line := range lines {
		if i > 3 { // Max 4 lines on a 32px display
			break
		}
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

func (f *FontDisplay) DisplayNumReadings(r *types.Readings) {
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
	f.font.XPos = int16(128-(len(r.CO2String)*6)) / 2
	f.font.YPos = 24
	f.font.PrintText(r.CO2String)
}

func (f *FontDisplay) DisplayTextReadings(r *types.Readings) {
	if f == nil {
		return
	}
	f.clear()

	// Display CO2 status string at the top right corner
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	status := r.CO2String
	f.font.XPos = int16(128 - (len(status) * 7))
	f.font.YPos = 0
	f.font.PrintText(status)

	// Bars for CO2 level
	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText(printVerticalBar(r.CO2))

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	isValidStr := ""
	if !r.IsValid {
		isValidStr = "*"
	}
	formatString := fmt.Sprintf("CO2 %d%s %.0fC %.0f%%", r.CO2, isValidStr, r.Temperature, r.Humidity)
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(formatString)
}

func printVerticalBar(count uint16) string {
	var result string
	cnt := int(count)

	cnt = cnt - 400 // Adjust count to start from 400

	if cnt <= 0 {
		return result
	}

	numBars := cnt / 100

	for i := 0; i < numBars; i++ {
		result += "|"
	}

	return result
}
