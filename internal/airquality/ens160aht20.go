package airquality

import (
	"fmt"
	"machine"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/aht20"

	"pico_co2/internal/types"
	"pico_co2/pkg/ens160"
)

// ENS160AHT20Adapter adapts the combination of an ENS160 and AHT20 sensor
// to the AirQualitySensor interface.
type ENS160AHT20Adapter struct {
	aht20    *aht20.Device
	ens160   *ens160.Device
	readings *types.Readings
}

// Verify that ENS160AHT20Adapter implements the AirQualitySensor interface.
var _ AirQualitySensor = (*ENS160AHT20Adapter)(nil)

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
	a.aht20.Configure()

	if err := a.ens160.Configure(); err != nil {
		return fmt.Errorf("failed to configure ENS160: %w", err)
	}

	return nil
}

// Read performs a sequential read: first the AHT20 to get temperature and
// humidity, then the ENS160 using those values for compensation.
func (a *ENS160AHT20Adapter) Read() (*types.Readings, error) {
	if err := a.aht20.Read(); err != nil {
		return nil, fmt.Errorf("failed to read AHT20: %w", err)
	}

	temp := a.aht20.Celsius()
	hum := a.aht20.RelHumidity()

	// Convert to integer representation for ENS160
	tempMilliC := int32(temp * 1000)
	humidityMilliPct := int32(hum * 1000)

	// Environmental compensation (ENS160)
	if err := a.ens160.SetEnvDataMilli(tempMilliC, humidityMilliPct); err != nil {
		return nil, fmt.Errorf("failed to set environment data for ENS160: %w", err)
	}

	if err := a.ens160.Update(drivers.Concentration); err != nil {
		return nil, fmt.Errorf("failed to read ENS160: %w", err)
	}

	co2 := a.ens160.ECO2()
	r := &types.Readings{
		Temperature: temp,
		Humidity:    hum,
		CO2:         co2,
		CO2String:   CO2String(co2),
		IsValid:     true,
	}

	return r, nil
}

func CO2String(value uint16) string {
	switch {
	case value < 400:
		return "No data"
	case value < 600:
		return "Excellent"
	case value < 800:
		return "Good"
	case value < 1000:
		return "Fair"
	case value < 1500:
		return "Poor"
	default:
		return "Bad"
	}
}
