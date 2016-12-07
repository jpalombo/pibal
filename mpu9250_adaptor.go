package main

import (
	"github.com/hybridgroup/gobot"

	//#include "sensor.h"
	"C"
)

// MPU9250Driver is the Gobot Adaptor for UDP communication
type MPU9250Driver struct {
	adaptorName string
	portName    string
	gobot.Eventer
}

// NewMPU9250Driver returns a new MPU9250Driver with specified name
//
func NewMPU9250Driver(name string, portname string) *MPU9250Driver {
	b := &MPU9250Driver{
		adaptorName: name,
		portName:    portname,
		Eventer:     gobot.NewEventer(),
	}
	b.AddEvent(Data)
	return b
}

// Connect opens the Bluetooth port and starts a data reader loop.
func (b *MPU9250Driver) Connect() (errs []error) {
	C.mpu_open()
	return
}

// Disconnect closes the Bluetooth port
func (b *MPU9250Driver) Disconnect() (err error) {
	return nil
}

// Finalize terminates the Bluetooth connection
func (b *MPU9250Driver) Finalize() (errs []error) {
	return
}

// Port returns the  MPU9250Drivers port
func (b *MPU9250Driver) Port() string { return b.portName }

// Name returns the Bluetooth Adaptors name
func (b *MPU9250Driver) Name() string { return b.adaptorName }

// SensorAngle call the C func
func (b *MPU9250Driver) SensorAngle(i int) int {
	return int(C.sensorAngle(C.int(i)))
}

// SensorGyro calls the C func
func (b *MPU9250Driver) SensorGyro(i int) int {
	return int(C.sensorGyro(C.int(i)))
}
