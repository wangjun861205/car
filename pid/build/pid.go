package main

import (
	"buxiong/car/driver"
	"buxiong/car/pid"
	"buxiong/car/speedmeter"
	"buxiong/car/utils"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/stianeikeland/go-rpio/v4"
)

// Input Input
type Input struct {
	*utils.Synchronizer
	target float64
}

// NewInput NewInput
func NewInput(numOfActor int, target float64) *Input {
	return &Input{
		utils.NewSynchronizer(numOfActor),
		target,
	}
}

// Target Target
func (i *Input) Target() float64 {
	return i.target
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		return
	}
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()
	rightDriver, err := driver.NewDriver(rpio.Pin(5), rpio.Pin(6), 0, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	rightMeter, err := speedmeter.NewSpeedMeter(26, 22)
	if err != nil {
		fmt.Println(err)
		return
	}
	leftDriver, err := driver.NewDriver(rpio.Pin(24), rpio.Pin(23), 1, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	leftMeter, err := speedmeter.NewSpeedMeter(17, 27)
	if err != nil {
		fmt.Println(err)
		return
	}
	kp, err := strconv.ParseFloat(os.Getenv("PID_KP"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	ki, err := strconv.ParseFloat(os.Getenv("PID_KI"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	kd, err := strconv.ParseFloat(os.Getenv("PID_KD"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cycle, err := time.ParseDuration(os.Getenv("PID_CYCLE"))
	if err != nil {
		fmt.Println(err)
		return
	}
	target, err := strconv.ParseFloat(os.Getenv("PID_TARGET"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	testDuration, err := time.ParseDuration(os.Getenv("PID_TEST_DURATION"))
	if err != nil {
		fmt.Println(err)
		return
	}
	leftPid := pid.NewPID(kp, ki, kd, cycle, leftDriver, leftMeter)
	rightPid := pid.NewPID(kp, ki, kd, cycle, rightDriver, rightMeter)
	defer leftPid.Close()
	defer rightPid.Close()
	go leftPid.Run()
	go rightPid.Run()
	sync := utils.NewSynchronizer(2)
	leftPid.SetTarget(sync, target)
	rightPid.SetTarget(sync, target)
	sync.WaitReady()
	sync.Active()
	sync.WaitDone()
	time.Sleep(testDuration)

}
