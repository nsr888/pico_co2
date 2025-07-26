package airquality

import "pico_co2/internal/types"

// AirQualitySensor defines the standard interface for a sensor module that
// provides environmental readings.
type AirQualitySensor interface {
	Configure() error
	Read() (*types.ENSRawReadings, error)
}
