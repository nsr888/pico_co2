package main

import (
	"image/color"
	"log"

	"machine"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"

	"pico_co2/pkg/miniplot"
)

const (
	displayWidth      int16 = 128
	displayHeight     int16 = 32
	displayAddress          = ssd1306.Address_128_32
	sampleTimeSeconds       = 60
	watchDogMillis          = 8388 // max for RP2040 is 8388ms
	i2cFrequency            = 400 * machine.KHz
	i2cSDA                  = machine.GP4
	i2cSCL                  = machine.GP5
)

func main() {
	// Configure I2C
	if err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFrequency,
		SDA:       i2cSDA,
		SCL:       i2cSCL,
	}); err != nil {
		log.Fatalf("Failed to configure I2C: %v", err)
	}

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Width:   displayWidth,
		Height:  displayHeight,
		Address: displayAddress,
	})

	// Create fake CO2 measurements (realistic values 400-2000 ppm)
	fakeMeasurements := make([]int16, 128)
	for i := 0; i < 128; i++ {
		// Simulate varying CO2 levels
		base := 400 + (i * 10) // Start from 400 ppm
		if i > 100 {
			base = 1200 + (i-100)*5 // Higher readings
		}
		if i > 120 {
			base = 1800 + (i-120)*2 // Peak readings
		}
		fakeMeasurements[i] = int16(base)
	}

	// Clear display
	display.ClearDisplay()

	// Create plot
	font := tinyfont.Fonter(&proggy.TinySZ8pt7b)
	plot := miniplot.NewMiniPlot(&display, font, 128, 32)

	// Draw plot at position (0, 0)
	plot.DrawLineChart(fakeMeasurements)

	// Display title using proggy font
	tinyfont.WriteLine(&display, font, 0, 8, "CO2 Plot", color.RGBA{255, 255, 255, 255})

	// Update display
	display.Display()

	// Keep running to see the display
	for {
	}
}
