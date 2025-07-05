package service

import (
	"fmt"
	"log"
	"time"

	"tinygo.org/x/drivers/ds3231"

	"pico_co2/internal/airquality"
	"pico_co2/internal/display"
)

type SensorReader struct {
	airSensor airquality.AirQualitySensor
	ds3231    *ds3231.Device
	display   *display.FontDisplay
}

type SensorReaderOption func(*SensorReader)

func WithDS3231(ds3231 *ds3231.Device) SensorReaderOption {
	return func(sr *SensorReader) {
		sr.ds3231 = ds3231
	}
}

func NewSensorReader(
	airSensor airquality.AirQualitySensor,
	display *display.FontDisplay,
	sensorReaderOptions ...SensorReaderOption,
) *SensorReader {
	sr := &SensorReader{
		airSensor: airSensor,
		ds3231:    nil,
		display:   display,
	}
	for _, opt := range sensorReaderOptions {
		opt(sr)
	}

	return sr
}

func (sr *SensorReader) ProcessSensorReadings() {
	logger := log.New(log.Writer(), "SensorReader: ", log.LstdFlags)

	readings, err := sr.airSensor.Read()
	if err != nil {
		log.Println("Error reading air quality sensor:", err)
		sr.display.DisplayError(fmt.Sprintf("air sensor: %s", err.Error()))
		return
	}

	if sr.ds3231 != nil {
		dt, err := sr.ds3231.ReadTime()
		if err != nil {
			log.Println("Error reading time from DS3231:", err)
			sr.display.DisplayError(fmt.Sprintf("DS3231: %s", err.Error()))
			return
		}
		logger = log.New(log.Writer(), dt.Format(time.RFC3339)+" ", 0)
	}

	logger.Println("Sensor readings:", readings)
	sr.display.DisplayTextReadings(readings)
}
