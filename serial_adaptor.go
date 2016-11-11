package main

import (
	"io"
	"log"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

// SerialAdaptor is the Gobot Adaptor for Serial based boards
type SerialAdaptor struct {
	name   string
	port   string
	conn   io.ReadWriteCloser
	openSP func(port string) (io.ReadWriteCloser, error)
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
func NewSerialAdaptor(name string, args ...interface{}) *SerialAdaptor {
	f := &SerialAdaptor{
		name: name,
		port: "",
		conn: nil,
		openSP: func(port string) (io.ReadWriteCloser, error) {
			return serial.OpenPort(&serial.Config{Name: port, Baud: 115200})
		},
		Eventer: gobot.NewEventer(),
	}

	for _, arg := range args {
		switch arg.(type) {
		case string:
			f.port = arg.(string)
		case io.ReadWriteCloser:
			f.conn = arg.(io.ReadWriteCloser)
		}
	}

	return f
}

// Connect starts a connection to the board.
func (f *SerialAdaptor) Connect() (errs []error) {
	if f.conn == nil {
		sp, err := f.openSP(f.Port())
		if err != nil {
			return []error{err}
		}
		f.conn = sp
	}
	return
}

// Disconnect closes the io connection to the board
func (f *SerialAdaptor) Disconnect() (err error) {
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
func (f *SerialAdaptor) Port() string { return f.port }

// Name returns the  SerialAdaptors name
func (f *SerialAdaptor) Name() string { return f.name }

// WriteCmd writes a command to the serial port.
func (f *SerialAdaptor) SerialWrite(cmd string) (err error) {
	log.Println("Sending : " + cmd)

	n, err := f.conn.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = f.conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d %q", n, buf[:n])

	return
}
