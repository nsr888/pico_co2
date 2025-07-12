package display

import (
	"fmt"
	"image/color"
	"log"

	"tinygo.org/x/tinydraw"

	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayReadingsWithHI(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	radiusFilled := int16(3)

	font7 := font.NewFont7(f.display)
	if r.ValidityError != "" {
		font7.Print(0, 0, r.ValidityError)
	} else {
		// CO2Status
		status := "CO2"
		font7.Print(0, 0, status)

		f.DrawBar(30, 4, r.CO2Index(), radiusFilled)

		// CO2 value
		co2Value := fmt.Sprintf("%d", r.CO2)
		YPos := int16(0)
		XPos := int16(128 - (len(co2Value) * 11))
		font7.Print(XPos, YPos, co2Value)
	}

	// Heat Index status
	hiStatus := "HI"
	var YPos int16 = 10
	var XPos int16 = 0
	font7.Print(XPos, YPos, hiStatus)

	f.DrawBar(30, 14, r.HeatIndexRating(), radiusFilled)

	font11 := font.NewFont11(f.display)

	tempHum := fmt.Sprintf("%.0f %.0f", r.Temperature, r.Humidity)
	YPos = 16
	XPos = int16(128 - (len(tempHum) * 11))
	font11.Print(XPos, YPos, tempHum)

	f.humFIFO.Enqueue(int16(r.Humidity))
	f.GraphHumidity(0, 18)
}

// GraphHumidity draws a vertical bar graph of the last 64 humidity readings
// x and y are the top-left corner of the graph area.
func (f *FontDisplay) GraphHumidity(x, y int16) {
	if f == nil {
		return
	}
	black := color.RGBA{1, 1, 1, 255}

	// tinydraw.Rectangle(f.display, x, y, graphWidth, graphHeight, black)

	y = y + graphHeight

	raw := f.humFIFO.CopyTo()
	if len(raw) < 3 {
		log.Printf("Not enough data to draw graph: %d readings\n", len(raw))
		return
	}

	result := f.sparkline.Process(raw)

	for _, hum := range result {
		barHeight := (hum * graphHeight) / 100
		topY := y - barHeight

		// Draw the bar
		f.display.FillRectangle(x, topY, barWidth, barHeight, black)
		// log.Printf("geo: x=%d, y=%d, toY=%d, width=%d, height=%d\n", x, y, topY, barWidth, barHeight)

		x += 1 // Move to the next bar position
		if x < 0 {
			return // Stop drawing if we go out of bounds
		}
	}

	f.display.Display()
}

// max bars is 4
func (f *FontDisplay) DrawBar(x, y, filledBars int16, filledRadius int16) {
	if f == nil {
		return
	}

	black := color.RGBA{1, 1, 1, 255}

	x = x + 2

	// Draw filled bars
	for i := int16(0); i < filledBars; i++ {
		tinydraw.FilledCircle(f.display, x+i*10, y, filledRadius, black)
	}

	// Draw empty bars
	for i := filledBars; i < 4; i++ {
		radius := int16(1)
		tinydraw.FilledCircle(f.display, x+i*10, y, radius, black)
	}
}
