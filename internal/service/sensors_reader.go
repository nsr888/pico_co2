package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"tinygo.org/x/drivers/ds3231"

	"pico_co2/internal/airquality"
	"pico_co2/internal/display"
	"pico_co2/internal/types"
	"pico_co2/pkg/ens160"
)

var ErrAirQualityReadError = errors.New("air quality sensor read error")

type SensorReader struct {
	airSensor airquality.AirQualitySensor
	ds3231    *ds3231.Device
	display   *display.FontDisplay
}

func NewSensorReader(
	airSensor airquality.AirQualitySensor,
	ds3231 *ds3231.Device,
	display *display.FontDisplay,
) *SensorReader {
	return &SensorReader{
		airSensor: airSensor,
		ds3231:    ds3231,
		display:   display,
	}
}

func (sr *SensorReader) Read() (types.Readings, error) {
	var r types.Readings

	dt, err := sr.ds3231.ReadTime()
	if err != nil {
		log.Printf("Error reading time: %v", err)
	}
	r.Timestamp = dt

	if err := sr.airSensor.Read(); err != nil {
		// Read temp/hum anyway if air quality part fails
		r.Temperature = sr.airSensor.Temperature()
		r.Humidity = sr.airSensor.Humidity()
		return r, fmt.Errorf("%w: %v", ErrAirQualityReadError, err)
	}

	r.Temperature = sr.airSensor.Temperature()
	r.Humidity = sr.airSensor.Humidity()
	r.CO2 = sr.airSensor.CO2()
	r.Status = ens160.CO2String(sr.airSensor.CO2())

	return r, nil
}

func (sr *SensorReader) ProcessSensorReadings() {
	readings, err := sr.Read()
	logger := log.New(log.Writer(), readings.Timestamp.Format(time.RFC3339)+" ", 0)
	logger.Printf("Readings: %+v", readings)
	switch {
	case err != nil && !errors.Is(err, ErrAirQualityReadError):
		logger.Panicf("Error reading sensors: %v", err)
	case errors.Is(err, ErrAirQualityReadError):
		logger.Println(err)
		sr.display.DisplayBasic(readings)
	case readings.CO2 == 0 && readings.Temperature != 0:
		logger.Println("CO2 readings are zero, displaying temperature data only")
		sr.display.DisplayBasic(readings)
	default:
		sr.display.DisplayFull(readings)
	}
}
