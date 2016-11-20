package main

import (
	"log"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/serial"
)

// SerialAdaptor is the Gobot Adaptor for Serial based boards
type SerialAdaptor struct {
	adaptorName string
	portName    string
	port        *serial.Port
	gobot.Eventer
}

// NewSerialAdaptor returns a new SerialAdaptor with specified name and optionally accepts:
//
//	string: port the SerialAdaptor uses to connect to a serial port with a baud rate of 115200
//	io.ReadWriteCloser: connection the SerialAdaptor uses to communication with the hardware
//
// If an io.ReadWriteCloser is not supplied, the SerialAdaptor will open a connection
// to a serial port with a baude rate of 57600. If an io.ReadWriteCloser
// is supplied, then the SerialAdaptor will use the provided io.ReadWriteCloser and use the
// string port as a label to be displayed in the log and api.
func NewSerialAdaptor(name string, portname string) *SerialAdaptor {
	f := &SerialAdaptor{
		adaptorName: name,
		portName:    portname,
		port:        nil,
		Eventer:     gobot.NewEventer(),
	}
	f.AddEvent(Data)
	f.AddEvent(Error)
	return f
}

// Connect starts a connection to the board and start a data reader loop.
func (f *SerialAdaptor) Connect() (errs []error) {
	if f.port == nil {
		c := &serial.Config{Name: f.portName, Baud: 115200}
		sp, err := serial.OpenPort(c)
		if err != nil {
			return []error{err}
		}
		f.port = sp
	}

	go func() {
		inchar := make([]byte, 1)
		line := make([]byte, 0, 128)
		for {
			n, err := f.port.Read(inchar)
			if err != nil || n != 1 {
				f.Publish(f.Event(Error), err)
				log.Fatal(err)
			} else if inchar[0] == '\n' || inchar[0] == '\r' {
				//log.Printf("%d %q", n, line)
				f.Publish(f.Event(Data), line)
				line = line[0:0]
			} else {
				line = append(line, inchar[0])
			}
		}
	}()

	return
}

// Disconnect closes the io connection to the board
func (f *SerialAdaptor) Disconnect() (err error) {
	f.port.Close()
	return nil
}

// Finalize terminates the Serial connection
func (f *SerialAdaptor) Finalize() (errs []error) {
	if err := f.Disconnect(); err != nil {
		return []error{err}
	}
	return
}

// Port returns the  SerialAdaptors port
func (f *SerialAdaptor) Port() string { return f.portName }

// Name returns the  SerialAdaptors name
func (f *SerialAdaptor) Name() string { return f.adaptorName }

// WriteCmd writes a command to the serial port.
func (f *SerialAdaptor) SerialWrite(cmd string) (err error) {
	log.Println("Sending : " + cmd)
	_, err = f.port.Write([]byte(cmd + "\n"))
	if err != nil {
		f.Publish(f.Event(Error), err)
		log.Fatal(err)
	}
	return
}
