package airquality

import (
	"errors"
	"fmt"
	"time"

	"machine"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/aht20"

	"pico_co2/internal/types"
	"pico_co2/pkg/ens160"
)

// ENS160AHT20Adapter adapts the combination of an ENS160 and AHT20 sensor
// to the AirQualitySensor interface.
type ENS160AHT20Adapter struct {
	aht20     *aht20.Device
	ens160    *ens160.Device
	startTime int64
}

// Verify that ENS160AHT20Adapter implements the AirQualitySensor interface.
var _ AirQualitySensor = (*ENS160AHT20Adapter)(nil)

// NewENS160AHT20Adapter creates a new composite sensor adapter.
func NewENS160AHT20Adapter(bus *machine.I2C) *ENS160AHT20Adapter {
	aht20Device := aht20.New(bus)
	return &ENS160AHT20Adapter{
		aht20:     &aht20Device,
		ens160:    ens160.New(bus, ens160.DefaultAddress),
		startTime: time.Now().UnixMilli(),
	}
}

// Configure initializes both underlying sensors.
func (a *ENS160AHT20Adapter) Configure() error {
	a.aht20.Reset()
	a.aht20.Configure()

	if err := a.ens160.Configure(); err != nil {
		return errors.New("failed to configure ENS160: " + err.Error())
	}

	return nil
}

// Read performs a sequential read: first the AHT20 to get temperature and
// humidity, then the ENS160 using those values for compensation.
func (a *ENS160AHT20Adapter) Read() (*types.ENSRawReadings, error) {
	if err := a.aht20.Read(); err != nil {
		return nil, errors.New("failed to read AHT20: " + err.Error())
	}

	temp := a.aht20.Celsius()
	hum := a.aht20.RelHumidity()

	// Convert to integer representation for ENS160
	tempMilliC := int32(temp * 1000)
	humidityMilliPct := int32(hum * 1000)

	// Environmental compensation (ENS160)
	if err := a.ens160.SetEnvDataMilli(tempMilliC, humidityMilliPct); err != nil {
		return nil, errors.New("failed to set environmental data for ENS160: " + err.Error())
	}

	if err := a.ens160.Update(drivers.Concentration); err != nil {
		return nil, errors.New("failed to update ENS160 readings: " + err.Error())
	}

	co2 := a.ens160.ECO2()
	aqi := a.ens160.AQI()
	tvoc := a.ens160.TVOC()
	validity := a.ens160.Validity()
	validityStr := a.ValidityError(validity)

	return &types.ENSRawReadings{
		AQI:                 aqi,
		Temperature:         temp,
		Humidity:            hum,
		CO2:                 co2,
		TVOC:                tvoc,
		DataValidityWarning: validityStr,
	}, nil
}

// ValidityError returns a human-readable error message based on the ENS160
// validity status.
func (a *ENS160AHT20Adapter) ValidityError(validity uint8) string {
	switch validity {
	case ens160.ValidityNormalOperation:
		return ""
	case ens160.ValidityWarmUpPhase:
		const waitWarmUpMinutes = 3
		now := time.Now().UnixMilli()
		minutesFromStart := ((now - a.startTime) / 60000) // Convert to minutes
		minutesLeft := waitWarmUpMinutes - minutesFromStart
		minutesLeftStr := fmt.Sprintf("~%d min", minutesLeft)
		if minutesLeft <= 0 {
			minutesLeftStr = ""
		}

		return fmt.Sprintf("WARM-UP %s", minutesLeftStr)
	case ens160.ValidityInitialStartUpPhase:
		const waitStartUpMinutes = 60
		minutesFromStart := (a.startTime / 60000) // Convert to minutes
		minutesLeft := waitStartUpMinutes - minutesFromStart
		minutesLeftStr := fmt.Sprintf("~%d min", minutesLeft)
		if minutesLeft <= 0 {
			minutesLeftStr = ""
		}
		return fmt.Sprintf("START-UP %s", minutesLeftStr)
	case ens160.ValidityInvalidOutput:
		return "INVALID OUTPUT"
	default:
		return "UNKNOWN VALIDITY"
	}
}
