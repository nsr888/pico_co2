package app

import (
	"fmt"
	"log"
	"machine"
	"pico_co2/internal/airquality"
	"pico_co2/internal/display"
	"pico_co2/internal/service"
	"time"

	"tinygo.org/x/drivers/ds3231"
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

func New(i2cFrequency uint32, sda machine.Pin, scl machine.Pin) (*App, error) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	if err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFrequency,
		SDA:       sda,
		SCL:       scl,
	}); err != nil {
		return nil, err
	}

	d, err := display.NewFontDisplay(machine.I2C0)
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

	sensorReader := service.NewSensorReader(airQualitySensor, &ds3231Sensor, d)

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
