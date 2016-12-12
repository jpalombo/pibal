package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
)

// BalanceDriver Represents a Joystick
type BalanceDriver struct {
	name        string
	mpu9250conn MPU9250Sender
	running     bool
	gobot.Eventer
}

// NewBalanceDriver return a new BalanceDriver given a UDPWriter and name
func NewBalanceDriver(a MPU9250Sender, name string) *BalanceDriver {
	b := &BalanceDriver{
		name:        name,
		mpu9250conn: a,
		running:     true,
		Eventer:     gobot.NewEventer(),
	}
	b.AddEvent(Balance)
	return b
}

// Name returns the BalanceDrivers name
func (b *BalanceDriver) Name() string { return b.name }

// Connection returns the BalanceDrivers Connection
func (b *BalanceDriver) Connection() gobot.Connection { return b.mpu9250conn.(gobot.Connection) }

// Start implements the Driver interface
func (b *BalanceDriver) Start() (errs []error) {
	go b.balanceLoop()
	return
}

// Halt implements the Driver interface
func (b *BalanceDriver) Halt() (errs []error) {
	b.running = false
	return
}

// balanceLoop does the work of balancing the robot
func (b *BalanceDriver) balanceLoop() {

	// First some housekeeping and utility routines
	// code to track how fast we are looping
	/*	loopcount := 0
		gobot.Every(time.Second, func() {
			log.Println("Loops per sec = ", loopcount)
			loopcount = 0
		}) */
	// a return code checker
	checkrc := func(ok error) {
		if ok != nil {
			log.Println("Error reading MPU9250 : ", ok)
			//b.running = false
		}
	}
	// an Abs function
	Abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	// tracking variables, some of which are monitored
	var (
		gAngle, gGyro, newspeed   int
		gP, gI, gD, gKp, gKi, gKd int
		gGyroint, gGyrointint     int64
		ok                        error
	)
	gKp = -10000
	gKi = -0
	gKd = 0
	started := false
	Monitor.Watch(&gP, "P")
	Monitor.Watch(&gI, "I")
	Monitor.Watch(&gD, "D")
	Monitor.Watch(&newspeed, "Speed")
	Monitor.Control(&gKp, "Kp", 50000, -50000, nil)
	Monitor.Control(&gKi, "Ki", 50000, -50000, nil)
	Monitor.Control(&gKd, "Kd", 50000, -50000, nil)

	// Now into the loop itself
	// Loop takes ~1.9ms to run, so max frequency is just over 500Hz
	// No point looping faster than can issue motor commands, which take 10-15ms
	tick := time.Tick(time.Millisecond * 15)
	for b.running {
		<-tick
		//loopcount++
		gAngle, ok = b.mpu9250conn.SensorAngle(2)
		checkrc(ok)
		gGyro, ok = b.mpu9250conn.SensorGyro(1)
		checkrc(ok)
		gGyroint += int64(gGyro)
		gGyrointint += int64(gGyroint)

		// Detect if we are upright and should start
		if !started && Abs(gAngle) < 500 && Abs(gGyro) < 500 {
			started = true
			gGyroint = 0
			gGyrointint = 0
			log.Println("Start Balancing")
		}

		gP = int((int64(gKp) * gGyroint) >> 22)
		gI = int((int64(gKi) * gGyrointint) >> 22)
		gD = int((int64(gKd) * int64(gGyro)) >> 22)

		if started {
			newspeed = gP + gI + gD
			// Detect if we are at the speed limit and stop if we are
			if Abs(newspeed) > 200 {
				started = false
				newspeed = 0
				log.Println("Stop Balancing")
			}
			b.Publish(Balance, newspeed)
		}
	}
}
