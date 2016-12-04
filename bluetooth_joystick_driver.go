package main

import "github.com/hybridgroup/gobot"

// BluetoothDriver Represents a Joystick
type BluetoothDriver struct {
	name       string
	connection gobot.Adaptor
	posX       int16
	posY       int16
	dmh        bool
	gobot.Eventer
}

// NewBluetoothDriver return a new BluetoothDriver given a UDPWriter and name
func NewBluetoothDriver(a gobot.Adaptor, name string) *BluetoothDriver {
	b := &BluetoothDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
	}
	b.AddEvent(Joystick)

	if eventer, ok := a.(gobot.Eventer); ok {
		eventer.On(eventer.Event(Data), func(data interface{}) {
			if joydata, ok := data.(BluetoothData); ok && joydata.bType == 2 {
				if joydata.bType == 2 {
					// type 2 = axis data
					if joydata.bNumber == 0x02 { // x axis
						b.posX = joydata.value
					} else if joydata.bNumber == 0x05 { // y axis
						b.posY = joydata.value
					} else if joydata.bNumber == 0x03 || joydata.bNumber == 0x04 { // dead mans handle
						b.dmh = joydata.value > 0
					}
					b.Publish(Joystick, JoystickData{float64(b.posX) / 32768, float64(b.posY) / 32768, b.dmh})
				}
			}
		})
	}
	return b
}

// Name returns the BluetoothDrivers name
func (b *BluetoothDriver) Name() string { return b.name }

// Connection returns the BluetoothDrivers Connection
func (b *BluetoothDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start implements the Driver interface
func (b *BluetoothDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (b *BluetoothDriver) Halt() (errs []error) { return }

// Speed sets the motor speeds
func (b *BluetoothDriver) Speed(value ...int16) (err error) {
	return
}
