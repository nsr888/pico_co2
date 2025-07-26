package types

import (
	"time"

	"pico_co2/internal/types/status"
	"pico_co2/pkg/fifo"
)

type Readings struct {
	Raw RawReadings `json:"raw_readings"`

	FirstReadingTime time.Time          `json:"first_reading_time"`
	Calculated       CalculatedReadings `json:"calculated_readings"`
	CO2History       MeasurementHistory `json:"co2_history"`

	Warning string `json:"warning,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RawReadings struct {
	CO2         uint16  `json:"eco2"`
	TVOC        uint16  `json:"tvoc"`
	AQI         uint8   `json:"aqi"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

type MeasurementHistory struct {
	Measurements *fifo.FIFO16  `json:"measurements"`
	AddedAt      time.Time     `json:"added_at"`
	Granularity  time.Duration `json:"granularity"`
}

type CalculatedReadings struct {
	HeatIndex status.HeatIndexStatus `json:"heat_index,omitempty"`
	CO2Status status.CO2Index        `json:"eco2_human,omitempty"`
}

func InitReadings(queueSize int) *Readings {
	return &Readings{
		CO2History: MeasurementHistory{
			Measurements: fifo.NewFIFO16(queueSize),
			Granularity:  time.Minute,
		},
	}
}

func (r *Readings) AddReadings(
	co2 uint16,
	tvoc uint16,
	aqi uint8,
	temperature float32,
	humidity float32,
	warning string,
) {
	// Initialize the first reading time if not set
	if r.FirstReadingTime.IsZero() {
		r.FirstReadingTime = time.Now()
	}

	// Update CO2 history if measurements are initialized
	if r.CO2History.Measurements == nil {
		return
	}

	if time.Since(r.CO2History.AddedAt) > r.CO2History.Granularity {
		r.CO2History.Measurements.Enqueue(int16(co2))
		r.CO2History.AddedAt = time.Now()
	}

	// Calculate derived readings
	heatIndex := status.HeatIndex(temperature, humidity)
	r.Calculated.HeatIndex = status.ToHeatIndexStatus(heatIndex)
	r.Calculated.CO2Status = status.ToCO2Index(co2)

	r.Warning = warning

	// Update the readings
	r.Raw = RawReadings{
		CO2:         co2,
		TVOC:        tvoc,
		AQI:         aqi,
		Temperature: temperature,
		Humidity:    humidity,
	}
}
