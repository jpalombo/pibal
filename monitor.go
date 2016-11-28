package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hybridgroup/gobot"
)

// Watchvar holds data about a watched variable
type Watchvar struct {
	ref    *int
	oldval int
	name   string
}

// Controlvar holds data about a controlled variable
type Controlvar struct {
	ref    *int
	oldval int
	maxval int
	minval int
}

// MonitorDriver holds data about the monitor driver
type MonitorDriver struct {
	name        string
	connection  UDPWriter
	watchvars   []*Watchvar
	controlvars map[string]Controlvar
	running     bool
}

// NewMonitorDriver initialises the MonitorDriver struct
func NewMonitorDriver(u UDPWriter, name string) *MonitorDriver {
	m := &MonitorDriver{
		name:        name,
		connection:  u,
		controlvars: make(map[string]Controlvar),
		running:     true,
	}

	if eventer, ok := u.(gobot.Eventer); ok {
		eventer.On(eventer.Event(Data), func(data interface{}) {
			if inbytes, ok := data.([]byte); ok {
				m.parseMonitorData(inbytes)
			}
		})
	}

	go m.watcher()

	return m
}

// Name returns the MonitorDrivers name
func (m *MonitorDriver) Name() string { return m.name }

// Connection returns the MonitorDrivers Connection
func (m *MonitorDriver) Connection() gobot.Connection { return m.connection.(gobot.Connection) }

// Start implements the Driver interface
func (m *MonitorDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (m *MonitorDriver) Halt() (errs []error) { return }

// Watch starts watching a variable
func (m *MonitorDriver) Watch(watchvar *int, name string) {
	m.watchvars = append(m.watchvars, &Watchvar{ref: watchvar, oldval: *watchvar, name: name})
}

// Control registers a control variable for remote control
func (m *MonitorDriver) Control(controlvar *int, name string, maxval int, minval int) {
	m.controlvars[name] = Controlvar{ref: controlvar, oldval: *controlvar, maxval: maxval, minval: minval}
}

// JSONControl struct used for communication
type JSONControl struct {
	Min int `json:"min"`
	Max int `json:"max"`
	Val int `json:"val"`
}

func (m *MonitorDriver) parseMonitorData(buf []byte) {
	var jdat map[string]interface{}
	if err := json.Unmarshal(buf, &jdat); err != nil {
		panic(err)
	}

	if setmap, ok := jdat["set"]; ok {
		for k, v := range setmap.(map[string]interface{}) {
			if c, ok := m.controlvars[k]; ok {
				*c.ref = int(v.(float64))
			}
		}
	}

	// Send back the current set of control variables
	jout := make(map[string]JSONControl)
	for name, c := range m.controlvars {
		jout[name] = JSONControl{Min: c.minval, Max: c.maxval, Val: *c.ref}
	}
	jdat["control"] = jout

	// Convert the data into JSON
	joutbytes, err := json.Marshal(jdat)
	if err != nil {
		panic(err)
	}

	//Reply to the client
	if err := m.connection.UDPWrite(joutbytes); err != nil {
		log.Println("Error: ", err)
	}
}

const updateinterval = 200

// JSONWatch struct used for communicatng watch variables
type JSONWatch struct {
	Millis    int            `json:"millis"`
	Variables map[string]int `json:"variables"`
}

func (m *MonitorDriver) watcher() {

	loopcount := 0
	dataToSend := false
	j := JSONWatch{}
	t0 := time.Now()

	for m.running {
		time.Sleep(2 * time.Millisecond)
		loopcount++
		j.Millis = int(time.Since(t0) / time.Millisecond)
		j.Variables = make(map[string]int)
		dataToSend = false

		for _, wv := range m.watchvars {
			if (loopcount%updateinterval) == 0 || wv.oldval != *wv.ref {
				// Add the watched variable to the report
				j.Variables[wv.name] = *wv.ref
				wv.oldval = *wv.ref
				dataToSend = true
			}
		}

		if dataToSend {
			if jdata, err := json.Marshal(j); err == nil {
				m.connection.UDPWrite(jdata)
			} else {
				log.Println("Error: ", err)
			}
		}
	}
}
