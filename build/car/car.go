package main

import (
	"buxiong/car/car"
	"buxiong/car/controller"
	"buxiong/car/model"
	"buxiong/car/network"
	"buxiong/car/pwm"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio"
)

func getConfig() (leftPinNum, rightPinNum int, base, step uint8, addr string, err error) {
	left, err := strconv.ParseInt(os.Getenv(model.LeftMotoPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse left moto pin num failed")
		return
	}
	right, err := strconv.ParseInt(os.Getenv(model.RightMotoPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse right moto pin num failed")
		return
	}
	b, err := strconv.ParseUint(os.Getenv(model.Base), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse base failed")
		return
	}
	s, err := strconv.ParseUint(os.Getenv(model.Step), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse step failed")
		return
	}
	a := os.Getenv(model.ListenAddr)
	return int(left), int(right), uint8(b), uint8(s), a, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()
	leftPinNum, rightPinNum, base, step, addr, err := getConfig()
	if err != nil {
		panic(err)
	}
	leftPin := rpio.Pin(leftPinNum)
	rightPin := rpio.Pin(rightPinNum)
	leftPin.Output()
	rightPin.Output()
	leftPWM, err := pwm.NewPWM(base, step, leftPin)
	if err != nil {
		panic(err)
	}
	rightPWM, err := pwm.NewPWM(base, step, rightPin)
	if err != nil {
		panic(err)
	}
	server, err := network.NewServer(addr, '\n')
	if err != nil {
		panic(err)
	}
	car := car.NewCar(
		controller.NewController(
			leftPWM,
			rightPWM,
		),
		server,
	)
	car.Run()
}
