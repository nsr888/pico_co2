package main

import (
	"math"
	"os"
	"time"

	"pico_co2/internal/display"
	"pico_co2/internal/types"
)

const (
	queueCapacity       = 128 // Number of readings to keep in memory
	displayWidth  int16 = 128
	displayHeight int16 = 32
)

func simulateSensor(x uint8) uint16 {
	radians := float32(x) * 2 * math.Pi / queueCapacity
	y := 1200 + 800*float32(math.Sin(float64(radians)))
	return uint16(y)
}

func main() {
	vd := display.NewVirtualDisplay(displayWidth, displayHeight)

	testReadings := types.InitReadings(queueCapacity)
	testReadings.FirstReadingAt = testReadings.FirstReadingAt.Add(-3 * time.Minute)

	countMeasurements := queueCapacity
	for i := range countMeasurements {
		testReadings.History.AddedAt = testReadings.FirstReadingAt.Add(-2 * time.Minute)
		// use formula to generate graph data with increasing and decreasing values
		co2 := simulateSensor(uint8(i))
		temperature := 22.5 + float64(i)/10.0
		humidity := 45.0 + float64(i)/10.0
		testReadings.AddReadings(
			uint16(co2),
			float32(temperature),
			float32(humidity),
		)
	}

	testCases := []struct {
		name     string
		mode     string
		readings func(*types.Readings) *types.Readings
	}{
		{
			name: "normal",
			readings: func(r *types.Readings) *types.Readings {
				r.Warning = ""
				r.Error = ""
				return r
			},
		},
		{
			name: "warning",
			readings: func(r *types.Readings) *types.Readings {
				r.Warning = "Waiting for data..."
				r.Error = ""
				return r
			},
		},
		{
			name: "error",
			readings: func(r *types.Readings) *types.Readings {
				r.Warning = ""
				r.Error = "Test error message for display with long text that should wrap correctly across multiple lines."
				return r
			},
		},
	}

	for _, tc := range testCases {
		vd.Clear()
		for _, method := range display.MethodRegistry {
			os.MkdirAll("images/"+tc.name, 0755)
			method.Fn(vd, tc.readings(testReadings))
			vd.SavePNG("images/" + tc.name + "/" + method.Name + ".png")
		}
	}
}
