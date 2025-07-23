package main

import (
	"log"
	"machine"
	"ssd1306"
	"color"
	"miniplot"
	"tinyfont"
	"proggy"
)

func main() {
	// Configure I2C
	if err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GP4,
		SCL:       machine.GP5,
	}); err != nil {
		log.Fatal("Failed to configure I2C:", err)
	}

	// Initialize SSD1306 display
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Width:   128,
		Height:  32,
		Address: 0x3C,
	})

	// Create fake CO2 measurements (realistic values 400-2000 ppm)
	fakeMeasurements := make([]int16, 128)
	for i := 0; i < 128; i++ {
		// Simulate varying CO2 levels
		base := 400 + (i * 10) // Start from 400 ppm
		if i > 100 {
			base = 1200 + (i - 100) * 5 // Higher readings
		}
		if i > 120 {
			base = 1800 + (i - 120) * 2 // Peak readings
		}
		fakeMeasurements[i] = int16(base)
	}

	// Clear display
	display.Clear(color.Black)
	
	// Create plot
	plot := miniplot.NewPlot(display, color.RGBA{255, 255, 255, 255}, 128, 32)
	
	// Draw plot at position (0, 0)
	plot.Draw(fakeMeasurements, 0, 0)
	
	// Display title using proggy font
	font := tinyfont.Fonter(&proggy.TinySZ8pt7b)
	tinyfont.WriteLine(display, font, 0, 8, "CO2 Plot", color.RGBA{255, 255, 255, 255})
	
	// Update display
	display.Update()
	
	// Keep running to see the display
	for {
	}
}
