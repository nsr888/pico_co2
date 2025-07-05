// Package ens160 provides a driver for the ScioSense ENS160 digital gas sensor.
//
// Datasheet: https://www.sciosense.com/wp-content/uploads/2023/12/ENS160-Datasheet.pdf
package ens160

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"tinygo.org/x/drivers"
)

const (
	defaultTimeout = 20 * time.Millisecond
	longTimeout    = 1 * time.Second
)

// Device wraps an I2C connection to an ENS160 device.
type Device struct {
	bus  drivers.I2C // I²C implementation
	addr uint16      // 7‑bit bus address, promoted to uint16 per drivers.I2C

	// shadow registers / last measurements
	tvocPPB  uint16
	eco2PPM  uint16
	aqiUBA   uint8
	validity uint8 // Store the latest validity status

	// pre‑allocated buffers (do **not** enlarge at runtime!)
	wbuf [6]byte // longest write: reg + 4 bytes (TEMP+RH)
	rbuf [5]byte // longest read: DATA burst (5 bytes)
}

// New returns a new ENS160 driver.
func New(bus drivers.I2C, address uint16) *Device {
	if address == 0 {
		address = DefaultAddress
	}
	return &Device{bus: bus, addr: address}
}

// Configure sets up the device for reading.
func (d *Device) Configure() error {
	// 1. Soft‑reset
	if err := d.write1(regOpMode, ModeReset); err != nil {
		return err
	}
	time.Sleep(defaultTimeout)

	// 2. Enter IDLE, clear GPR registers, then go STANDARD.
	if err := d.write1(regOpMode, ModeIdle); err != nil {
		return err
	}
	time.Sleep(defaultTimeout)

	if err := d.write1(regCommand, cmdClrGPR); err != nil {
		return err
	}
	time.Sleep(defaultTimeout)

	if err := d.write1(regOpMode, ModeStandard); err != nil {
		return err
	}
	time.Sleep(longTimeout)

	return nil
}

// SetEnvDataMilli sets the ambient temperature and humidity for compensation.
//
// tempMilliC is the temperature in milli-degrees Celsius.
// rhMilliPct is the relative humidity in milli-percent.
func (d *Device) SetEnvDataMilli(tempMilliC, rhMilliPct int32) error {
	// Clip temperature
	const (
		minC = -40 * 1000
		maxC = 85 * 1000
	)
	if tempMilliC < minC {
		tempMilliC = minC
	} else if tempMilliC > maxC {
		tempMilliC = maxC
	}

	// Clip humidity
	if rhMilliPct < 0 {
		rhMilliPct = 0
	} else if rhMilliPct > 100*1000 {
		rhMilliPct = 100 * 1000
	}

	// Integer fixed‑point conversion
	tempRaw := uint16(((tempMilliC + 273150) * 64) / 1000) // Kelvin×64
	humRaw := uint16((rhMilliPct * 512) / 1000)            // %RH×512

	d.wbuf[0] = regTempIn // start address (auto‑increment)
	binary.LittleEndian.PutUint16(d.wbuf[1:3], tempRaw)
	binary.LittleEndian.PutUint16(d.wbuf[3:5], humRaw)

	return d.bus.Tx(d.addr, d.wbuf[:5], nil)
}

// Update refreshes the concentration measurements.
func (d *Device) Update(which drivers.Measurement) error {
	if which&drivers.Concentration == 0 {
		return nil // nothing requested
	}

	const maxTries = 1000
	var (
		status   uint8
		validity uint8
	)
	var gotData bool

	// Poll DEVICE_STATUS until NEWDAT or timeout
	for range maxTries {
		var err error
		status, err = d.read1(regStatus)
		if err != nil {
			return err
		}
		if status&statusSTATER != 0 {
			return errors.New("ENS160: error (STATER set)")
		}
		validity = (status & statusValidityMask) >> statusValidityShift

		if status&statusNEWDAT != 0 {
			gotData = true
			break // Always break when data available
		}
		time.Sleep(time.Millisecond)
	}
	if !gotData {
		return errors.New("ENS160: timeout waiting for NEWDAT")
	}

	// Burst-read data regardless of validity state
	d.wbuf[0] = regAQI
	if err := d.bus.Tx(d.addr, d.wbuf[:1], d.rbuf[:5]); err != nil {
		return fmt.Errorf("ENS160: burst read failed: %w", err)
	}

	d.aqiUBA = d.rbuf[0]
	d.tvocPPB = binary.LittleEndian.Uint16(d.rbuf[1:3])
	d.eco2PPM = binary.LittleEndian.Uint16(d.rbuf[3:5])
	d.validity = validity // Store the validity status

	return nil
}

// TVOC returns the last total‑VOC concentration in parts‑per‑billion.
func (d *Device) TVOC() uint16 { return d.tvocPPB }

// ECO2 returns the last equivalent CO₂ concentration in parts‑per‑million.
func (d *Device) ECO2() uint16 { return d.eco2PPM }

// AQI returns the last Air‑Quality Index according to UBA (1–5).
func (d *Device) AQI() uint8 { return d.aqiUBA }

// Validity returns the current operating state of the sensor.
func (d *Device) Validity() uint8 {
	return d.validity
}

func (d *Device) ValidityString() string {
	switch d.validity {
	case ValidityNormalOperation:
		return "OK: data is valid"
	case ValidityWarmUpPhase:
		return "WARM-UP: needs ~3 min until valid data"
	case ValidityInitialStartUpPhase:
		return "INITIAL START-UP: needs ~1 h until valid data"
	case ValidityInvalidOutput:
		return "INVALID OUTPUT: signals give unexpected values"
	default:
		return "UNKNOWN STATUS"
	}
}

// write1 writes a single byte to a register.
func (d *Device) write1(reg, val uint8) error {
	d.wbuf[0] = reg
	d.wbuf[1] = val
	return d.bus.Tx(d.addr, d.wbuf[:2], nil)
}

// read1 reads a single byte from a register.
func (d *Device) read1(reg uint8) (uint8, error) {
	d.wbuf[0] = reg
	if err := d.bus.Tx(d.addr, d.wbuf[:1], d.rbuf[:1]); err != nil {
		return 0, err
	}
	return d.rbuf[0], nil
}
