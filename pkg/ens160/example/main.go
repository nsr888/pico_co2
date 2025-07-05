// This example demonstrates ENS160 usage.
//
// Wiring:
// - VCC to 3.3V, GND to ground
// - SDA to board SDA, SCL to board SCL

package main

import (
	"fmt"
	"log"
	"time"

	"machine"
	"tinygo.org/x/drivers"

	"pico_co2/pkg/ens160"
)

func main() {
	err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 200 * machine.KHz,
	})
	if err != nil {
		log.Fatal("Failed to configure I2C:", err)
	}

	dev := ens160.New(machine.I2C0, ens160.DefaultAddress)
	if err := dev.Configure(); err != nil {
		log.Fatal(err)
	}

	for {
		err := dev.Update(drivers.Concentration)
		if err != nil {
			fmt.Printf("Error reading ENS160: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf(
			"AQI=%d, TVOC=%dppb, eCO₂=%dppm\n",
			dev.AQI(),
			dev.TVOC(),
			dev.ECO2(),
		)

		time.Sleep(10 * time.Second)
	}
}
