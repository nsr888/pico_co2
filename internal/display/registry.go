package display

import (
	"pico_co2/internal/types"
)

type RenderMethod struct {
	Name string
	Fn   func(Renderer, *types.Readings)
}

var MethodRegistry = []RenderMethod{
	{"RenderCO2BarWithNums", RenderCO2BarWithNums},
	{"RenderBars", RenderBars},
	{"RenderBarsWithLargeNums", RenderBarsWithLargeNums},
	{"RenderBasic", RenderBasic},
	{"RenderCO2Graph", RenderCO2Graph},
	{"RenderError", RenderError},
	{"RenderHeatIndexStatus", RenderHeatIndexStatus},
	{"RenderLargeBar", RenderLargeBar},
	{"RenderNums", RenderNums},
	{"RenderSparkline", RenderSparkline},
	{"RenderTempHumid", RenderTempHumid},
}
