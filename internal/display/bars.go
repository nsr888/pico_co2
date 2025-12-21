package display

import (
	"fmt"
	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderBars(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	var (
		lineY int16 = 0
		sf          = renderer.GetFont(font.ProggySZ8)
	)

	renderer.DrawTwoSideBar(
		36,
		lineY,
		int16(status.CO2Index(r.Raw.CO2)),
		"CO2",
		0,
		4,
	)
	co2 := fmt.Sprintf("%d", r.Raw.CO2)
	sf.Print(128-renderer.CalcSmallTextWidth(co2), lineY, co2)

	lineY = 11
	renderer.DrawTwoSideBar(
		0,
		lineY,
		int16(status.CalculateComfortIndex(r.Raw.Temperature, r.Raw.Humidity)),
		"TEM",
		3,
		4,
	)
	tem := fmt.Sprintf("%.0f", r.Raw.Temperature)
	sf.Print(128-renderer.CalcSmallTextWidth(tem), lineY, tem)

	lineY = 22
	renderer.DrawTwoSideBar(
		0,
		lineY,
		int16(status.HumidityComfortIndex(r.Raw.Humidity)),
		"HUM",
		3,
		4,
	)
	hum := fmt.Sprintf("%.0f", r.Raw.Humidity)
	sf.Print(128-renderer.CalcSmallTextWidth(hum), lineY, hum)

	renderer.Display()
}
