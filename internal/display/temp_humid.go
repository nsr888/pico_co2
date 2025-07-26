package display

import (
	"fmt"
	"image/color"

	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderTempHumid(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	var (
		xPos int16
		yPos int16
	)

	var verticalBarWidth int16 = 4
	var spacing int16 = 20

	temp := fmt.Sprintf("%.0f", r.Raw.Temperature)
	tempWidth := renderer.CalcXLargeTextWidth(temp)
	xPos = int16(0)
	yPos = int16(8)
	renderer.DrawXLargeText(xPos, yPos, temp)
	// TODO: move to driver
	DrawVerticalBar(renderer, tempWidth+4, yPos, int16(r.Calculated.CO2Status))

	hum := fmt.Sprintf("%.0f", r.Raw.Humidity)
	humWidth := renderer.CalcXLargeTextWidth(hum)
	xPos = tempWidth + verticalBarWidth + spacing
	renderer.DrawXLargeText(xPos, yPos, hum)
	DrawVerticalBar(renderer, xPos+humWidth+4, yPos, status.HumidityComfortIndex(r.Raw.Humidity))

	xPos = int16(0)
	yPos = int16(0)
	renderer.DrawSmallText(xPos, yPos, "Temp")
	xPos = tempWidth + verticalBarWidth + spacing
	renderer.DrawSmallText(xPos, yPos, "Humidity")

	renderer.Display()
}

func DrawVerticalBar(renderer Renderer, x, y, filledBars int16) {
	if renderer == nil {
		return
	}

	black := color.RGBA{255, 255, 255, 255}

	const maxBars = 4
	y += 5
	unfilledBars := maxBars - filledBars

	coef := func(i int16) int16 {
		return (i * 5)
	}

	// Draw filled bars
	for i := int16(0); i < unfilledBars; i++ {
		renderer.SetPixel(x, y+coef(i), black)
	}

	// Draw empty bars
	for i := unfilledBars; i < 4; i++ {
		drawBlock(renderer, x, y+coef(i), x+3, y+coef(i), black)
	}
}

func drawBlock(renderer Renderer, x0, y0, x1, y1 int16, color color.RGBA) {
	// Simple line drawing implementation
	dx := x1 - x0
	dy := y1 - y0
	steps := max(abs(dx), abs(dy))

	for i := int16(0); i <= steps; i++ {
		x := x0 + (dx*i)/steps
		y := y0 + (dy*i)/steps
		renderer.SetPixel(x, y, color)
		renderer.SetPixel(x, y+1, color)
		renderer.SetPixel(x, y+2, color)
	}
}

func abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}
