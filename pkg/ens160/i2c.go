package ens160

// I2C is a hardware-agnostic interface for I2C bus access.
type I2C interface {
	ReadRegister(addr uint8, reg uint8, buf []byte) error
	WriteRegister(addr uint8, reg uint8, buf []byte) error
}
