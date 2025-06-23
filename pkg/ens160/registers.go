package ens160

const (
	// Default I2C address for ENS160
	DefaultAddress = 0x53

	// Register addresses
	regOperatingMode = 0x10
	regCommand       = 0x12
	regTempIn        = 0x13
	regRHumIn        = 0x15
	regStatus        = 0x20
	regAQI           = 0x21
	regTVOC          = 0x22
	regECO2          = 0x24

	// Operating modes
	ModeDeepSleep = 0x00 // low-power standby
	ModeIdle      = 0x01 // low power
	ModeStandard  = 0x02 // Gas Sensing Mode
	ModeReset     = 0xF0

	// Sensor validity flags uint8
	ValidityNormalOperation     = 0x00
	ValidityWarmUpPhase         = 0x01
	ValidityInitialStartUpPhase = 0x02
	ValidityInvalidOutput       = 0x03

	// Data status flags
	ENS160_DATA_STATUS_STATAS   = 0x80
	ENS160_DATA_STATUS_STATER   = 0x40
	ENS160_DATA_STATUS_VALIDITY = 0x0C
	ENS160_DATA_STATUS_NEWDAT   = 0x02
	ENS160_DATA_STATUS_NEWGPR   = 0x01

	// Command codes
	ENS160_COMMAND_NOP        = 0x00
	ENS160_COMMAND_GET_APPVER = 0x0E // Get FW Version
	ENS160_COMMAND_CLRGPR     = 0xCC // Clears GPR Read Registers
)
