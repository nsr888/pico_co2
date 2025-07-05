// This example demonstrates ENS160 usage with temperature and humidity compensation
// using an AHT20 sensor for improved accuracy.
//
// Wiring:
// ENS160+AHT21:
// - VCC to 3.3V, GND to ground
// - SDA to board SDA, SCL to board SCL
//
// Both sensors share the same I2C bus.

package main

import (
	"fmt"
	"log"
	"pico_co2/pkg/ens160"
	"time"

	"machine"

	"tinygo.org/x/drivers/aht20"
)

func main() {
	err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 200000,
	})
	if err != nil {
		log.Fatal("Failed to configure I2C:", err)
	}

	// Initialize AHT20
	aht21Sensor := aht20.New(machine.I2C0)
	aht21Sensor.Reset()
	aht21Sensor.Configure()

	// Initialize ENS160 with default address
	ens160Sensor := ens160.New(machine.I2C0, ens160.DefaultAddress)
	if err := ens160Sensor.Configure(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("ENS160 + AHT20 Compensation Example")
	fmt.Println("===================================")
	fmt.Println("Warming up... (this takes 3 minutes)")
	fmt.Println()

	for {
		// Read temperature and humidity from AHT20
		if err := aht21Sensor.Read(); err != nil {
			fmt.Printf("Error reading AHT20: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		temperature := aht21Sensor.Celsius()
		humidity := aht21Sensor.RelHumidity()

		// Set environmental compensation data for ENS160
		tempMilliC := int32(temperature * 1000)
		humidityMilliPct := int32(humidity * 1000)
		if err := ens160Sensor.SetEnvData(tempMilliC, humidityMilliPct); err != nil {
			fmt.Printf("Error setting environmental data: %v\n", err)
		}

		// Read air quality data
		err := ens160Sensor.Read(ens160.ReadConfig{
			WaitForNew:        true,
			WithValidityCheck: true,
		})
		if err != nil {
			fmt.Printf("Error reading ENS160: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Get readings
		co2 := ens160Sensor.LastCO2()
		tvoc := ens160Sensor.LastTVOC()
		aqi := ens160Sensor.LastAQI()

		// Display all readings
		fmt.Printf("Time: %s\n", time.Now().Format("15:04:05"))
		fmt.Printf("Temperature: %.1f°C\n", temperature)
		fmt.Printf("Humidity: %.1f%%\n", humidity)
		fmt.Printf("eCO2: %d ppm (%s)\n", co2, ens160.CO2String(co2))
		fmt.Printf("TVOC: %d ppb\n", tvoc)
		fmt.Printf("AQI: %d (%s)\n", aqi, ens160.AQIString(aqi))

		// Show detailed sensor status
		validityFlag, err := ens160Sensor.ReadValidityFlag()
		if err == nil {
			fmt.Printf("Validity: %s\n", ens160.ValidityFlagToString(validityFlag))
		}

		fmt.Println("---")

		// Wait before next reading
		time.Sleep(10 * time.Second)
	}
}
