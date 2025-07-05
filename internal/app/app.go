package app

import (
	"fmt"
	"log"
	"time"

	"machine"
	"tinygo.org/x/drivers/ds3231"

	"pico_co2/internal/airquality"
	"pico_co2/internal/display"
	"pico_co2/internal/service"
)

// Application Logic
const (
	sampleTimeSeconds = 60
	watchDogMillis    = 8388 // max for RP2040 is 8388ms
	i2cFrequency      = 400 * machine.KHz
	i2cSDA            = machine.GP4
	i2cSCL            = machine.GP5
)

type App struct {
	sensorReader *service.SensorReader
}

type Config struct {
	IsAdvancedSetup bool
}

func New(cfg Config) (*App, error) {
	if err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFrequency,
		SDA:       i2cSDA,
		SCL:       i2cSCL,
	}); err != nil {
		return nil, err
	}

	fontDisplay, err := display.NewFontDisplay(machine.I2C0)
	if err != nil {
		return nil, err
	}

	airQualitySensor := airquality.NewENS160AHT20Adapter(machine.I2C0)
	if err := airQualitySensor.Configure(); err != nil {
		return nil, fmt.Errorf("failed to configure air quality sensor: %w", err)
	}

	ds3231Sensor := ds3231.New(machine.I2C0)
	if ok := ds3231Sensor.Configure(); !ok {
		return nil, fmt.Errorf("failed to configure DS3231 sensor")
	}

	opts := []service.SensorReaderOption{}
	if cfg.IsAdvancedSetup {
		opts = append(opts, service.WithDS3231(&ds3231Sensor))
	}

	sensorReader := service.NewSensorReader(
		airQualitySensor,
		fontDisplay,
		opts...,
	)

	return &App{
		sensorReader: sensorReader,
	}, nil
}

// Run starts the main application loop.
func (a *App) Run() {
	wd := machine.Watchdog
	config := machine.WatchdogConfig{
		TimeoutMillis: watchDogMillis,
	}
	wd.Configure(config)
	wd.Start()
	log.Printf("starting loop")

	for {
		a.sensorReader.Process()
		waitNextSample(sampleTimeSeconds)
	}
}

// waitNextSample pauses execution for a given number of seconds
// while periodically updating the watchdog.
func waitNextSample(seconds int) {
	wd := machine.Watchdog
	for i := 0; i < seconds; i++ {
		wd.Update()
		time.Sleep(time.Second)
	}
}
