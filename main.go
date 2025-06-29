package main

import (
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
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	playBoardLed(machine.LED, 3)
	log.Printf("Ready to go")

	// I2C setup
	err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFreq,
		SDA:       SDAPin,
		SCL:       SCLPin,
	})
	if err != nil {
		return err
	}
	bus := machine.I2C0
	log.Printf("I2C configuration: SDA=%v, SCL=%v, Frequency=%dHz", SDAPin, SCLPin, i2cFreq)

	// Dependency creation
	sensors, err := NewSensorReader(bus)
	if err != nil {
		return err
	}

	display, err := NewFontDisplay(bus)
	if err != nil {
		return err
	}

	// App creation and execution
	app := NewApp(machine.LED, sensors, display)
	app.Run()
	return nil
}

func playBoardLed(led machine.Pin, count int) {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 500)
		led.Low()
		time.Sleep(time.Millisecond * 500)
	}
}
