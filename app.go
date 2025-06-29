package main

import (
	"errors"
	"log"
	"machine"
	"time"
)

// Application Logic
const (
	sampleTimeSeconds = 60
	watchDogMillis    = 8388 // max for RP2040 is 8388ms
)

var ErrAirQualityReadError = errors.New("air quality sensor read error")

type App struct {
	led     machine.Pin
	sensors *SensorReader
	display *FontDisplay
}

// NewApp creates a new App instance with its dependencies.
func NewApp(led machine.Pin, sensors *SensorReader, display *FontDisplay) *App {
	return &App{
		led:     led,
		sensors: sensors,
		display: display,
	}
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

		a.processSensorReadings()

		time.Sleep(time.Millisecond * 200)
		a.led.Low()

		waitNextSample(sampleTimeSeconds)
	}
}

func (a *App) processSensorReadings() {
	readings, err := a.sensors.Read()
	logger := log.New(log.Writer(), readings.Timestamp.Format(time.RFC3339)+" ", 0)
	logger.Printf("Readings: %+v", readings)
	switch {
	case err != nil && !errors.Is(err, ErrAirQualityReadError):
		logger.Panicf("Error reading sensors: %v", err)
	case errors.Is(err, ErrAirQualityReadError):
		logger.Println(err)
		a.display.DisplayBasic(readings)
	case readings.CO2 == 0 && readings.Temperature != 0:
		logger.Println("CO2 readings are zero, displaying temperature data only")
		a.display.DisplayBasic(readings)
	default:
		a.display.DisplayFull(readings)
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
