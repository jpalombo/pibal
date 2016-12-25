package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
)

// BalanceDriver used to balance the robot
type BalanceDriver struct {
	name        string
	mpu9250conn MPU9250Sender
	running     bool
	angleOffset int
	gyroFilter  *Filter
	accelFilter *Filter
	gobot.Eventer
}

// NewBalanceDriver return a new BalanceDriver given a UDPWriter and name
func NewBalanceDriver(a MPU9250Sender, name string) *BalanceDriver {
	b := &BalanceDriver{
		name:        name,
		mpu9250conn: a,
		running:     true,
		gyroFilter:  NewFilter(20),
		accelFilter: NewFilter(20),
		Eventer:     gobot.NewEventer(),
	}
	b.AddEvent(BalanceSpeed)
	b.AddEvent(Balancing)
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

// SetAngleOffset sets angle offset
func (b *BalanceDriver) SetAngleOffset(offset int) {
	b.angleOffset = offset
}

// balanceLoop does the work of balancing the robot
func (b *BalanceDriver) balanceLoop() {

	// tracking variables, some of which are monitored
	var (
		gAccelRaw, gAccel              int
		gGyroRaw, gGyro                int
		motorSpeed                     int
		gP, gI, gD                     int
		gKp, gKi, gKd                  int
		gAngle, gLastAngle, gAngleDiff int
		gAngleInt                      int64
		ok                             error
	)
	gKp = -780
	gKi = -2200
	gKd = -4400
	started := false
	Monitor.Watch(&gP, "P")
	Monitor.Watch(&gI, "I")
	Monitor.Watch(&gD, "D")
	Monitor.Watch(&gAccel, "Accel")
	Monitor.Watch(&gAccelRaw, "AccelRaw")
	Monitor.Watch(&motorSpeed, "Speed")
	Monitor.Watch(&gGyro, "Gyro")
	Monitor.Watch(&gGyroRaw, "GyroRaw")
	Monitor.Watch(&gAngle, "Angle")
	Monitor.Control(&gKp, "Kp", 50000, -50000, nil)
	Monitor.Control(&gKi, "Ki", 50000, -50000, nil)
	Monitor.Control(&gKd, "Kd", 50000, -50000, nil)

	// Now into the loop itself
	// Loop takes ~1.9ms to run, so max frequency is just over 500Hz.  We'll set
	// loop time to exactly 2ms, i.e rate of 500 Hz
	//
	// The loop reads the Acceleration and Gyro settings and derives the angle
	// of lean from these.
	// The sensitivities of the two sensors are :
	//		Accel : full scale range (fsr) = 2 G
	//    Gyro  : fsr = 250 degrees per second
	//
	// The Accel reading maps directly to angle.  Specifically
	//   sin(Angle) = Accel reading * 2 / fsr
	// for small Angles sin(Angle) ~= Angle (in Radians).  Error <2% for Angle <20 degrees
	//
	// The angle can also be found from the integral of the Gyro reading
	// if Gyro is sampled 500 times a second and the values summed to give Gyroint then
	//   Angle (degrees) = (Gyroint / 500) / (fsr / 250) = Gyroint * 250 / (fsr * 500) = Gyroint / 2 * fsr
	//
	// Angle degrees = radians * 180 / Pi
	//
	//  => Gyroint * Pi / 2 * fsr * 180 = Accel * 2 /fsr
	//     Gyroint = Accel * 2 * 2 * 180 / Pi ~= 229 * Accel
	//
	// We'll combine the two each by using 98% of the Gyroint number and 2% of the Accel number
	//

	tick := time.Tick(time.Millisecond * 2)
	for b.running {
		<-tick
		//loopcount++
		gAccelRaw, ok = b.mpu9250conn.SensorAccel(2)
		if ok != nil {
			log.Println("Error reading MPU9250 Accel : ", ok)
		}
		gGyroRaw, ok = b.mpu9250conn.SensorGyro(1)
		if ok != nil {
			log.Println("Error reading MPU9250 Gyro : ", ok)
		}

		gAccel = b.accelFilter.Add(gAccelRaw)
		gGyro = b.gyroFilter.Add(gGyroRaw)
		gAngle = (gAngle+gGyro)*98/100 + (gAccel * 229 * 2 / 100)
		gAngleInt += int64(gAngle)
		gAngleDiff = gAngle - gLastAngle
		gLastAngle = gAngle

		// Detect if we are upright and should start
		if !started && Abs(gAccel) < 500 && Abs(gGyro) < 500 {
			started = true
			gAngleInt = 0
			gAngle = 0
			gLastAngle = 0
			log.Println("Start Balancing")
			b.Publish(Balancing, true)
		}

		gP = (gKp * gAngle) >> 22
		gI = int((int64(gKi) * gAngleInt) >> 28)
		gD = (gKd * gAngleDiff) >> 22

		if started {
			motorSpeed = gP + gI + gD
			b.Publish(BalanceSpeed, motorSpeed)
			// Detect if we are at the speed limit and stop if we are
			if Abs(motorSpeed) > 300 {
				started = false
				motorSpeed = 0
				log.Println("Stop Balancing")
				b.Publish(Balancing, false)
				b.Publish(BalanceSpeed, motorSpeed)
			}
		}
	}
}
