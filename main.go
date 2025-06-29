package main

import (
	"errors"
	"log"
	"machine"
	"time"
)

// I2C Configuration
const (
	i2cFreq = 200000
	SDAPin  = machine.GP4
	SCLPin  = machine.GP5
)

func main() {
	app, err := NewApp()
	if err != nil {
		log.Printf("Failed to initialize hardware:", err.Error())
		return
	}

	wd := machine.Watchdog
	config := machine.WatchdogConfig{
		TimeoutMillis: watchDogMillis,
	}
	wd.Configure(config)
	wd.Start()
	log.Printf("starting loop")

	app.led.Low()
	for {
		app.led.High()

		readings, err := app.readSensors()
		log := log.New(log.Writer(), readings.Timestamp.Format(time.RFC3339)+" ", 0)
		log.Printf("Readings: %+v", readings)
		switch {
		case err != nil && !errors.Is(err, ErrENS160ReadError):
			log.Panicf("Error reading sensors: %v", err)
		case errors.Is(err, ErrENS160ReadError):
			log.Println(err)
			app.fontDisplay.DisplayAHT20Readings(readings)
		case readings.AQI == 0 && readings.ECO2 == 0 && readings.TVOC == 0:
			log.Println("ENS160 readings are zero, displaying AHT20 data only")
			app.fontDisplay.DisplayAHT20Readings(readings)
		default:
			app.fontDisplay.DisplayFullReadings(readings)
		}

		time.Sleep(time.Millisecond * 200)
		app.led.Low()

		for i := 0; i < sampleTimeSeconds; i++ {
			wd.Update()
			time.Sleep(time.Second)
		}
	}
}
