package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
)

// ConfigDriver used to configure the robot
type ConfigDriver struct {
	name        string
	mpu9250conn MPU9250Sender
	running     bool
	gobot.Eventer
}

// NewConfigDriver return a new ConfigDriver given a UDPWriter and name
func NewConfigDriver(a MPU9250Sender, name string) *ConfigDriver {
	b := &ConfigDriver{
		name:        name,
		mpu9250conn: a,
		running:     true,
		Eventer:     gobot.NewEventer(),
	}
	b.AddEvent(Balance)
	return b
}

// Name returns the ConfigDrivers name
func (b *ConfigDriver) Name() string { return b.name }

// Connection returns the ConfigDrivers Connection
func (b *ConfigDriver) Connection() gobot.Connection { return b.mpu9250conn.(gobot.Connection) }

// Start implements the Driver interface
func (b *ConfigDriver) Start() (errs []error) {
	go b.balanceLoop()
	return
}

// Halt implements the Driver interface
func (b *ConfigDriver) Halt() (errs []error) {
	b.running = false
	return
}

// utility Functions

func checkrc(ok error) {
	if ok != nil {
		log.Println("Error reading MPU9250 : ", ok)
		//b.running = false
	}
}

// Abs function
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Sign check
func Sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

// balanceLoop does the work of balancing the robot
func (b *ConfigDriver) balanceLoop() {

	//***  Variables + Monitoring
	var (
		gAccel, gGyro    int = 1000, 1000
		newG, gGyrodelta int
		speed            int
		gAngle           int
		gAngleInt        int64
		ok               error
	)
	Monitor.Watch(&gAccel, "Accel")
	Monitor.Watch(&gGyro, "Gyro")
	Monitor.Watch(&gAngle, "Angle")
	Monitor.Watch(&gGyrodelta, "Gyrodelta")
	Monitor.Watch(&speed, "Speed")

	//***  MPU monitoring go routine
	runMPU := true
	go func() {
		for runMPU {
			gAccel, ok = b.mpu9250conn.SensorAccel(2)
			checkrc(ok)
			newG, ok = b.mpu9250conn.SensorGyro(1)
			checkrc(ok)
			if newG != gGyro {
				// ignore repeated readings
				gGyrodelta = newG - gGyro
			}
			gGyro = newG

			gAngle = (gAngle+gGyro)*98/100 + (gAccel * 229 * 2 / 100)
			gAngleInt += int64(gAngle)
		}
	}()

	//***  Now start doing some work.  First wait for the robot to be vertical
	log.Println("Hold robot vertically...")
	for Abs(gAccel) > 500 || Abs(gGyro) > 500 {
		time.Sleep(time.Millisecond * 2)
	}
	gAngle = 0
	gAngleInt = 0

	//***  Next give some controlled kicks and see how the robot reacts
	//***  We're looking for the reaction time from kick to action and the shape
	//***  of the action when it happens
	log.Println("Now support loosely in the upright position")
	/*
		speed = 50
		b.Publish(Balance, speed)
		time.Sleep(time.Millisecond * 100)
		speed = -50
		b.Publish(Balance, speed)
		time.Sleep(time.Millisecond * 200)
		setangle := Abs(gAngle)
		dir := Sign(gAngle)

		for {
			speed = 50
			b.Publish(Balance, speed)
			for gAngle*dir*-1 < setangle {
				time.Sleep(time.Millisecond * 2)
			}
			speed = -50
			b.Publish(Balance, speed)
			for gAngle*dir < setangle {
				time.Sleep(time.Millisecond * 2)
			}
		}*/

	for i := 0; i <= 10; i++ {
		speed = i * 20
		if i&1 == 1 {
			speed = -speed
		}
		b.Publish(Balance, speed)
		log.Println("Speed :", speed)
		time.Sleep(time.Millisecond * 100)
	}
	speed = 0
	b.Publish(Balance, speed)

	log.Println("Config complete")
}
