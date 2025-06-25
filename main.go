package main

import (
	"errors"
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

const (
	sampleTimeSeconds = 60
	watchDogMillis    = 8388 // max for RP2040 is 8388ms

	i2cFreq = 200000
	SDAPin  = machine.GP4
	SCLPin  = machine.GP5

	displayWidth   = 128
	displayHeight  = 32
	displayAddress = ssd1306.Address_128_32
)

var ErrENS160ReadError = fmt.Errorf("ENS160 read error")

type SensorDevice struct {
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

func (d *SensorDevice) ClearDisplay() {
	if d.display != nil {
		d.display.ClearBuffer()
		d.display.ClearDisplay()
	}
}

type FontDisplay struct {
	font  *font.Display
	clear func()
}

func (f *FontDisplay) DisplayAHT20Readings(r Readings) {
	if f == nil {
		return
	}
	f.clear()
	f.font.Configure(font.Config{FontType: font.FONT_16x26})
	tempStr := fmt.Sprintf("%.0f", r.Temperature)
	f.font.XPos = int16((128 - (len(tempStr) * 16)) / 2)
	f.font.YPos = 0
	f.font.PrintText(tempStr)

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_6x8})
	formatString := fmt.Sprintf("Temp %.1fC Hum %.1f%%", r.Temperature, r.Humidity)
	f.font.XPos = int16((128 - (len(formatString) * 6)) / 2)
	f.font.YPos = 24
	f.font.PrintText(formatString)
}

func (f *FontDisplay) DisplayFullReadings(r Readings) {
	if f == nil {
		return
	}
	f.clear()

	// Big numbers for eCO2 and AQI
	f.font.Configure(font.Config{FontType: font.FONT_16x26})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.ECO2))
	tempStr := fmt.Sprintf("%.0f", r.Temperature)
	f.font.XPos = int16(128 - (len(tempStr) * 16))
	f.font.YPos = 0
	f.font.PrintText(tempStr)

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_6x8})
	co2Str := "eCO2"
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(co2Str)
	tempTitleStr := "Temp"
	f.font.XPos = int16(128 - (len(tempTitleStr) * 6))
	f.font.YPos = 24
	f.font.PrintText(tempTitleStr)
	f.font.XPos = int16(128-(len(r.Status)*6)) / 2
	f.font.YPos = 24
	f.font.PrintText(r.Status)
}

func (f *FontDisplay) DisplayFullReadingsCO2andAQI(r Readings) {
	if f == nil {
		return
	}
	f.clear()

	// Big numbers for eCO2 and AQI
	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	f.font.XPos = 30
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.ECO2))
	f.font.XPos = 110
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.AQI))

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_6x8})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText("eCO2")
	f.font.XPos = 87
	f.font.YPos = 0
	f.font.PrintText("AQI")
	f.font.XPos = 0
	f.font.YPos = 16
	f.font.PrintText("-----------------------")
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(fmt.Sprintf("T %.0f H %.0f", r.Temperature, r.Humidity))
	f.font.XPos = int16(128 - (len(r.Status) * 6))
	f.font.YPos = 24
	f.font.PrintText(r.Status)
}

// Readings represents sensor data
type Readings struct {
	AQI         uint8     `json:"aqi"`
	ECO2        uint16    `json:"eco2"`
	TVOC        uint16    `json:"tvoc"`
	Humidity    float32   `json:"humidity"`
	Temperature float32   `json:"temperature"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

// Initialize hardware and return device instance
func initHardware() (*SensorDevice, error) {
	d := &SensorDevice{
		led: machine.LED,
	}
	log.Printf("Setting up led")
	d.led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	const timeout = 3
	d.playBoardLed(timeout)
	log.Printf("Ready to go")

	err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: i2cFreq,
		SDA:       SDAPin,
		SCL:       SCLPin,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to configure I2C: %w", err)
	}
	d.i2c = machine.I2C0
	log.Printf("I2C configuration: SDA=%v, SCL=%v, Frequency=%dHz", SDAPin, SCLPin, i2cFreq)

	display := ssd1306.NewI2C(d.i2c)
	display.Configure(ssd1306.Config{
		Width:   displayWidth,
		Height:  displayHeight,
		Address: displayAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to configure display: %w", err)
	}
	log.Printf("Display configured: Width=%d, Height=%d, Address=%d", displayWidth, displayHeight, displayAddress)

	d.display = &display
	d.ClearDisplay()

	fontLib := font.NewDisplay(display)
	d.fontDisplay = &FontDisplay{
		font:  &fontLib,
		clear: d.ClearDisplay,
	}

	aht20Sensor := aht20.New(d.i2c)
	aht20Sensor.Reset()
	aht20Sensor.Configure()
	d.aht20Sensor = &aht20Sensor

	ens160Sensor := ens160.New(d.i2c, ens160.DefaultAddress)
	if err := ens160Sensor.Reset(); err != nil {
		panic(err)
	}
	d.ens160Sensor = ens160Sensor

	ds3231Sensor := ds3231.New(d.i2c)
	ds3231Sensor.Configure()
	d.ds3231Sensor = &ds3231Sensor

	return d, nil
}

func (d *SensorDevice) playBoardLed(count int) {
	for i := 0; i < count; i++ {
		d.led.High()
		time.Sleep(time.Millisecond * 500)
		d.led.Low()
		time.Sleep(time.Millisecond * 500)
	}
}

func (d *SensorDevice) readSensors() (Readings, error) {
	var r Readings

	dt, err := d.ds3231Sensor.ReadTime()
	if err != nil {
		log.Printf("Error reading time: %v", err)
	}

	r.Timestamp = dt

	if d.aht20Sensor == nil {
		return r, fmt.Errorf("AHT20 sensor not initialized")
	}

	if err := d.aht20Sensor.Read(); err != nil {
		return r, fmt.Errorf("failed to read AHT20 sensor: %w", err)
	}

	r.Temperature = d.aht20Sensor.Celsius()
	r.Humidity = d.aht20Sensor.RelHumidity()

	if err := d.ens160Sensor.SetEnvData(r.Temperature, r.Humidity); err != nil {
		return r, fmt.Errorf("failed to set environment data for ENS160: %w", err)
	}

	// err = d.ens160Sensor.Read(ens160.WithValidityCheck(), ens160.WithWaitForNew())
	err = d.ens160Sensor.Read()
	if err != nil {
		return r, fmt.Errorf("%w: %v", ErrENS160ReadError, err)
	}

	r.AQI = d.ens160Sensor.LastAQI()
	r.ECO2 = d.ens160Sensor.LastCO2()
	r.TVOC = d.ens160Sensor.LastTVOC()
	r.Status = ens160.CO2String(d.ens160Sensor.LastCO2())

	return r, nil
}

func main() {
	device, err := initHardware()
	if err != nil {
		log.Printf("Failed to initialize hardware:", err.Error())
		return
	}

	wd := machine.Watchdog
	config := machine.WatchdogConfig{
		TimeoutMillis: watchDogMillis,
	}
	wd.Configure(config)
	wd.Start()
	log.Printf("starting loop")

	device.led.Low()
	for {
		device.led.High()

		readings, err := device.readSensors()
		log := log.New(log.Writer(), readings.Timestamp.Format(time.RFC3339)+" ", 0)
		log.Printf("Readings: %+v", readings)
		switch {
		case err != nil && !errors.Is(err, ErrENS160ReadError):
			log.Panicf("Error reading sensors: %v", err)
		case errors.Is(err, ErrENS160ReadError):
			log.Println(err)
			device.fontDisplay.DisplayAHT20Readings(readings)
		case readings.AQI == 0 && readings.ECO2 == 0 && readings.TVOC == 0:
			log.Println("ENS160 readings are zero, displaying AHT20 data only")
			device.fontDisplay.DisplayAHT20Readings(readings)
		default:
			device.fontDisplay.DisplayFullReadings(readings)
		}

		time.Sleep(time.Millisecond * 200)
		device.led.Low()

		for i := 0; i < sampleTimeSeconds; i++ {
			wd.Update()
			time.Sleep(time.Second)
		}
	}
}
