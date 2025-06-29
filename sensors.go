package main

import (
	"fmt"
	"log"
	"time"

	"pico_co2/pkg/ens160"
)

// Readings represents sensor data
type Readings struct {
	AQI         uint8     `json:"aqi"`
	ECO2        uint16    `json:"eco2"`
	TVOC        uint16    `json:"tvoc"`
	Humidity    float32   `json:"humidity"`
	Temperature float32   `json:"temperature"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

func (a *App) readSensors() (Readings, error) {
	var r Readings

	dt, err := a.ds3231Sensor.ReadTime()
	if err != nil {
		log.Printf("Error reading time: %v", err)
	}

	r.Timestamp = dt

	if a.aht20Sensor == nil {
		return r, fmt.Errorf("AHT20 sensor not initialized")
	}

	if err := a.aht20Sensor.Read(); err != nil {
		return r, fmt.Errorf("failed to read AHT20 sensor: %w", err)
	}

	r.Temperature = a.aht20Sensor.Celsius()
	r.Humidity = a.aht20Sensor.RelHumidity()

	if err := a.ens160Sensor.SetEnvData(r.Temperature, r.Humidity); err != nil {
		return r, fmt.Errorf("failed to set environment data for ENS160: %w", err)
	}

	// err = a.ens160Sensor.Read(ens160.WithValidityCheck(), ens160.WithWaitForNew())
	err = a.ens160Sensor.Read()
	if err != nil {
		return r, fmt.Errorf("%w: %v", ErrENS160ReadError, err)
	}

	r.AQI = a.ens160Sensor.LastAQI()
	r.ECO2 = a.ens160Sensor.LastCO2()
	r.TVOC = a.ens160Sensor.LastTVOC()
	r.Status = ens160.CO2String(a.ens160Sensor.LastCO2())

	return r, nil
}
