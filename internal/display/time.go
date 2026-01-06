package display

import (
	"fmt"
	"math"
	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderTime(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	var (
		y         int16 = 1
		x         int16
		co2status int16
		lf        = renderer.GetFont(font.FreemonoRegular18)
		sf        = renderer.GetFont(font.ProggySZ8)
	)

	width, _ := renderer.Size()

	// First line

	heatIndex := status.GetHeatIndex(r.Raw.Temperature, r.Raw.Humidity)
	x = renderer.DrawTwoSideBar(x, y, int16(heatIndex), "H", 0, 2)

	temp := fmt.Sprintf("%.0f", math.Round(float64(r.Raw.Temperature)))
	hum := fmt.Sprintf("%.0f", math.Round(float64(r.Raw.Humidity)))
	x = width/2 - sf.CalcWidth(temp) - 2 - 1
	sf.Print(x, y, temp)
	x = width/2 + 2
	sf.Print(x, y, hum)

	// https://backend.orbit.dtu.dk/ws/portalfiles/portal/348932926/1-s2.0-S0360132323011459-main_1_.pdf
	switch {
	case r.Raw.CO2 < 800:
		co2status = 0
	case r.Raw.CO2 < 1000:
		co2status = 1
	default:
		co2status = 2
	}
	x = 97
	renderer.DrawTwoSideBar(x, y, co2status, "C", 0, 2)

	// second line
	y = 10
	timeStr := fmt.Sprintf("%d:%02d", r.Time.Hour, r.Time.Minute)
	xTime := (width - lf.CalcWidth(timeStr)) / 2
	lf.Print(xTime, y, timeStr)

	renderer.Display()
}
