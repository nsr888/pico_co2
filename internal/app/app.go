package app

import (
	"fmt"
	"log"
	"time"

	"machine"

	"pico_co2/internal/airquality"
	"pico_co2/internal/display"
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
	airqualitySensor airquality.AirQualitySensor
	fontDisplay      *display.FontDisplay
}

func New() (*App, error) {
	if err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFrequency,
		SDA:       i2cSDA,
		SCL:       i2cSCL,
	}); err != nil {
		return nil, err
	}

	fontDisplay, err := display.NewFontDisplay(machine.I2C0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize display: %w", err)
	}

	airQualitySensor := airquality.NewENS160AHT20Adapter(machine.I2C0)
	if err := airQualitySensor.Configure(); err != nil {
		return nil, fmt.Errorf("failed to configure air quality sensor: %w", err)
	}

	return &App{
		airqualitySensor: airQualitySensor,
		fontDisplay:      fontDisplay,
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
	log.Println("Starting loop")

	for {
		a.UpdateReadingsAndDisplay()
		waitNextSample(sampleTimeSeconds)
	}
}

func (a *App) UpdateReadingsAndDisplay() {
	readings, err := a.airqualitySensor.Read()
	if err != nil {
		log.Println("Error reading air quality sensor:", err)
		a.fontDisplay.DisplayError(err.Error())
		return
	}

	if readings.ValidityError != "" {
		log.Printf("Sensor readings is invalid: %+v\n", readings)
		a.fontDisplay.DisplayReadingsWithHI(readings)
		return
	}

	log.Printf("Sensor readings: %+v\n", readings)
	a.fontDisplay.DisplayReadingsWithHI(readings)
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
