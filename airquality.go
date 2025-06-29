package main

import (
	"fmt"
	"machine"

	"pico_co2/pkg/ens160"
	"tinygo.org/x/drivers/aht20"
)

// AirQualitySensor defines the standard interface for a sensor module that
// provides environmental readings.
type AirQualitySensor interface {
	Configure() error
	Read() error
	Temperature() float32
	Humidity() float32
	CO2() uint16
	TVOC() uint16
	AQI() uint8
}

// ENS160AHT20Adapter adapts the combination of an ENS160 and AHT20 sensor
// to the AirQualitySensor interface.
type ENS160AHT20Adapter struct {
	aht20    *aht20.Device
	ens160   *ens160.Device
	lastTemp float32
	lastHum  float32
}

// NewENS160AHT20Adapter creates a new composite sensor adapter.
func NewENS160AHT20Adapter(bus *machine.I2C) *ENS160AHT20Adapter {
	aht20Device := aht20.New(bus)
	return &ENS160AHT20Adapter{
		aht20:  &aht20Device,
		ens160: ens160.New(bus, ens160.DefaultAddress),
	}
}

// Configure initializes both underlying sensors.
func (a *ENS160AHT20Adapter) Configure() error {
	a.aht20.Reset()
	if err := a.aht20.Configure(); err != nil {
		return fmt.Errorf("failed to configure AHT20: %w", err)
	}
	if err := a.ens160.Configure(); err != nil {
		return fmt.Errorf("failed to configure ENS160: %w", err)
	}
	return nil
}

// Read performs a sequential read: first the AHT20 to get temperature and
// humidity, then the ENS160 using those values for compensation.
func (a *ENS160AHT20Adapter) Read() error {
	if err := a.aht20.Read(); err != nil {
		return fmt.Errorf("failed to read from AHT20: %w", err)
	}
	a.lastTemp = a.aht20.Celsius()
	a.lastHum = a.aht20.RelHumidity()

	if err := a.ens160.SetEnvData(a.lastTemp, a.lastHum); err != nil {
		return fmt.Errorf("failed to set env data for ENS160: %w", err)
	}
	return a.ens160.Read(ens160.WithValidityCheck(), ens160.WithWaitForNew())
}

// Temperature returns the last measured temperature.
func (a *ENS160AHT20Adapter) Temperature() float32 {
	return a.lastTemp
}

// Humidity returns the last measured humidity.
func (a *ENS160AHT20Adapter) Humidity() float32 {
	return a.lastHum
}

// CO2 returns the last measured eCO2 value.
func (a *ENS160AHT20Adapter) CO2() uint16 {
	return a.ens160.LastCO2()
}

// TVOC returns the last measured TVOC value.
func (a *ENS160AHT20Adapter) TVOC() uint16 {
	return a.ens160.LastTVOC()
}

// AQI returns the last measured AQI value.
func (a *ENS160AHT20Adapter) AQI() uint8 {
	return a.ens160.LastAQI()
}
