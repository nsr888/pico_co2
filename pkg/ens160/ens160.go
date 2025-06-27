/*
Package ens160 provides a driver for the ENS160 Digital Metal-Oxide Multi-Gas
Sensor manufactured by ScioSense.

Example of usage:

	device := ens160.New(machine.I2C1, ens160.DefaultAddress)
	if err := device.Configure(); err != nil {
		log.Fatal(err)
	}

	for {
		if err := device.Read(ens160.WithWaitForNew(), ens160.WithValidityCheck()); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("eCO2: %d, TVOC: %d, AQI: %d\n", device.LastCO2(), device.LastTVOC(), device.LastAQI())
		time.Sleep(5 * time.Second)
	}
*/
package ens160

import (
	"errors"
	"fmt"
	"time"

	"machine"
)

// ErrInitialStartUpPhase indicates the device is in its initial startup phase.
// The device requires 1 hour of continuous operation after first power-on.
// After 24 hours of continuous operation, this status is stored in non-volatile memory.
// If unpowered before 24 hours, the device will resume initial startup mode after re-powering.
var ErrInitialStartUpPhase = errors.New("initial startup required")

// ErrWarmUpPhase indicates readings are unavailable during the 3-minute warm-up period.
var ErrWarmUpPhase = errors.New("warmup in progress")

// ErrNoValidOutput indicates sensor signals are out of range or invalid.
var ErrNoValidOutput = errors.New("no valid output")

// Device wraps an I2C connection to an ENS160 device.
type Device struct {
	bus      *machine.I2C
	address  uint8
	lastCO2  uint16
	lastTVOC uint16
	lastAQI  uint8
}

// New creates a new ENS160 device with the given I2C bus and address.
func New(bus *machine.I2C, address uint8) *Device {
	if address == 0 {
		address = DefaultAddress
	}
	return &Device{
		bus:     bus,
		address: address,
	}
}

// Configure sets up the sensor by resetting it and putting it into standard mode.
// Should be called once after New.
func (d *Device) Configure() error {
	if err := d.Reset(); err != nil {
		return err
	}
	return d.SetOperatingMode(ModeStandard)
}

// GetOperatingMode reads the current operating mode of the device.
func (d *Device) GetOperatingMode() (uint8, error) {
	data := []uint8{0}
	err := d.bus.ReadRegister(d.address, regOperatingMode, data)

	return data[0], err
}

// SetOperatingMode sets the device's operating mode.
func (d *Device) SetOperatingMode(mode uint8) error {
	if mode != ModeDeepSleep &&
		mode != ModeIdle &&
		mode != ModeStandard &&
		mode != ModeReset {
		return errors.New("invalid operating mode")
	}

	return d.bus.WriteRegister(d.address, regOperatingMode, []uint8{mode})
}

// GetRawCO2 reads the calculated equivalent CO2 concentration in PPM.
func (d *Device) GetRawCO2() (uint16, error) {
	data := []uint8{0, 0}
	err := d.bus.ReadRegister(d.address, regECO2, data)
	if err != nil {
		return 0, err
	}

	return uint16(data[0]) | uint16(data[1])<<8, nil
}

func CO2String(value uint16) string {
	switch {
	case value < 400:
		return "No data"
	case value < 600:
		return "Excellent"
	case value < 800:
		return "Good"
	case value < 1000:
		return "Fair"
	case value < 1500:
		return "Poor"
	default:
		return "Bad"
	}
}

// GetRawTVOC reads the calculated Total Volatile Organic Compounds
// concentration in ppb.
func (d *Device) GetRawTVOC() (uint16, error) {
	data := []uint8{0, 0}
	err := d.bus.ReadRegister(d.address, regTVOC, data)
	if err != nil {
		return 0, err
	}

	return uint16(data[0]) | uint16(data[1])<<8, nil
}

// GetRawAQI reads the calculated Air Quality Index (1-5).
func (d *Device) GetRawAQI() (uint8, error) {
	data := []uint8{0}
	err := d.bus.ReadRegister(d.address, regAQI, data)

	return data[0], err
}

func AQIString(value uint8) string {
	switch value {
	case 1:
		return "Excellent"
	case 2:
		return "Good"
	case 3:
		return "Moderate"
	case 4:
		return "Poor"
	case 5:
		return "Unhealthy"
	default:
		return "Unknown"
	}
}

// Reset performs a device reset and leaves the ENS160 in IDLE mode.
func (d *Device) Reset() error {
	// 1) Trigger reset
	if err := d.SetOperatingMode(ModeReset); err != nil {
		return err
	}
	time.Sleep(time.Second)

	// 2) Go to IDLE (default state after reset)
	if err := d.SetOperatingMode(ModeIdle); err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)

	// 3) Clear any old GPR data
	if err := d.bus.WriteRegister(d.address, regCommand, []uint8{ENS160_COMMAND_NOP}); err != nil {
		return err
	}
	time.Sleep(150 * time.Millisecond)
	if err := d.bus.WriteRegister(d.address, regCommand, []uint8{ENS160_COMMAND_CLRGPR}); err != nil {
		return err
	}
	time.Sleep(350 * time.Millisecond)

	return nil
}

// SetEnvData sets the temperature and humidity data for improved accuracy.
// temperature is in Celsius, humidity is in percentage.
func (d *Device) SetEnvData(temperature float32, humidity float32) error {
	// Convert temperature to Kelvin * 64
	tempRaw := uint16((temperature + 273.15) * 64)
	// Convert humidity to percentage * 512
	humRaw := uint16(humidity * 512)

	tempLSB := byte(tempRaw & 0xFF)
	tempMSB := byte(tempRaw >> 8)
	humLSB := byte(humRaw & 0xFF)
	humMSB := byte(humRaw >> 8)

	// write 2 bytes temp @ 0x13, then 2 bytes humidity @ 0x15
	return d.bus.WriteRegister(d.address, regTempIn, []byte{
		tempLSB, tempMSB,
		humLSB, humMSB,
	})
}

// ReadStatus reads the status register of the ENS160.
func (d *Device) ReadStatus() (byte, error) {
	data := []byte{0}
	err := d.bus.ReadRegister(d.address, regStatus, data)
	if err != nil {
		return 0, err
	}

	return data[0], nil
}

/**
 DEVICE_STATUS (Address 0x20)
 This 1-byte register indicates the current status of the ENS160.
 Register structure:
 ------------------------------------------------------------------------
 |   b7   |   b6   |   b5   |   b4  |   b3   |   b2   |   b1   |   b0   |
 ------------------------------------------------------------------------
 | STATAS | STATER |    reserved    |  VALIDITY FLAG  | NEWDAT | NEWGPR |
 ------------------------------------------------------------------------
 Where:
   STATAS:        1 bit  - High indicates that an OPMODE is running
   STATER:        1 bit  - High indicates that an error is detected.
                           E.g. Invalid Operating Mode has been selected.
   VALIDITY FLAG: 2 bits - Status
                           0: Normal operation
                           1: Warm-Up phase
                           2: Initial Start-Up phase
                           3: Invalid output
   reserved:      2 bits - Reserved bits
   NEWDAT:        1 bit  - High indicates that a new data is available
						   in the DATA_x registers. Cleared automatically at
						   first DATA_x read
   NEWGPR:        1 bit  - High indicates that a new data is available
						   in the GPR_READx registers. Cleared automatically
						   at first GPR_READx read.
**/

// ReadGPRDrdy reads the general purpose register data ready flag from the
// status register.
func (d *Device) ReadGPRDrdy() (bool, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return false, err
	}

	gprDrdy := (status & ENS160_DATA_STATUS_NEWGPR) != 0 // Extract bit 0

	return gprDrdy, nil
}

// ReadDataDrdy reads measured data ready flag from the status register.
func (d *Device) ReadDataDrdy() (bool, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return false, err
	}

	dataDrdy := (status & ENS160_DATA_STATUS_NEWDAT) != 0 // Extract bit 1

	return dataDrdy, nil
}

// ReadStatusText reads the status register and returns a human-readable
// string.
func (d *Device) ReadValidityFlag() (uint8, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return 0, err
	}
	// Extract bits 2 and 3
	validityFlag := (status & ENS160_DATA_STATUS_VALIDITY) >> 2

	return validityFlag, nil
}

func ValidityFlagToString(flag uint8) string {
	switch flag {
	case ValidityNormalOperation:
		return "Normal operation"
	case ValidityWarmUpPhase:
		return "Warm-Up phase"
	case ValidityInitialStartUpPhase:
		return "Initial Start-Up phase"
	case ValidityInvalidOutput:
		return "Invalid output"
	default:
		return "Unknown validity flag"
	}
}

// ReadStater reads the stater flag from the status register.
func (d *Device) ReadStater() (bool, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return false, err
	}

	stater := (status & ENS160_DATA_STATUS_STATER) != 0 // Extract bit 6

	return stater, nil
}

// ReadStatas reads the status flag from the status register.
func (d *Device) ReadStatas() (bool, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return false, err
	}

	statusFlag := (status & ENS160_DATA_STATUS_STATAS) != 0 // Extract bit 7

	return statusFlag, nil
}

// ReadStatusText show complete status information from the ENS160.
func (d *Device) ReadStatusText() (string, error) {
	gprDrdy, err := d.ReadGPRDrdy()
	if err != nil {
		return "", err
	}

	dataDrdy, err := d.ReadDataDrdy()
	if err != nil {
		return "", err
	}

	validityFlag, err := d.ReadValidityFlag()
	if err != nil {
		return "", err
	}

	stater, err := d.ReadStater()
	if err != nil {
		return "", err
	}

	statusFlag, err := d.ReadStatas()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"General purpose register data ready: %t, "+
			"Measured data ready: %t, "+
			"Validity: %s, "+
			"Stater: %t, "+
			"Statas: %t",
		gprDrdy, dataDrdy,
		ValidityFlagToString(validityFlag), stater, statusFlag), nil
}

// ReadOptions holds configuration for the Read operation
type ReadOptions struct {
	waitForNew        bool
	withValidityCheck bool
}

// ReadOption is a function that configures ReadOptions
type ReadOption func(*ReadOptions)

// WithWaitForNew configures the Read operation to wait for new data
func WithWaitForNew() ReadOption {
	return func(o *ReadOptions) {
		o.waitForNew = true
	}
}

// WithValidityCheck configures the Read operation to check data validity
func WithValidityCheck() ReadOption {
	return func(o *ReadOptions) {
		o.withValidityCheck = true
	}
}

// Read reads the sensor status and updates all sensor values.
// It accepts optional ReadOption parameters to configure the operation.
func (d *Device) Read(opts ...ReadOption) error {
	options := &ReadOptions{
		waitForNew:        false,
		withValidityCheck: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	status, err := d.ReadStatus()
	if err != nil {
		return fmt.Errorf("error reading status register: %v", err)
	}

	validityFlag := (status & ENS160_DATA_STATUS_VALIDITY) >> 2
	stater := (status & ENS160_DATA_STATUS_STATER) != 0
	dataReady := (status & ENS160_DATA_STATUS_NEWDAT) != 0

	// Check for fatal error state first
	if stater {
		return fmt.Errorf("fatal sensor error (stater flag)")
	}

	if options.waitForNew {
		for !dataReady {
			time.Sleep(time.Millisecond)
			status, err = d.ReadStatus()
			if err != nil {
				return fmt.Errorf("error reading status register: %v", err)
			}
			validityFlag = (status & ENS160_DATA_STATUS_VALIDITY) >> 2
			stater = (status & ENS160_DATA_STATUS_STATER) != 0
			dataReady = (status & ENS160_DATA_STATUS_NEWDAT) != 0

			if stater {
				return fmt.Errorf("fatal sensor error during wait (stater flag)")
			}
		}
	}

	if options.withValidityCheck {
		switch validityFlag {
		case ValidityInitialStartUpPhase:
			return ErrInitialStartUpPhase
		case ValidityWarmUpPhase:
			return ErrWarmUpPhase
		case ValidityInvalidOutput:
			return ErrNoValidOutput
		}
	}

	co2, err := d.GetRawCO2()
	if err != nil {
		return fmt.Errorf("error reading eCO2 data register: %v", err)
	}

	tvoc, err := d.GetRawTVOC()
	if err != nil {
		return fmt.Errorf("error reading TVOC data register: %v", err)
	}

	aqi, err := d.GetRawAQI()
	if err != nil {
		return fmt.Errorf("error reading AQI data register: %v", err)
	}

	d.lastCO2 = co2
	d.lastTVOC = tvoc
	d.lastAQI = aqi

	return nil
}

func (d *Device) LastCO2() uint16 {
	return d.lastCO2
}

func (d *Device) LastTVOC() uint16 {
	return d.lastTVOC
}

func (d *Device) LastAQI() uint8 {
	return d.lastAQI
}
