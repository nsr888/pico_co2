package main

import (
	"fmt"
	"log"
	"time"

	"machine"
	"tinygo.org/x/drivers/aht20"
	"tinygo.org/x/drivers/ds3231"

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

type SensorReader struct {
	aht20  *aht20.Device
	ens160 *ens160.Device
	ds3231 *ds3231.Device
}

func NewSensorReader(bus *machine.I2C) (*SensorReader, error) {
	aht20Sensor := aht20.New(bus)
	aht20Sensor.Reset()
	aht20Sensor.Configure()

	ens160Sensor := ens160.New(bus, ens160.DefaultAddress)
	if err := ens160Sensor.Configure(); err != nil {
		return nil, fmt.Errorf("failed to configure ENS160 sensor: %w", err)
	}

	ds3231Sensor := ds3231.New(bus)
	if ok := ds3231Sensor.Configure(); !ok {
		return nil, fmt.Errorf("failed to configure DS3231 sensor")
	}

	return &SensorReader{
		aht20:  &aht20Sensor,
		ens160: ens160Sensor,
		ds3231: &ds3231Sensor,
	}, nil
}

func (sr *SensorReader) Read() (Readings, error) {
	var r Readings

	dt, err := sr.ds3231.ReadTime()
	if err != nil {
		log.Printf("Error reading time: %v", err)
	}

	r.Timestamp = dt

	if sr.aht20 == nil {
		return r, fmt.Errorf("AHT20 sensor not initialized")
	}

	if err := sr.aht20.Read(); err != nil {
		return r, fmt.Errorf("failed to read AHT20 sensor: %w", err)
	}

	r.Temperature = sr.aht20.Celsius()
	r.Humidity = sr.aht20.RelHumidity()

	if err := sr.ens160.SetEnvData(r.Temperature, r.Humidity); err != nil {
		return r, fmt.Errorf("failed to set environment data for ENS160: %w", err)
	}

	err = sr.ens160.Read(ens160.WithValidityCheck(), ens160.WithWaitForNew())
	if err != nil {
		return r, fmt.Errorf("%w: %v", ErrENS160ReadError, err)
	}

	r.AQI = sr.ens160.LastAQI()
	r.ECO2 = sr.ens160.LastCO2()
	r.TVOC = sr.ens160.LastTVOC()
	r.Status = ens160.CO2String(sr.ens160.LastCO2())

	return r, nil
}
