package display

import (
	"fmt"
	"math"
	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderLevel(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	width, _ := renderer.Size()

	var (
		lf = renderer.GetFont(font.FreemonoRegular12)
		sf = renderer.GetFont(font.Notosans)
		xf = renderer.GetFont(font.ProggySZ8)
		ne = renderer.GetFont(font.Notoemoji)
	)

	var (
		arrow    string
		decision string
		trend    = r.Calculated.CO2Trend
	)

	switch trend {
	case status.RisingCO2:
		arrow = "⬆"
	case status.FallingCO2:
		arrow = "⬇"
	case status.StableCO2:
		arrow = "➡"
	default:
		arrow = ""
	}

	// https://backend.orbit.dtu.dk/ws/portalfiles/portal/348932926/1-s2.0-S0360132323011459-main_1_.pdf
	switch {
	case r.Raw.CO2 <= 800:
		decision = "OK CO"
	case r.Raw.CO2 <= 1000:
		decision = "SOON CO"
	default:
		decision = "VENT CO"
	}

	lf.Print(0, 0, decision)
	decisionWidth := lf.CalcWidth(decision) + 2

	xf.Print(decisionWidth, 12, "2")
	decisionWidth += xf.CalcWidth("2") + 2

	ne.Print(decisionWidth, 0, arrow)

	// Line 2: Three metrics (small font) - Temperature, Humidity, CO2
	tempStr := fmt.Sprintf("%.0f C", math.Round(float64(r.Raw.Temperature)))
	humStr := fmt.Sprintf("%.0f %%", math.Round(float64(r.Raw.Humidity)))
	co2Str := fmt.Sprintf("%d", r.Raw.CO2)

	// Calculate widths for proper spacing
	tempWidth := sf.CalcWidth(tempStr)
	co2Width := sf.CalcWidth(co2Str)

	// Position: Temperature | Humidity | CO2
	var x1, x2, x3 int16
	x1 = 0
	x2 = x1 + tempWidth + 17
	x3 = width - co2Width

	sf.Print(x1, 22, tempStr)
	sf.Print(x2, 22, humStr)
	sf.Print(x3, 22, co2Str)

	renderer.Display()
}
