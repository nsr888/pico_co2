package display

import (
	"fmt"
	"image/color"
	"log"

	font "github.com/Nondzu/ssd1306_font"
	"tinygo.org/x/tinydraw"

	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayReadingsWithHI(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	if r.ValidityError != "" {
		f.font.Configure(font.Config{FontType: font.FONT_7x10})
		f.font.YPos = 0
		f.font.XPos = 0
		f.font.XPos = 0
		f.font.PrintText(r.ValidityError)
	} else {
		// CO2Status
		f.font.Configure(font.Config{FontType: font.FONT_7x10})
		status := "CO2"
		f.font.YPos = 0
		f.font.XPos = 0
		f.font.PrintText(status)

		f.DrawBar(30, 4, r.CO2Rating())

		// CO2 value
		f.font.Configure(font.Config{FontType: font.FONT_11x18})
		co2Value := fmt.Sprintf("%d", r.CO2)
		f.font.YPos = 0
		f.font.XPos = int16(128 - (len(co2Value) * 11))
		f.font.PrintText(co2Value)
	}

	// Heat Index status
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	hiStatus := "HI"
	f.font.YPos = 10
	f.font.XPos = 0
	f.font.PrintText(hiStatus)

	f.DrawBar(30, 14, r.HeatIndexRating())

	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	tempHum := fmt.Sprintf("%.0f %.0f", r.Temperature, r.Humidity)
	f.font.YPos = 16
	f.font.XPos = int16(128 - (len(tempHum) * 11))
	f.font.PrintText(tempHum)

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
func (f *FontDisplay) DrawBar(x, y, filledBars int16) {
	if f == nil {
		return
	}

	black := color.RGBA{1, 1, 1, 255}

	x = x + 2

	// Draw filled bars
	for i := int16(0); i < filledBars; i++ {
		radius := int16(3)
		tinydraw.FilledCircle(f.display, x+i*10, y, radius, black)
	}

	// Draw empty bars
	for i := filledBars; i < 4; i++ {
		radius := int16(1)
		tinydraw.FilledCircle(f.display, x+i*10, y, radius, black)
	}
}
