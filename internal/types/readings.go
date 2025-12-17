package types

import (
	"pico_co2/internal/types/status"
	"pico_co2/pkg/fifo"
	"time"
)

type Readings struct {
	Raw RawReadings `json:"raw_readings"`

	FirstReadingAt time.Time          `json:"first_reading_time"`
	Calculated       CalculatedReadings `json:"calculated_readings"`
	History          MeasurementHistory `json:"co2_history"`
	LastUpdateAt     time.Time          `json:"created_at"`
	IsDrawen         bool               `json:"is_drawn"`
	LastRaw          RawReadings        `json:"last_raw"`

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
	CO2         *fifo.FIFO16  `json:"co2"`
	Temperature *fifo.FIFO16  `json:"temperature"`
	Humidity    *fifo.FIFO16  `json:"humidity"`
	AddedAt     time.Time     `json:"added_at"`
	Granularity time.Duration `json:"granularity"`
}

type CalculatedReadings struct {
	HeatIndex       status.HeatIndexStatus `json:"heat_index,omitempty"`
	CO2Status       status.CO2Index        `json:"eco2_human,omitempty"`
	CO215MinAverage uint16                 `json:"co2_15min_average,omitempty"`
}

func InitReadings(queueSize int) *Readings {
	return &Readings{
		History: MeasurementHistory{
			CO2:         fifo.NewFIFO16(queueSize),
			Temperature: fifo.NewFIFO16(queueSize),
			Humidity:    fifo.NewFIFO16(queueSize),
			Granularity: time.Minute,
		},
	}
}

func (r *Readings) AddReadings(
	co2 uint16,
	temperature float32,
	humidity float32,
) {
	r.Error = ""
	r.LastUpdateAt = time.Now()

	if r.FirstReadingAt.IsZero() {
		r.FirstReadingAt = time.Now()
	}

	if r.History.CO2 == nil {
		return
	}

	if time.Since(r.History.AddedAt) > r.History.Granularity {
		r.History.CO2.Enqueue(int16(co2))
		r.History.Temperature.Enqueue(int16(temperature))
		r.History.Humidity.Enqueue(int16(humidity))
		r.History.AddedAt = time.Now()
	}

	heatIndex := status.HeatIndex(temperature, humidity)
	r.Calculated.HeatIndex = status.ToHeatIndexStatus(heatIndex)
	r.Calculated.CO2Status = status.ToCO2Index(co2)

	// Calculate 15-minute average of last 15 readings
	if r.History.CO2.Len() >= 15 {
		var sum uint32
		count := 0

		r.History.CO2.PeekAll(func(val int16) {
			sum += uint32(val)
			count++
			if count >= 15 {
				return
			}
		})

		if count > 0 {
			r.Calculated.CO215MinAverage = uint16(sum / uint32(count))
		}
	} else {
		// Not enough data yet, use current reading as initial average
		if r.Calculated.CO215MinAverage == 0 {
			r.Calculated.CO215MinAverage = co2
		}
	}

	// Store last measurements before updating with new ones
	r.LastRaw = r.Raw

	r.Raw = RawReadings{
		CO2:         co2,
		Temperature: temperature,
		Humidity:    humidity,
	}
}

func (r *Readings) MeasurementsChanged() bool {
	// If this is the first reading, consider it as changed
	if r.LastRaw.CO2 == 0 && r.LastRaw.Temperature == 0 &&
		r.LastRaw.Humidity == 0 {
		return true
	}

	return r.Raw.CO2 != r.LastRaw.CO2 ||
		r.Raw.Temperature != r.LastRaw.Temperature ||
		r.Raw.Humidity != r.LastRaw.Humidity
}
