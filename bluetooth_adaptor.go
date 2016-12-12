package main

import (
	"log"
	"os"

	"github.com/hybridgroup/gobot"
)

// BluetoothAdapter is the Gobot Adaptor for UDP communication
type BluetoothAdapter struct {
	adaptorName string
	portName    string
	port        *os.File
	gobot.Eventer
}

// NewBluetoothAdapter returns a new BluetoothAdapter with specified name
//
func NewBluetoothAdapter(name string, portname string) *BluetoothAdapter {
	b := &BluetoothAdapter{
		adaptorName: name,
		portName:    portname,
		Eventer:     gobot.NewEventer(),
	}
	b.AddEvent(Data)
	return b
}

// Connect opens the Bluetooth port and starts a data reader loop.
func (b *BluetoothAdapter) Connect() (errs []error) {
	if b.port == nil {
		var err error
		if b.port, err = os.Open(b.portName); err != nil {
			log.Println("Error opening Bluetooth Adapter on port", b.portName, err)
		}
	}

	go func() {
		inbuf := make([]byte, 8)
		bluebuf := BluetoothData{}
		for {
			n, err := b.port.Read(inbuf)
			if err != nil || n != 8 {
				log.Println("Error : ", err)
			}
			bluebuf.value = int16(inbuf[4]) | int16(inbuf[5])<<8
			bluebuf.bType = inbuf[6]
			bluebuf.bNumber = inbuf[7]
			if bluebuf.bType == 1 || (bluebuf.bType == 2 && bluebuf.bNumber < 8) {
				b.Publish(b.Event(Data), bluebuf)
			}
		}
	}()

	return
}

// Disconnect closes the Bluetooth port
func (b *BluetoothAdapter) Disconnect() (err error) {
	return nil
}

// Finalize terminates the Bluetooth connection
func (b *BluetoothAdapter) Finalize() (errs []error) {
	return
}

// Port returns the  BluetoothAdapters port
func (b *BluetoothAdapter) Port() string { return b.portName }

// Name returns the Bluetooth Adaptors name
func (b *BluetoothAdapter) Name() string { return b.adaptorName }
