package airquality

// AirQualitySensor defines the standard interface for a sensor module that
// provides environmental readings.
type AirQualitySensor interface {
	Configure() error
	Read() error
	Temperature() float32
	Humidity() float32
	CO2() uint16
}
