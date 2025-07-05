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
)

type App struct {
	led          machine.Pin
	sensorReader *service.SensorReader
}

type Config struct {
	I2cFrequency    uint32
	I2cSDA          machine.Pin
	I2cSCL          machine.Pin
	IsAdvancedSetup bool
}

func New(cfg Config) (*App, error) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	if err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: cfg.I2cFrequency,
		SDA:       cfg.I2cSDA,
		SCL:       cfg.I2cSCL,
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
		led:          led,
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

	a.led.Low()
	for {
		a.led.High()

		a.sensorReader.ProcessSensorReadings()

		time.Sleep(time.Millisecond * 200)
		a.led.Low()

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
