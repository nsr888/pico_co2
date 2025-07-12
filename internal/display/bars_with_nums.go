package display

import (
	"fmt"
	"image/color"

	"pico_co2/internal/display/bar"
	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayBarsWithLargeNums(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	var (
		y              int16
		x              int16
		radiusFilled   int16 = 3
		barSpacing     int16 = 5
		leftBarsCount  int16 = 0
		rightBarsCount int16 = 4
	)

	hbar := bar.NewTwoSideBar(
		f.display,
		radiusFilled,
		barSpacing,
		leftBarsCount,
		rightBarsCount,
		font.NewProggy(f.display, color.RGBA{1, 1, 1, 255}),
	)

	font16 := font.NewFreemonoRegular12(f.display, color.RGBA{1, 1, 1, 255})

	x = 0
	y = 0
	// Print first line with CO2 and AQI
	if r.ValidityError != "" {
		hbar.PrintText(x, y, r.ValidityError)
	} else {
		x = hbar.Draw(x, y, r.CO2Index(), "C")
		hbar.Draw(x+barSpacing, y, int16(r.AQI-1), "A")
	}

	x = 0
	y = 8
	// Print second line with temperature and humidity comfort index
	x = hbar.Draw(x, y, r.CalculateComfortIndex(), "T")
	hbar.Draw(x+barSpacing, y, r.HumidityComfortIndex(), "H")

	x = 0
	y = 16
	// Print CO2 value
	if r.ValidityError == "" {
		co2 := fmt.Sprintf("%d", r.CO2)
		font16.Print(x, y, co2)
	}

	// Print humidity
	hum := fmt.Sprintf("%.0f", r.Humidity)
	x = 128 - font16.CalcWidth(hum)
	font16.Print(x, y, hum)

	// Print temperature
	temp := fmt.Sprintf("%.0f", r.Temperature)
	x = int16(x - 10 - font16.CalcWidth(temp))
	font16.Print(x, y, temp)

	f.display.Display()
}
