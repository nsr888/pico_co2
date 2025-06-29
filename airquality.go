package main

import (
	"machine"

	"pico_co2/pkg/ens160"
)

// AirQualitySensor defines the standard interface for an air quality sensor.
type AirQualitySensor interface {
	Configure() error
	// ReadAndCompensate triggers a sensor reading, using temperature and
	// humidity for compensation if supported by the sensor.
	ReadAndCompensate(temperature, humidity float32) error
	CO2() uint16
	TVOC() uint16
	AQI() uint8
}

// ENS160Adapter adapts the ens160.Device to the AirQualitySensor interface.
type ENS160Adapter struct {
	device *ens160.Device
}

// NewENS160 creates a new adapter for the ENS160 sensor.
func NewENS160(bus *machine.I2C) *ENS160Adapter {
	return &ENS160Adapter{
		device: ens160.New(bus, ens160.DefaultAddress),
	}
}

// Configure initializes the sensor.
func (a *ENS160Adapter) Configure() error {
	return a.device.Configure()
}

// ReadAndCompensate sets environment data and reads the sensor.
func (a *ENS160Adapter) ReadAndCompensate(temperature, humidity float32) error {
	if err := a.device.SetEnvData(temperature, humidity); err != nil {
		return err
	}
	return a.device.Read(ens160.WithValidityCheck(), ens160.WithWaitForNew())
}

// CO2 returns the last measured eCO2 value.
func (a *ENS160Adapter) CO2() uint16 {
	return a.device.LastCO2()
}

// TVOC returns the last measured TVOC value.
func (a *ENS160Adapter) TVOC() uint16 {
	return a.device.LastTVOC()
}

// AQI returns the last measured AQI value.
func (a *ENS160Adapter) AQI() uint8 {
	return a.device.LastAQI()
}
