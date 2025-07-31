package display

import (
	"sort"

	"pico_co2/internal/types"
)

// DisplayMethod defines a type for the display method keys to improve type safety.
type DisplayMethod string

// Constants for all available display methods.
const (
	AqiBarWithNums    DisplayMethod = "aqi_bar_with_nums"
	TempHumid         DisplayMethod = "temp_humid"
	Error             DisplayMethod = "error"
	Basic             DisplayMethod = "basic"
	CO2Graph          DisplayMethod = "co2_graph"
	Bars              DisplayMethod = "bars"
	BarsWithLargeNums DisplayMethod = "bars_with_large_nums"
	HeatIndexStatus   DisplayMethod = "render_heat_index_status"
	Sparkline         DisplayMethod = "sparkline"
	LargeBar          DisplayMethod = "large_bar"
	Nums              DisplayMethod = "nums"
)

func (d DisplayMethod) String() string {
	return string(d)
}

var MethodRegistry = map[DisplayMethod]func(Renderer, *types.Readings){
	AqiBarWithNums:    RenderAqiBarWithNums,
	TempHumid:         RenderTempHumid,
	Error:             RenderError,
	Basic:             RenderBasic,
	CO2Graph:          RenderCO2Graph,
	Bars:              RenderBars,
	BarsWithLargeNums: RenderBarsWithLargeNums,
	HeatIndexStatus:   RenderHeatIndexStatus,
	Sparkline:         RenderSparkline,
	LargeBar:          RenderLargeBar,
	Nums:              RenderNums,
}

// GetAllDisplayMethods returns a sorted slice of all display methods
func GetAllDisplayMethods() []DisplayMethod {
	methods := make([]DisplayMethod, 0, len(MethodRegistry))

	for method := range MethodRegistry {
		// Exclude Error method as it's only for error conditions
		if method != Error {
			methods = append(methods, method)
		}
	}

	// Sort alphabetically
	sort.Slice(methods, func(i, j int) bool {
		return methods[i] < methods[j]
	})

	return methods
}
