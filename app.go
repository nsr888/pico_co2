package main

import (
	"fmt"
	"log"
	"time"

	font "github.com/Nondzu/ssd1306_font"
	"machine"
	"tinygo.org/x/drivers/aht20"
	"tinygo.org/x/drivers/ds3231"
	"tinygo.org/x/drivers/ssd1306"

	"pico_co2/pkg/ens160"
)

// Application Logic
const (
	sampleTimeSeconds = 60
	watchDogMillis    = 8388 // max for RP2040 is 8388ms
)

// I2C Configuration
const (
	i2cFreq = 200000
	SDAPin  = machine.GP4
	SCLPin  = machine.GP5
)

// Display Configuration
const (
	displayWidth   = 128
	displayHeight  = 32
	displayAddress = ssd1306.Address_128_32
)

var ErrENS160ReadError = fmt.Errorf("ENS160 read error")

type App struct {
	i2c               *machine.I2C
	font              *font.Display
	fontDisplay       *FontDisplay
	display           *ssd1306.Device
	displayScreenNum  int
	led               machine.Pin
	ensCalibrated     bool
	ensStateSaved     bool
	lastValues        *Readings
	samplesUploaded   uint32
	startupCalTime    int64
	nextStateSaveTime int64
	aht20Sensor       *aht20.Device
	ens160Sensor      *ens160.Device
	ds3231Sensor      *ds3231.Device
}

func (a *App) ClearDisplay() {
	if a.display != nil {
		a.display.ClearBuffer()
		a.display.ClearDisplay()
	}
}

// NewApp initializes hardware and returns a new App instance.
func NewApp() (*App, error) {
	app := &App{
		led: machine.LED,
	}
	log.Printf("Setting up led")
	app.led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	const timeout = 3
	app.playBoardLed(timeout)
	log.Printf("Ready to go")

	if err := app.initI2C(); err != nil {
		return nil, err
	}

	if err := app.initDisplay(); err != nil {
		return nil, err
	}

	if err := app.initSensors(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initI2C() error {
	err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFreq,
		SDA:       SDAPin,
		SCL:       SCLPin,
	})
	if err != nil {
		return fmt.Errorf("failed to configure I2C: %w", err)
	}
	a.i2c = machine.I2C0
	log.Printf("I2C configuration: SDA=%v, SCL=%v, Frequency=%dHz", SDAPin, SCLPin, i2cFreq)
	return nil
}

func (a *App) initDisplay() error {
	display := ssd1306.NewI2C(a.i2c)
	display.Configure(ssd1306.Config{
		Width:   displayWidth,
		Height:  displayHeight,
		Address: displayAddress,
	})
	log.Printf("Display configured: Width=%d, Height=%d, Address=%d", displayWidth, displayHeight, displayAddress)

	a.display = &display
	a.ClearDisplay()

	fontLib := font.NewDisplay(display)
	a.fontDisplay = &FontDisplay{
		font:  &fontLib,
		clear: a.ClearDisplay,
	}
	return nil
}

func (a *App) initSensors() error {
	aht20Sensor := aht20.New(a.i2c)
	aht20Sensor.Reset()
	aht20Sensor.Configure()
	a.aht20Sensor = &aht20Sensor

	ens160Sensor := ens160.New(a.i2c, ens160.DefaultAddress)
	if err := ens160Sensor.Configure(); err != nil {
		return fmt.Errorf("failed to configure ENS160 sensor: %w", err)
	}
	a.ens160Sensor = ens160Sensor

	ds3231Sensor := ds3231.New(a.i2c)
	ds3231Sensor.Configure()
	a.ds3231Sensor = &ds3231Sensor

	return nil
}

func (a *App) playBoardLed(count int) {
	for i := 0; i < count; i++ {
		a.led.High()
		time.Sleep(time.Millisecond * 500)
		a.led.Low()
		time.Sleep(time.Millisecond * 500)
	}
}
