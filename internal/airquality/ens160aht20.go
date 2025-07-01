package airquality

import (
	"errors"
	"fmt"

	"machine"
	"tinygo.org/x/drivers/aht20"

	"pico_co2/internal/types"
	"pico_co2/pkg/ens160"
)

const (
	StatusStartUp   = "Start up"
	StatusWarmingUp = "Warming up"
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
func (a *ENS160AHT20Adapter) Read() (*Readings, error) {
	var r Readings

	if err := a.aht20.Read(); err != nil {
		return nil, fmt.Errorf("read from AHT20: %w", err)
	}
	r.Temperature = a.Temperature()
	r.Humidity = a.Humidity()

	if err := a.ens160.SetEnvData(
		r.Temperature,
		r.Humidity,
	); err != nil {
		return &r, fmt.Errorf("set env data for ENS160: %w", err)
	}

	err := a.ens160.Read(ens160.WithValidityCheck(), ens160.WithWaitForNew())
	switch {
	case errors.Is(err, ens160.ErrInitialStartUpPhase) ||
		errors.Is(err, ens160.ErrWarmUpPhase):
		r.CO2 = a.CO2()
		r.CO2Interpretation = ens160.CO2String(r.CO2)
		r.Quality.Status = StatusNotValid
		r.Quality.Description = fmt.Sprintf("read ENS160: %s", err.Error())

		return &r, nil
	case err != nil:
		r.Quality.Status = StatusNotValid
		r.Quality.Description = fmt.Sprintf("read ENS160: %s", err.Error())

		return &r, fmt.Errorf("read ENS160: %w", err)
	}

	r.CO2 = a.CO2()
	r.CO2Interpretation = ens160.CO2String(r.CO2)
	r.Quality.Status = StatusOK

	return &r, nil
}

// Temperature returns the last measured temperature.
func (a *ENS160AHT20Adapter) Temperature() float32 {
	return a.aht20.Celsius()
}

// Humidity returns the last measured humidity.
func (a *ENS160AHT20Adapter) Humidity() float32 {
	return a.aht20.RelHumidity()
}

// CO2 returns the last measured eCO2 value.
func (a *ENS160AHT20Adapter) CO2() uint16 {
	return a.ens160.LastCO2()
}
