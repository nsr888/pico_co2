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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	application, err := app.New(i2cFrequency, i2cSDA, i2cSCL)
	if err != nil {
		log.Fatalf("application setup failed: %v", err)
	}

	application.Run()
}
