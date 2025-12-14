package app

import (
	"errors"
	"fmt"
	"machine"
	"pico_co2/internal/button"
	"pico_co2/internal/display"
	"pico_co2/internal/types"
	"pico_co2/pkg/ens160"
	"time"

	"tinygo.org/x/drivers/aht20"
	"tinygo.org/x/drivers/scd4x"
	"tinygo.org/x/drivers/ssd1306"
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
	button1Pin     = machine.GP10
	button2Pin     = machine.GP11
)

const (
	displayWidth   int16 = 128
	displayHeight  int16 = 32
	displayAddress       = ssd1306.Address_128_32
	queueCapacity        = 128 // Number of readings to keep in memory
)

type App struct {
	renderer            display.Renderer
	aht20Sensor         *aht20.Device
	ens160Sensor        *ens160.Device
	scd4xSensor         *scd4x.Device
	currentDisplayIndex int
	button1             *button.TouchButton
	button2             *button.TouchButton
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
	// reduce contrast for night time viewing
	ssd1306disp.Command(ssd1306.SETCONTRAST)
	ssd1306disp.Command(0x01)

	renderer := display.NewSSD1306Adapter(&ssd1306disp)

	aht20Sensor := aht20.New(machine.I2C0)
	aht20Sensor.Reset()
	aht20Sensor.Configure()
	ens160Sensor := ens160.New(machine.I2C0, ens160.DefaultAddress)
	if err := ens160Sensor.Sleep(); err != nil {
		return nil, errors.New("failed to put ENS160 to sleep: " + err.Error())
	}

	scd4xSensor := scd4x.New(machine.I2C0)
	time.Sleep(1500 * time.Millisecond)
	if err := scd4xSensor.Configure(); err != nil {
		return nil, errors.New("failed to configure SCD4x: " + err.Error())
	}

	time.Sleep(1500 * time.Millisecond)

	if err := scd4xSensor.StartPeriodicMeasurement(); err != nil {
		return nil, errors.New(
			"failed to start SCD4x periodic measurement: " + err.Error(),
		)
	}

	time.Sleep(1500 * time.Millisecond)

	button1 := button.NewTouchButton(button1Pin)
	button2 := button.NewTouchButton(button2Pin)

	return &App{
		renderer:     renderer,
		aht20Sensor:  &aht20Sensor,
		ens160Sensor: ens160Sensor,
		scd4xSensor:  scd4xSensor,
		button1:      button1,
		button2:      button2,
		// Display method 0 : co2_bar_with_nums
		// Display method 1 : bars
		// Display method 2 : basic
		// Display method 3 : temp_humid
		// Display method 4 : bars_with_large_nums
		// Display method 5 : co2_graph
		// Display method 6 : render_heat_index_status
		// Display method 7 : large_bar
		// Display method 8 : nums
		// Display method 9 : sparkline
		currentDisplayIndex: 3,
	}, nil
}

// Run starts the main application loop.
func (a *App) Run() {
	r := types.InitReadings(queueCapacity)

	wd := machine.Watchdog
	wd.Configure(machine.WatchdogConfig{
		TimeoutMillis: watchDogMillis,
	})
	wd.Start()
	println("starting loop")

	displayMethods := display.GetAllDisplayMethods()

	for i, method := range displayMethods {
		println("Display method", i, ":", method)
	}

	lastSampleTime := time.Now()
	displayDirty := true // Render on the first loop iteration

	readAndRecord := func() {
		lastSampleTime = time.Now()
		errAht := a.aht20Sensor.Read()
		if errAht != nil {
			r.Error = "AHT20 read error: " + errAht.Error()
			displayDirty = true
			return
		}
		temp := a.aht20Sensor.Celsius()
		hum := a.aht20Sensor.RelHumidity()
		co2, err := a.scd4xSensor.ReadCO2()
		if err != nil {
			r.Error = "SCD4x read error: " + err.Error()
			displayDirty = true
			return
		}
		fmt.Printf(
			"CO2: %d ppm, Temp: %.2f °C, Humidity: %.2f %%\n",
			co2,
			temp,
			hum,
		)
		r.AddReadings(uint16(co2), temp, hum)
		r.Error = ""
		displayDirty = true // Flag that the display needs to be updated.
	}

	// Initial sensor read.
	readAndRecord()

	for {
		wd.Update()

		// Button 1 press changes the display method and flags a re-render.
		if a.button1.Consume() {
			a.currentDisplayIndex = (a.currentDisplayIndex + 1) % len(
				displayMethods,
			)
			displayDirty = true
		}

		button2Pressed := a.button2.Consume()

		// New sensor data is read on a schedule or when button 2 is pressed.
		if time.Since(lastSampleTime) >= minuteTimeout || button2Pressed {
			readAndRecord()
		}

		// Re-render the display only if something has changed.
		if displayDirty {
			if r.Error != "" {
				display.RenderError(a.renderer, r)
			} else {
				currentMethod := displayMethods[a.currentDisplayIndex]
				if renderFunc, exists := display.MethodRegistry[currentMethod]; exists {
					renderFunc(a.renderer, r)
				} else {
					display.RenderBasic(a.renderer, r)
				}
			}
			displayDirty = false // Reset the flag until the next change.
		}

		time.Sleep(50 * time.Millisecond)
	}
}
