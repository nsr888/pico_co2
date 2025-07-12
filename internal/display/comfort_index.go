package display

import (
	"fmt"
	"image/color"

	"tinygo.org/x/tinydraw"

	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayComfortIndex(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	font7 := font.NewFont7(f.display)

	status := r.ComfortStatus()
	if r.ValidityError != "" {
		status = r.ValidityError
	}
	font7.Print(0, 0, status)

	f.DrawHorizontalBar(0, 13, int16(r.AQI))

	font11 := font.NewFont11(f.display)
	if r.ValidityError == "" {
	}
	co2Str := fmt.Sprintf("%d", r.CO2)
	co2end := font11.Print(0, 16, co2Str)
	f.DrawVerticalBar(co2end+3, 16, r.CO2Index())

	var verticalBarWidth int16 = 9
	hum := fmt.Sprintf("%.0f", r.Humidity)
	humWidth := font11.CalcWidth(hum) + verticalBarWidth
	humPos := int16(128 - humWidth)
	font11.Print(humPos, 16, hum)
	f.DrawVerticalBar(128-6, 16, r.HumidityComfortIndex())

	spaceBetween := int16(6)

	temp := fmt.Sprintf("%.0f", r.Temperature)
	tempWidth := font11.CalcWidth(temp) + verticalBarWidth
	font11.Print(int16(humPos-spaceBetween-tempWidth), 16, temp)
	f.DrawVerticalBar(humPos-spaceBetween-6, 16, r.TempComfortIndex())
	f.display.Display()
}

func (f *FontDisplay) DrawVerticalBar(x, y, filledBars int16) {
	if f == nil {
		return
	}

	black := color.RGBA{1, 1, 1, 255}

	const maxBars = 4

	y += 2

	unfilledBars := maxBars - filledBars

	// Draw filled bars
	for i := int16(0); i < unfilledBars; i++ {
		f.display.SetPixel(x, y+i*4, black)
	}

	// Draw empty bars
	for i := unfilledBars; i < 4; i++ {
		tinydraw.Line(f.display, x, y+i*4, x+3, y+i*4, black)
	}
}

func (r *FontDisplay) DrawHorizontalBar(x, y, filledBars int16) {
	if r == nil {
		return
	}

	black := color.RGBA{1, 1, 1, 255}

	const maxBars = 5
	spaceBetween := int16(4)
	sectionWidth := int16(128-spaceBetween*(maxBars-1)) / maxBars

	// Draw filled bars
	for i := int16(0); i < filledBars; i++ {
		tinydraw.Line(r.display, x+i*sectionWidth+spaceBetween*i, y, x+i*sectionWidth+spaceBetween*i+sectionWidth, y, black)
	}
}
