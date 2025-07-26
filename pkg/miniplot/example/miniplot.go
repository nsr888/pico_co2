package main

import (
	"image/color"
	"log"
	"time"

	"machine"
	"tinygo.org/x/drivers/ssd1306"
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
	countMeasurements := 128 // Number of measurements to simulate
	fakeMeasurements := make([]int16, countMeasurements)
	for i := 0; i < len(fakeMeasurements); i++ {
		// Simulate varying CO2 levels
		base := 400 + (i * 10) // Start from 400 ppm
		if i > 80 {
			base = 1200 + (i-100)*5 // Higher readings
		}
		if i > 100 {
			base = 1800 + (i-120)*2 // Peak readings
		}
		fakeMeasurements[i] = int16(base)
	}

	time.Sleep(2 * time.Second)

	// Create plot
	font := &proggy.TinySZ8pt7b
	plot, err := miniplot.NewMiniPlot(&display, font, 128, 32, color.RGBA{255, 255, 255, 255})
	if err != nil {
		log.Fatalf("Failed to create MiniPlot: %v", err)
	}

	display.ClearDisplay()

	// Draw plot at position (0, 0)
	plot.DrawLineChart(fakeMeasurements, "fake")

	// Keep running to see the display
	for {
	}
}
