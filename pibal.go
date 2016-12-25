package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hybridgroup/gobot"
)

var (
	monitorudp *UDPAdaptor
	// Monitor global used to monitor variables
	Monitor *MonitorDriver

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	configure  = flag.Bool("configure", false, "auto-configure motion parameters")
	monflag    = flag.Bool("monitor", false, "enable remote monitoring")
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	monperiod := 0
	if *monflag {
		monperiod = 2
	}

	gbot := gobot.NewGobot()
	monitorudp = NewUDPAdaptor("monitor UDP", ":25045")
	Monitor = NewMonitorDriver(monitorudp, "monitor", monperiod)
	SpeedMonitor := NewSpeedMon()
	propserial := NewSerialAdaptor("propeller", "/dev/ttyAMA0")
	motor := NewMotorDriver(propserial, "motor")
	joystickudp := NewUDPAdaptor("joystick UDP", ":10000")
	remotejoystick := NewPCJoystickDriver(joystickudp, "PC joystick")
	bluetooth := NewBluetoothAdapter("bluetooth", "/dev/input/js0")
	bluetoothjoystick := NewBluetoothDriver(bluetooth, "Bluetooth joystick")
	mpu9250 := NewMPU9250Driver("MPU9250", "I2C")
	balance := NewBalanceDriver(mpu9250, "Balance")

	work := func() {
		// Initialize vars
		balancing := false
		motordiff := 0
		motor.Stop()

		// Start a go routine to poll for motor position
		positionPoll := true
		go func() {
			tick := time.Tick(time.Millisecond * 100)
			for positionPoll {
				<-tick
				motor.GetPosition()
			}
		}()

		deadzone := func(i float64) int {
			if i < 15 && i > -15 {
				i = 0
			}
			return int(i)
		}

		// Handler for commands from a local or remote joystick
		joystickhandler := func(data interface{}) {
			j := data.(JoystickData)
			if j.deadManHandle {
				pprof.StopCPUProfile()
				gbot.Stop()
			}
			motordiff = deadzone(j.posX * 50)
			motorspeed := deadzone(j.posY * 200)
			if !balancing {
				motor.Speed(motorspeed+motordiff, motorspeed-motordiff)
			} else {
				angleoffset := SpeedMonitor.RequestSpeed(motorspeed / 2)
				balance.SetAngleOffset(angleoffset)
			}
		}

		// Event Handlers

		// remote Joystick
		remotejoystick.On(remotejoystick.Event(Joystick), joystickhandler)

		// Bluetooth Joystick
		bluetoothjoystick.On(bluetoothjoystick.Event(Joystick), joystickhandler)

		// Changed balance state, reset everything
		balance.On(balance.Event(Balancing), func(data interface{}) {
			balancing = data.(bool)
			motordiff = 0
			motor.Stop()
			SpeedMonitor.Reset()
		})

		// Motor speeds from the balance unit
		balance.On(balance.Event(BalanceSpeed), func(data interface{}) {
			if balancing {
				speed := data.(int)
				motor.Speed(speed+motordiff, speed-motordiff)
			}
		})

		// Updated position info
		motor.On(motor.Event(MotorPosition), func(data interface{}) {
			if balancing {
				mp := data.(MotorPositionData)
				angleoffset := SpeedMonitor.UpdatePosition((mp.position[0] + mp.position[3]) / 2)
				balance.SetAngleOffset(angleoffset)
			}
		})
	} // end work function

	robot := gobot.NewRobot("PiBal",
		[]gobot.Connection{monitorudp, propserial, joystickudp, bluetooth, mpu9250},
		[]gobot.Device{Monitor, motor, remotejoystick, bluetoothjoystick, balance},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
