package main

import (
	"fmt"
	"log"
	"time"

	"machine"
	"tinygo.org/x/drivers/ds3231"

	"pico_co2/pkg/ens160"
)

// Readings represents sensor data
type Readings struct {
	CO2         uint16    `json:"eco2"`
	Temperature float32   `json:"temperature"`
	Humidity    float32   `json:"humidity"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

type SensorReader struct {
	airSensor AirQualitySensor
	ds3231    *ds3231.Device
}

func NewSensorReader(bus *machine.I2C) (*SensorReader, error) {
	airQualitySensor := NewENS160AHT20Adapter(bus)
	if err := airQualitySensor.Configure(); err != nil {
		return nil, fmt.Errorf("failed to configure air quality sensor: %w", err)
	}

	ds3231Sensor := ds3231.New(bus)
	if ok := ds3231Sensor.Configure(); !ok {
		return nil, fmt.Errorf("failed to configure DS3231 sensor")
	}

	return &SensorReader{
		airSensor: airQualitySensor,
		ds3231:    &ds3231Sensor,
	}, nil
}

func (sr *SensorReader) Read() (Readings, error) {
	var r Readings

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
