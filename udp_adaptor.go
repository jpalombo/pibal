package main

import (
	"net"

	"github.com/hybridgroup/gobot"
)

const bufsize = 1024

// UDPAdaptor is the Gobot Adaptor for UDP communication
type UDPAdaptor struct {
	adaptorName string
	portName    string
	localAddr   *net.UDPAddr
	remoteAddr  *net.UDPAddr
	conn        *net.UDPConn
	gobot.Eventer
}

// NewUDPAdaptor returns a new UDPAdaptor with specified name
//
func NewUDPAdaptor(name string, portname string) *UDPAdaptor {
	u := &UDPAdaptor{
		adaptorName: name,
		portName:    portname,
		localAddr:   nil,
		remoteAddr:  nil,
		conn:        nil,
		Eventer:     gobot.NewEventer(),
	}
	u.AddEvent(Data)
	return u
}

// Connect opens a UDP port and starts a data reader loop.
func (u *UDPAdaptor) Connect() (errs []error) {
	if u.localAddr == nil {
		u.localAddr, _ = net.ResolveUDPAddr("udp", u.portName)
	}
	if u.conn == nil {
		u.conn, _ = net.ListenUDP("udp", u.localAddr)
	}

	go func() {
		buf := make([]byte, bufsize)
		var err error
		n := 0
		for err == nil {
			if n, u.remoteAddr, err = u.conn.ReadFromUDP(buf); err == nil {
				u.Publish(u.Event(Data), buf[:n])
			}
		}
	}()

	return
}

// Disconnect closes the UDP port
func (u *UDPAdaptor) Disconnect() (err error) {
	u.conn.Close()
	return nil
}

// Finalize terminates the UDP connection
func (u *UDPAdaptor) Finalize() (errs []error) {
	return
}

// Port returns the  UDPAdaptors port
func (u *UDPAdaptor) Port() string { return u.portName }

// Name returns the  SerialAdaptors name
func (u *UDPAdaptor) Name() string { return u.adaptorName }

// UDPWrite sends a UDP packet back to the last place that we received one from.
func (u *UDPAdaptor) UDPWrite(data []byte) (err error) {
	if u.conn != nil && u.remoteAddr != nil {
		_, err = u.conn.WriteToUDP(data, u.remoteAddr)
	}
	return
}
