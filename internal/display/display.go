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
	displayWidth   int16 = 128
	displayHeight  int16 = 32
	displayAddress       = ssd1306.Address_128_32
)

func NewFontDisplay(bus *machine.I2C) (*FontDisplay, error) {
	display := ssd1306.NewI2C(bus)
	display.Configure(ssd1306.Config{
		Width:   displayWidth,
		Height:  displayHeight,
		Address: displayAddress,
	})
	log.Printf("Display configured: Width=%d, Height=%d, Address=%#x\n", displayWidth, displayHeight, displayAddress)

	fontLib := font.NewDisplay(display)
	return &FontDisplay{
		font:         &fontLib,
		clearDisplay: display.ClearDisplay,
	}, nil
}

type FontDisplay struct {
	font         *font.Display
	clearDisplay func()
}

func (f *FontDisplay) DisplayBasic(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	lines := strings.Split(r.ValidityError, ": ")

	f.printLines(lines[:1])
	f.printTempHumid(r.Humidity, r.Temperature)
}

func (f *FontDisplay) DisplayError(msg string) {
	// Split longer messages into multiple lines
	lines := []string{msg}
	maxCharPerLine := 128 / 7
	if len(msg) > maxCharPerLine {
		lines = splitStringToLines(msg, maxCharPerLine)
	}

	f.printLines(lines)
}

func (f *FontDisplay) printLines(lines []string) {
	if f == nil || len(lines) == 0 {
		return
	}

	f.clearDisplay()

	f.font.Configure(font.Config{FontType: font.FONT_7x10})

	for i, line := range lines {
		if i >= 3 { // Max 3 lines on a 32px display
			break
		}

		f.font.XPos = 0
		f.font.YPos = int16(i * 11)
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
	f.clearDisplay()

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
	f.font.XPos = int16(128-(len(r.CO2Status)*6)) / 2
	f.font.YPos = 24
	f.font.PrintText(r.CO2Status)
}

func (f *FontDisplay) DisplayTextReadings(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	// Display CO2 status string at the top right corner
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	status := r.CO2Status
	f.font.XPos = int16(128 - (len(status) * 7))
	f.font.YPos = 0
	f.font.PrintText(status)

	// Bars for CO2 level
	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText(printBar(r.CO2))

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	co2Str := fmt.Sprintf("CO2 %d", r.CO2)
	if r.ValidityError != "" {
		co2Str = fmt.Sprintf("CO2 %d*", r.CO2)
	}
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(co2Str)

	f.printTempHumid(r.Humidity, r.Temperature)
}

func (f *FontDisplay) printTempHumid(humidity, temperature float32) {
	if f == nil {
		return
	}

	humStr := fmt.Sprintf("H %.0f", humidity)
	f.font.XPos = int16(128 - (len(humStr) * 7))
	f.font.YPos = 24
	f.font.PrintText(humStr)
	tempStr := fmt.Sprintf("T %.0f", temperature)
	f.font.XPos = int16(128 - ((len(humStr) * 7) + (len(tempStr) * 7) + 8)) // 8 for padding
	f.font.YPos = 24
	f.font.PrintText(tempStr)
}

func printBar(count uint16) string {
	cnt := int(count)

	cnt = cnt - 400 // Adjust count to start from 400

	if cnt <= 0 {
		return ""
	}

	numBars := cnt / 100 // Corrected value 1 - 11

	return strings.Repeat("|", numBars)
}

func (f *FontDisplay) DisplayBigReadings(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	// Status top right corner
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	status := r.CO2Status
	f.font.YPos = 0
	f.font.XPos = int16(128 - (len(status) * 7))
	if r.ValidityError != "" {
		status = r.ValidityError
		f.font.XPos = 0
	}
	f.font.PrintText(status)

	// Bars for CO2 level
	if r.ValidityError == "" {
		f.font.Configure(font.Config{FontType: font.FONT_7x10})
		f.font.YPos = 0
		f.font.XPos = 0
		f.font.PrintText(printBar(r.CO2))
	}

	// Bottom big numbers
	if r.ValidityError == "" {
		f.font.Configure(font.Config{FontType: font.FONT_11x18})
		f.font.YPos = 16
		f.font.XPos = 0
		f.font.PrintText(fmt.Sprintf("%d", r.CO2))
	}

	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	temp := fmt.Sprintf("%.0f", r.Temperature)
	f.font.YPos = 16
	f.font.XPos = 64
	f.font.PrintText(temp)

	hum := fmt.Sprintf("%.0f", r.Humidity)
	f.font.YPos = 16
	f.font.XPos = 128 - 22 - 6
	f.font.PrintText(hum)
}

func (f *FontDisplay) DisplayReadingsWithHI(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	// CO2Status
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	status := fmt.Sprintf(
		"CO2 %s",
		strings.Repeat("*", r.CO2Rating()),
	)
	f.font.YPos = 0
	f.font.XPos = 0
	if r.ValidityError != "" {
		status = r.ValidityError
		f.font.XPos = 0
	}
	f.font.PrintText(status)

	if r.ValidityError == "" {
		// CO2 value
		f.font.Configure(font.Config{FontType: font.FONT_11x18})
		co2Value := fmt.Sprintf("%d", r.CO2)
		f.font.YPos = 0
		f.font.XPos = int16(128 - (len(co2Value) * 11))
		f.font.PrintText(co2Value)
	}

	// Heat Index status
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	hiStatus := fmt.Sprintf(
		"HI %s",
		strings.Repeat("*", r.HeatIndexRating()),
	)
	f.font.YPos = 24
	f.font.XPos = 0
	f.font.PrintText(hiStatus)

	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	temp := fmt.Sprintf("%.0f", r.Temperature)
	f.font.YPos = 16
	f.font.XPos = 70
	f.font.PrintText(temp)

	hum := fmt.Sprintf("%.0f", r.Humidity)
	f.font.YPos = 16
	f.font.XPos = 128 - 22
	f.font.PrintText(hum)
}
