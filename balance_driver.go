package main

import "github.com/hybridgroup/gobot"

// BalanceDriver Represents a Joystick
type BalanceDriver struct {
	name       string
	connection MPU9250Sender
	gobot.Eventer
}

// NewBalanceDriver return a new BalanceDriver given a UDPWriter and name
func NewBalanceDriver(a MPU9250Sender, name string) *BalanceDriver {
	b := &BalanceDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
	}
	return b
}

// Name returns the BalanceDrivers name
func (b *BalanceDriver) Name() string { return b.name }

// Connection returns the BalanceDrivers Connection
func (b *BalanceDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start implements the Driver interface
func (b *BalanceDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (b *BalanceDriver) Halt() (errs []error) { return }

// Speed sets the motor speeds
func (b *BalanceDriver) Speed(value ...int16) (err error) {
	return
}
