package main

import (
	"log"

	"machine"

	"pico_co2/internal/app"
)

const (
	i2cFrequency = 200 * machine.KHz
	i2cSDA       = machine.GP4
	i2cSCL       = machine.GP5
)

func main() {
	log.SetFlags(log.Lshortfile)

	config := app.Config{
		I2cFrequency:    i2cFrequency,
		I2cSDA:          i2cSDA,
		I2cSCL:          i2cSCL,
		IsAdvancedSetup: false,
	}

	application, err := app.New(config)
	if err != nil {
		log.Fatalf("application setup failed: %v", err)
	}

	application.Run()
}
