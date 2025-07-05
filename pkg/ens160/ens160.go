/*
Package ens160 provides a driver for the ENS160 Digital Metal-Oxide Multi-Gas
Sensor manufactured by ScioSense.

Datasheet: https://www.sciosense.com/wp-content/uploads/2023/12/ENS160-Datasheet.pdf
*/
package ens160

import (
	"errors"
	"fmt"
	"time"
)

// ErrInitialStartUpPhase indicates the device is in its initial startup phase.
// The device requires 1 hour of continuous operation after initial start-up 
// for adequate readings. After 24 hours of continuous operation, this status
// is stored in non-volatile memory. If unpowered before 24 hours, the device 
// will resume initial startup mode after re-powering.
var ErrInitialStartUpPhase = errors.New("initial startup required")

// ErrWarmUpPhase indicates readings are unavailable during the 3-minute
// warm-up period.
var ErrWarmUpPhase = errors.New("warmup in progress")

// ErrNoValidOutput indicates sensor signals are out of range or invalid.
var ErrNoValidOutput = errors.New("no valid output")

type Device struct {
	bus      I2C
	address  uint8
	buf      [4]byte
	lastCO2  uint16
	lastTVOC uint16
	lastAQI  uint8
}

func New(bus I2C, address uint8) *Device {
	if address == 0 {
		address = DefaultAddress
	}
	return &Device{
		bus:     bus,
		address: address,
	}
}

func (d *Device) Configure() error {
	if err := d.Reset(); err != nil {
		return err
	}
	return d.SetOperatingMode(ModeStandard)
}

func (d *Device) GetOperatingMode() (uint8, error) {
	buf := d.buf[:1]
	err := d.bus.ReadRegister(d.address, regOperatingMode, buf)
	return buf[0], err
}

func (d *Device) SetOperatingMode(mode uint8) error {
	if mode != ModeDeepSleep &&
		mode != ModeIdle &&
		mode != ModeStandard &&
		mode != ModeReset {
		return errors.New("invalid operating mode")
	}
	d.buf[0] = mode
	return d.bus.WriteRegister(d.address, regOperatingMode, d.buf[:1])
}

func (d *Device) GetRawCO2() (uint16, error) {
	buf := d.buf[:2]
	err := d.bus.ReadRegister(d.address, regECO2, buf)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) | uint16(buf[1])<<8, nil
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
	buf := d.buf[:2]
	err := d.bus.ReadRegister(d.address, regTVOC, buf)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) | uint16(buf[1])<<8, nil
}

// GetRawAQI reads the calculated Air Quality Index (1-5).
func (d *Device) GetRawAQI() (uint8, error) {
	buf := d.buf[:1]
	err := d.bus.ReadRegister(d.address, regAQI, buf)
	return buf[0], err
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
	d.buf[0] = CommandNop
	if err := d.bus.WriteRegister(
		d.address,
		regCommand,
		d.buf[:1],
	); err != nil {
		return err
	}
	time.Sleep(150 * time.Millisecond)
	d.buf[0] = CommandClrGpr
	if err := d.bus.WriteRegister(
		d.address,
		regCommand,
		d.buf[:1],
	); err != nil {
		return err
	}
	time.Sleep(350 * time.Millisecond)

	return nil
}

// SetEnvData sets the temperature and humidity data for improved accuracy.
// tempMilliC is temperature in milli-Celsius, humidityMilliPct is humidity in milli-percent (e.g. 45670 = 45.67%).
func (d *Device) SetEnvData(tempMilliC int32, humidityMilliPct int32) error {
	// Convert temperature to Kelvin * 64, using integer math:
	// tempKelvin = tempMilliC/1000 + 273.15
	// tempRaw = (tempKelvin * 64)
	//         = ((tempMilliC + 273150) * 64) / 1000
	tempRaw := uint16(((tempMilliC + 273150) * 64) / 1000)
	// Convert humidity to percentage * 512, using integer math:
	// humRaw = (humidityMilliPct * 512) / 1000
	humRaw := uint16((humidityMilliPct * 512) / 1000)

	d.buf[0] = byte(tempRaw & 0xFF)
	d.buf[1] = byte(tempRaw >> 8)
	d.buf[2] = byte(humRaw & 0xFF)
	d.buf[3] = byte(humRaw >> 8)

	return d.bus.WriteRegister(d.address, regTempIn, d.buf[:4])
}

// ReadStatus reads the status register of the ENS160.
func (d *Device) ReadStatus() (byte, error) {
	buf := d.buf[:1]
	err := d.bus.ReadRegister(d.address, regStatus, buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
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
	gprDrdy := (status & DataStatusNewGpr) != 0 // Extract bit 0
	return gprDrdy, nil
}

// ReadDataDrdy reads measured data ready flag from the status register.
func (d *Device) ReadDataDrdy() (bool, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return false, err
	}
	dataDrdy := (status & DataStatusNewDat) != 0 // Extract bit 1
	return dataDrdy, nil
}

// ReadValidityFlag reads the status register and returns the validity flag.
func (d *Device) ReadValidityFlag() (uint8, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return 0, err
	}
	// Extract bits 2 and 3
	validityFlag := (status & DataStatusValidity) >> 2
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
	stater := (status & DataStatusStater) != 0 // Extract bit 6
	return stater, nil
}

// ReadStatas reads the status flag from the status register.
func (d *Device) ReadStatas() (bool, error) {
	status, err := d.ReadStatus()
	if err != nil {
		return false, err
	}
	statusFlag := (status & DataStatusStatas) != 0 // Extract bit 7
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
		"general purpose register data ready: %t, "+
			"measured data ready: %t, "+
			"validity: %s, "+
			"stater: %t, "+
			"statas: %t",
		gprDrdy, dataDrdy,
		ValidityFlagToString(validityFlag), stater, statusFlag), nil
}

// ReadConfig holds configuration for the Read operation
type ReadConfig struct {
	WaitForNew        bool
	WithValidityCheck bool
}

// Read reads the sensor status and updates all sensor values.
func (d *Device) Read(cfg ReadConfig) error {
	status, err := d.ReadStatus()
	if err != nil {
		return fmt.Errorf("error reading status register: %w", err)
	}

	validityFlag := (status & DataStatusValidity) >> 2
	stater := (status & DataStatusStater) != 0
	dataReady := (status & DataStatusNewDat) != 0

	// Check for fatal error state first
	if stater {
		return fmt.Errorf("fatal sensor error (stater flag)")
	}

	if cfg.WaitForNew {
		retries := 500 // ~500ms timeout
		for !dataReady && retries > 0 {
			time.Sleep(time.Millisecond)
			status, err = d.ReadStatus()
			if err != nil {
				return fmt.Errorf("error reading status register: %w", err)
			}
			validityFlag = (status & DataStatusValidity) >> 2
			stater = (status & DataStatusStater) != 0
			dataReady = (status & DataStatusNewDat) != 0

			if stater {
				return fmt.Errorf("fatal sensor error during wait (stater flag)")
			}
			retries--
		}
		if retries == 0 {
			return errors.New("timeout waiting for new data")
		}
	}

	co2, err := d.GetRawCO2()
	if err != nil {
		return fmt.Errorf("error reading eCO2 data register: %w", err)
	}

	tvoc, err := d.GetRawTVOC()
	if err != nil {
		return fmt.Errorf("error reading TVOC data register: %w", err)
	}

	aqi, err := d.GetRawAQI()
	if err != nil {
		return fmt.Errorf("error reading AQI data register: %w", err)
	}

	d.lastCO2 = co2
	d.lastTVOC = tvoc
	d.lastAQI = aqi

	if cfg.WithValidityCheck {
		switch validityFlag {
		case ValidityInitialStartUpPhase:
			return fmt.Errorf(
				"allow 60 minutes for adequate readings: %w",
				ErrInitialStartUpPhase,
			)
		case ValidityWarmUpPhase:
			return fmt.Errorf(
				"allow 3 minutes for adequate readings: %w",
				ErrWarmUpPhase,
			)
		case ValidityInvalidOutput:
			return ErrNoValidOutput
		}
	}

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
