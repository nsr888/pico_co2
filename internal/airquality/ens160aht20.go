package airquality

import (
	"errors"
	"fmt"

	"machine"
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

	r := &types.Readings{
		Temperature: a.aht20.Celsius(),
		Humidity:    a.aht20.RelHumidity(),
	}

	// Environmental compensation (ENS160)
	if err := a.ens160.SetEnvData(r.Temperature, r.Humidity); err != nil {
		return r, fmt.Errorf("failed to set environment data for ENS160: %w", err)
	}

	if err := a.ens160.Read(ens160.WithValidityCheck(), ens160.WithWaitForNew()); err != nil {
		r.Description = err.Error()
		if errors.Is(err, ens160.ErrInitialStartUpPhase) || errors.Is(err, ens160.ErrWarmUpPhase) {
			r.CO2 = a.ens160.LastCO2()
			r.CO2String = ens160.CO2String(r.CO2)
			r.IsValid = false

			return r, nil
		}

		return r, fmt.Errorf("failed to read ENS160: %w", err)
	}

	r.IsValid = true
	r.CO2 = a.ens160.LastCO2()
	r.CO2String = ens160.CO2String(r.CO2)

	return r, nil
}
