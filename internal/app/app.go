package app

import (
	"cmp"
	"fmt"
	"machine"
	"pico_co2/internal/button"
	"pico_co2/internal/display"
	"pico_co2/internal/types"
	"pico_co2/pkg/ens160"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/aht20"
	"tinygo.org/x/drivers/scd4x"
	"tinygo.org/x/drivers/ssd1306"
)

type Config struct {
	Display struct {
		Width   int16
		Height  int16
		Address uint16
	}
	I2C struct {
		Frequency uint32
		SDA       machine.Pin
		SCL       machine.Pin
	}
	Buttons struct {
		Button1 machine.Pin
		Button2 machine.Pin
	}
	Timeouts struct {
		Startup         time.Duration
		Minute          time.Duration
		MinimalInterval time.Duration
	}
	QueueCapacity       int
	DefaultDisplayIndex int
}

func DefaultConfig() Config {
	cfg := Config{}
	cfg.Display.Width = 128
	cfg.Display.Height = 32
	cfg.Display.Address = ssd1306.Address_128_32
	cfg.I2C.Frequency = 400 * machine.KHz
	cfg.I2C.SDA = machine.GP4
	cfg.I2C.SCL = machine.GP5
	cfg.Buttons.Button1 = machine.GP10
	cfg.Buttons.Button2 = machine.GP11
	cfg.Timeouts.Startup = 1 * time.Minute
	cfg.Timeouts.Minute = 60 * time.Second
	cfg.Timeouts.MinimalInterval = 1 * time.Second
	cfg.QueueCapacity = 128
	cfg.DefaultDisplayIndex = 2
	return cfg
}

func (c Config) initI2C() error {
	return machine.I2C0.Configure(machine.I2CConfig{
		Frequency: c.I2C.Frequency,
		SDA:       c.I2C.SDA,
		SCL:       c.I2C.SCL,
	})
}

func (c Config) initDisplay(bus drivers.I2C) (display.Renderer, error) {
	disp := ssd1306.NewI2C(bus)
	disp.Configure(ssd1306.Config{
		Width:   c.Display.Width,
		Height:  c.Display.Height,
		Address: c.Display.Address,
	})
	// REDUCE BRIGHTNESS
	// reduce contrast for night viewing
	disp.Command(ssd1306.SETCONTRAST)
	disp.Command(0x01)
	// precharge period
	disp.Command(ssd1306.SETPRECHARGE)
	disp.Command(
		0xE1,
	) // 0xF1 default, 0xE1 for lower power, 0xD2 for even lower
	// VCOMH deselect level
	disp.Command(ssd1306.SETVCOMDETECT)
	disp.Command(
		0x30,
	) // 0x20 default, 0x30 for lower power, 0x40 for even lower
	return display.NewSSD1306Adapter(&disp), nil
}

type RawReadings struct {
	CO2         uint16
	Temperature float32
	Humidity    float32
}

type Sensors struct {
	aht20  *aht20.Device
	ens160 *ens160.Device
	scd4x  *scd4x.Device
}

func NewSensors(bus drivers.I2C) (*Sensors, error) {
	s := &Sensors{}

	if err := s.initAHT20(bus); err != nil {
		return nil, fmt.Errorf("aht20 init: %w", err)
	}

	if err := s.initENS160(bus); err != nil {
		return nil, fmt.Errorf("ens160 init: %w", err)
	}

	if err := s.initSCD4x(bus); err != nil {
		return nil, fmt.Errorf("scd4x init: %w", err)
	}

	return s, nil
}

func (s *Sensors) initAHT20(bus drivers.I2C) error {
	aht20Sensor := aht20.New(bus)
	aht20Sensor.Reset()
	aht20Sensor.Configure()
	s.aht20 = &aht20Sensor
	return nil
}

func (s *Sensors) initENS160(bus drivers.I2C) error {
	ens160Sensor := ens160.New(bus, ens160.DefaultAddress)
	if err := ens160Sensor.Sleep(); err != nil {
		return err
	}
	s.ens160 = ens160Sensor
	return nil
}

func (s *Sensors) initSCD4x(bus drivers.I2C) error {
	scd4xSensor := scd4x.New(bus)
	time.Sleep(1500 * time.Millisecond)
	if err := scd4xSensor.Configure(); err != nil {
		return err
	}

	time.Sleep(1500 * time.Millisecond)

	if err := scd4xSensor.StartPeriodicMeasurement(); err != nil {
		return err
	}

	time.Sleep(1500 * time.Millisecond)

	s.scd4x = scd4xSensor
	return nil
}

func (s *Sensors) Read() (*RawReadings, error) {
	if err := s.aht20.Read(); err != nil {
		return nil, fmt.Errorf("aht20 read: %w", err)
	}

	co2, err := s.scd4x.ReadCO2()
	if err != nil {
		return nil, fmt.Errorf("scd4x read: %w", err)
	}

	return &RawReadings{
		CO2:         uint16(co2),
		Temperature: s.aht20.Celsius(),
		Humidity:    s.aht20.RelHumidity(),
	}, nil
}

type DisplayManager struct {
	renderer     display.Renderer
	currentIndex int
}

func NewDisplayManager(
	renderer display.Renderer,
	currentIndex int,
) *DisplayManager {
	return &DisplayManager{
		renderer:     renderer,
		currentIndex: currentIndex,
	}
}

func (dm *DisplayManager) NextDisplay() {
	dm.currentIndex = (dm.currentIndex + 1) % len(display.MethodRegistry)
}

func (dm *DisplayManager) Render(readings *types.Readings) {
	if readings.Error != "" {
		display.RenderError(dm.renderer, readings)
		return
	}

	if dm.currentIndex >= len(display.MethodRegistry) {
		dm.currentIndex = 0
	}

	renderMethod := display.MethodRegistry[dm.currentIndex]
	renderMethod.Fn(dm.renderer, readings)
}

type App struct {
	config         Config
	sensors        *Sensors
	displayManager *DisplayManager
	button1        *button.TouchButton
	button2        *button.TouchButton
}

func New(cfg Config) (*App, error) {
	err := cfg.initI2C()
	if err != nil {
		return nil, fmt.Errorf("i2c init: %w", err)
	}

	renderer, err := cfg.initDisplay(machine.I2C0)
	if err != nil {
		return nil, fmt.Errorf("display init: %w", err)
	}

	sensors, err := NewSensors(machine.I2C0)
	if err != nil {
		return nil, fmt.Errorf("sensors init: %w", err)
	}

	return &App{
		config:         cfg,
		sensors:        sensors,
		displayManager: NewDisplayManager(renderer, cfg.DefaultDisplayIndex),
		button1:        button.NewTouchButton(cfg.Buttons.Button1),
		button2:        button.NewTouchButton(cfg.Buttons.Button2),
	}, nil
}

func (a *App) Run() {
	readings := types.InitReadings(a.config.QueueCapacity)

	wd := machine.Watchdog
	wd.Configure(machine.WatchdogConfig{
		TimeoutMillis: machine.WatchdogMaxTimeout,
	})
	wd.Start()

	println("starting loop")

	for {
		wd.Update()

		a.handleInput(readings)
		a.updateReadings(readings)
		a.render(readings)

		time.Sleep(50 * time.Millisecond)
	}
}

func (a *App) handleInput(readings *types.Readings) {
	if a.button1.Consume() {
		a.displayManager.NextDisplay()
		readings.IsDrawen = false
	}
}

func (a *App) shouldAlwaysUpdateDisplay() bool {
	// Displays that should always update regardless of measurement changes
	alwaysUpdateDisplays := []string{
		"RenderCO2Graph",
		"RenderSparkline",
	}

	if a.displayManager.currentIndex >= len(display.MethodRegistry) {
		return false
	}

	currentDisplayName := display.MethodRegistry[a.displayManager.currentIndex].Name
	for _, name := range alwaysUpdateDisplays {
		if currentDisplayName == name {
			return true
		}
	}
	return false
}

func (a *App) updateReadings(readings *types.Readings) {
	shouldUpdate := cmp.Or(
		readings.LastUpdateAt.IsZero(),
		a.button2.Consume(),
		time.Since(readings.LastUpdateAt) >= a.config.Timeouts.MinimalInterval &&
			time.Since(readings.FirstReadingAt) < a.config.Timeouts.Startup,
		time.Since(readings.LastUpdateAt) >= a.config.Timeouts.Minute,
	)

	if shouldUpdate {
		if raw, err := a.sensors.Read(); err != nil {
			readings.Error = err.Error()
			readings.IsDrawen = false
		} else {
			readings.AddReadings(raw.CO2, raw.Temperature, raw.Humidity)
			fmt.Printf("time: %s, CO2: %d ppm, T: %.2f °C, H: %.2f %%, co2 len: %d, temp len: %d, hum len: %d\n",
				time.Now().Format("15:04:05"),
				raw.CO2, raw.Temperature, raw.Humidity,
				readings.History.CO2.Len(),
				readings.History.Temperature.Len(),
				readings.History.Humidity.Len(),
			)
			readings.Error = ""

			// Mark as needing redraw if measurements changed or display should always update
			if readings.MeasurementsChanged() || a.shouldAlwaysUpdateDisplay() {
				readings.IsDrawen = false
			}
		}
	}
}

func (a *App) render(readings *types.Readings) {
	if !readings.IsDrawen {
		a.displayManager.Render(readings)
		readings.IsDrawen = true
	}
}
