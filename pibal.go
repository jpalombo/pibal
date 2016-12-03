package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
)

func main() {
	gbot := gobot.NewGobot()

	var s0, s1, s2, s3 int
	var kp, ki, kd int = 20, 2, 10
	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")
	monitorudp := NewUDPAdaptor("monitor UDP", ":25045")
	monitor := NewMonitorDriver(monitorudp, "monitor")
	joystickudp := NewUDPAdaptor("joystick UDP", ":10000")
	remotejoystick := NewPCJoystickDriver(joystickudp, "PC joystick")

	change := func() {
		log.Printf("Kp = %d, Ki = %d, Kd = %d", kp, ki, kd)
		motor.SetPid(kp, ki, kd)
	}

	monitor.Watch(&s0, "s0")
	monitor.Watch(&s1, "s1")
	monitor.Watch(&s2, "s2")
	monitor.Watch(&s3, "s3")
	monitor.Control(&kp, "Kp", 100, 0, change)
	monitor.Control(&ki, "Ki", 100, 0, change)
	monitor.Control(&kd, "Kd", 100, 0, change)

	work := func() {
		gobot.Every(time.Millisecond*10, func() { motor.GetSpeed() })
		//gobot.Every(time.Millisecond*300, func() { t += inc })

		remotejoystick.On(remotejoystick.Event(Joystick), func(data interface{}) {
			j := data.(JoystickData)
			if j.deadManHandle {
				gbot.Stop()
			}
			motor.Speed(int16(j.posY) * 200)
		})

		motor.On(motor.Event(MotorSpeed), func(data interface{}) {
			mdata := data.(MotorSpeedData)
			s0 = mdata.speed[0]
			s1 = mdata.speed[1]
			s2 = mdata.speed[2]
			s3 = mdata.speed[3]
		})

	}

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{propserial, joystickudp, monitorudp},
		[]gobot.Device{motor, remotejoystick, monitor},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
