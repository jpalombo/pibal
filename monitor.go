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
	missed bool
	name   string
}

// Controlvar holds data about a controlled variable
type Controlvar struct {
	ref        *int
	oldval     int
	maxval     int
	minval     int
	changefunc func()
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
func NewMonitorDriver(u UDPWriter, name string, periodms int) *MonitorDriver {
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

	if periodms > 0 {
		go m.watcher(periodms)
	}

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
func (m *MonitorDriver) Control(controlvar *int, name string, maxval int, minval int, changefn func()) {
	m.controlvars[name] = Controlvar{ref: controlvar, oldval: *controlvar, maxval: maxval, minval: minval, changefunc: changefn}
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
		log.Println("Error in Monitor:", err)
		return
	}

	if setmap, ok := jdat["set"]; ok {
		for k, v := range setmap.(map[string]interface{}) {
			if c, ok := m.controlvars[k]; ok {
				*c.ref = int(v.(float64))
				if c.changefunc != nil {
					c.changefunc()
				}
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
		log.Println("Error in Monitor:", err)
		return
	}

	//Reply to the client
	if err := m.connection.UDPWrite(joutbytes); err != nil {
		log.Println("Error in Monitor: ", err)
	}
}

const updateinterval = 200

// JSONWatch struct used for communicatng watch variables
type JSONWatch struct {
	Millis    int            `json:"millis"`
	Variables map[string]int `json:"variables"`
}

func (m *MonitorDriver) watcher(period int) {

	loopcount := 0 // used to send values occassionally, even if they haven't changed
	newval := JSONWatch{Variables: make(map[string]int)}
	oldval := JSONWatch{Variables: make(map[string]int)}

	t0 := time.Now()

	tick := time.Tick(time.Millisecond * time.Duration(period))
	for m.running {
		<-tick
		loopcount++
		oldval.Millis = newval.Millis // use the time value from the last loop
		newval.Millis = int(time.Since(t0) / time.Millisecond)

		// Look for variables that have changed value
		for _, wv := range m.watchvars {
			if (loopcount%updateinterval) == 0 || wv.oldval != *wv.ref {
				// Add the watched variable to the report
				if wv.missed {
					oldval.Variables[wv.name] = wv.oldval
				}
				newval.Variables[wv.name] = *wv.ref
				wv.oldval = *wv.ref
				wv.missed = false
			} else {
				wv.missed = true
			}
		}

		// send any variables that have changed value
		if len(newval.Variables) > 0 {
			// first send the old value if any have been missed
			if len(oldval.Variables) > 0 {
				if jdata, err := json.Marshal(oldval); err == nil {
					m.connection.UDPWrite(jdata)
				} else {
					log.Println("Error: ", err)
				}
				oldval.Variables = make(map[string]int)
			}

			// now send the new value
			if jdata, err := json.Marshal(newval); err == nil {
				m.connection.UDPWrite(jdata)
			} else {
				log.Println("Error: ", err)
			}
			newval.Variables = make(map[string]int)
		}
	}
}
