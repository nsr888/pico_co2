package app

import (
	"errors"
	"fmt"
	"time"

	"machine"
	"tinygo.org/x/drivers/ssd1306"

	"pico_co2/internal/airquality"
	"pico_co2/internal/display"
	"pico_co2/internal/types"
)

// Application Logic
const (
	startupTimeout = 3 * time.Minute
	minuteTimeout  = 60 * time.Second
	shortTimeout   = 5 * time.Second
	watchDogMillis = machine.WatchdogMaxTimeout
	i2cFrequency   = 400 * machine.KHz
	i2cSDA         = machine.GP4
	i2cSCL         = machine.GP5
)

const (
	displayWidth   int16 = 128
	displayHeight  int16 = 32
	displayAddress       = ssd1306.Address_128_32
	queueCapacity        = 128 // Number of readings to keep in memory
)

type App struct {
	renderer         display.Renderer
	airqualitySensor airquality.AirQualitySensor
}

func New() (*App, error) {
	if err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFrequency,
		SDA:       i2cSDA,
		SCL:       i2cSCL,
	}); err != nil {
		return nil, err
	}

	ssd1306disp := ssd1306.NewI2C(machine.I2C0)
	ssd1306disp.Configure(ssd1306.Config{
		Width:   displayWidth,
		Height:  displayHeight,
		Address: displayAddress,
	})

	renderer := display.NewSSD1306Adapter(&ssd1306disp)

	airQualitySensor := airquality.NewENS160AHT20Adapter(machine.I2C0)
	if err := airQualitySensor.Configure(); err != nil {
		return nil, errors.New("failed to configure air quality sensor: " + err.Error())
	}

	return &App{
		renderer:         renderer,
		airqualitySensor: airQualitySensor,
	}, nil
}

// Run starts the main application loop.
func (a *App) Run() {
	r := types.InitReadings(queueCapacity)

	wd := machine.Watchdog
	config := machine.WatchdogConfig{
		TimeoutMillis: watchDogMillis,
	}
	wd.Configure(config)
	wd.Start()
	println("starting loop")

	var isDisplayGraph bool
	for {
		readings, err := a.airqualitySensor.Read()
		if err != nil {
			println("error reading air quality sensor:", err)
			readings := types.Readings{
				Error: err.Error(),
			}
			display.RenderError(a.renderer, &readings)

			return
		}

		fmt.Printf("Readings: %+v\n", readings)

		r.AddReadings(
			readings.CO2,
			readings.TVOC,
			readings.AQI,
			readings.Temperature,
			readings.Humidity,
			readings.DataValidityWarning,
		)

		// on startup, wait shorter to get updates quickly
		if time.Since(r.FirstReadingTime) < startupTimeout {
			display.RenderTempHumid(a.renderer, r)
			waitNextSample(shortTimeout)
		} else {
			switch isDisplayGraph {
			case false:
				display.RenderAqiBarWithNums(a.renderer, r)
			case true:
				display.RenderCO2Graph(a.renderer, r)
			}
			isDisplayGraph = !isDisplayGraph
			waitNextSample(minuteTimeout)
		}
	}
}

// waitNextSample pauses execution for a given number of seconds
// while periodically updating the watchdog.
func waitNextSample(timeout time.Duration) {
	seconds := int(timeout.Seconds())
	wd := machine.Watchdog
	for i := 0; i < seconds; i++ {
		wd.Update()
		time.Sleep(time.Second)
	}
}
