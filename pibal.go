package main

import (
	"time"

	"github.com/hybridgroup/gobot"
)

func main() {
	gbot := gobot.NewGobot()

	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")
	monitorudp := NewUDPAdaptor("monitor UDP", ":25045")
	monitor := NewMonitorDriver(monitorudp, "monitor")
	joystickudp := NewUDPAdaptor("joystick UDP", ":10000")
	remotejoystick := NewPCJoystickDriver(joystickudp, "PC joystick")
	bluetooth := NewBluetoothAdapter("bluetooth", "/dev/input/js0")
	bluetoothjoystick := NewBluetoothDriver(bluetooth, "Bluetooth joystick")
	mpu9250 := NewMPU9250Driver("MPU9250", "I2C")
	balance := NewBalanceDriver(mpu9250, "Balance")

	var a2, g1 int
	monitor.Watch(&a2, "Angle[2]")
	monitor.Watch(&g1, "Gyro[1]")

	work := func() {
		motor.Stop()

		// Handler for commands from a remote joystick
		joystickhandler := func(data interface{}) {
			j := data.(JoystickData)
			if j.deadManHandle {
				gbot.Stop()
			}
			motor.Speed(int16(j.posY*200+j.posX*50), int16(j.posY*200-j.posX*50))
		}
		remotejoystick.On(remotejoystick.Event(Joystick), joystickhandler)
		bluetoothjoystick.On(bluetoothjoystick.Event(Joystick), joystickhandler)

		gobot.Every(20*time.Millisecond, func() {
			a2 = mpu9250.SensorAngle(2)
			g1 = mpu9250.SensorGyro(1)
		})
	}

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{propserial, joystickudp, bluetooth, mpu9250, monitorudp},
		[]gobot.Device{motor, remotejoystick, bluetoothjoystick, balance, monitor},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
