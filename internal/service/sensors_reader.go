package service

import (
	"fmt"
	"log"
	"time"

	"tinygo.org/x/drivers/ds3231"

	"pico_co2/internal/airquality"
	"pico_co2/internal/display"
	"pico_co2/internal/types"
)

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

func (sr *SensorReader) ProcessSensorReadings() {
	readings, err := sr.ReadAll()
	if err != nil {
		log.Println("Critical sensor error:", err)
		sr.display.DisplayError(fmt.Sprintf("Fatal: %v", err))
		return
	}

	logger := log.New(log.Writer(), readings.Timestamp.Format(time.RFC3339)+" ", 0)
	logger.Println("Sensor readings:", readings)

	// Show detailed descriptions when status is not OK
	if readings.Description != "" {
		sr.display.DisplayError(readings.Description)
	} else {
		sr.display.DisplayReadings(readings)
	}
}

func (sr *SensorReader) ReadAll() (*types.Readings, error) {
	// Collect RTC time first
	dt, err := sr.ds3231.ReadTime()
	if err != nil {
		return nil, fmt.Errorf("RTC error: %w", err)
	}

	// Collect air quality sensor data
	airReadings, err := sr.airSensor.Read()
	if err != nil {
		return nil, fmt.Errorf("air sensor error: %w", err)
	}

	return &types.Readings{
		Timestamp:   dt,
		Temperature: airReadings.Temperature,
		Humidity:    airReadings.Humidity,
		CO2:         airReadings.CO2,
		Status:      airReadings.Interpretation(),
		Description: airReadings.Quality.Description,
	}, nil
}
