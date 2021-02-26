package main

import (
	"buxiong/car/car"
	"buxiong/car/controller"
	"buxiong/car/driver"
	"buxiong/car/model"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio/v4"
)

func getConfig() (cfg model.Config, err error) {
	leftA, err := strconv.ParseInt(os.Getenv(model.LeftAPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse left moto pin num failed")
		return
	}
	leftB, err := strconv.ParseInt(os.Getenv(model.LeftBPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse left moto pin num failed")
		return
	}
	leftPWM, err := strconv.ParseInt(os.Getenv(model.LeftPWMPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse left moto pin num failed")
		return
	}
	rightA, err := strconv.ParseInt(os.Getenv(model.RightAPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse right moto pin num failed")
		return
	}
	rightB, err := strconv.ParseInt(os.Getenv(model.RightBPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse right moto pin num failed")
		return
	}
	rightPWM, err := strconv.ParseInt(os.Getenv(model.RightPWMPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse right moto pin num failed")
		return
	}
	m, err := strconv.ParseUint(os.Getenv(model.Max), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse max failed")
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
	cfg.Addr = a
	cfg.Max = m
	cfg.Base = b
	cfg.LeftAPin = uint8(leftA)
	cfg.LeftBPin = uint8(leftB)
	cfg.LeftPWMNum = uint8(leftPWM)
	cfg.RightAPin = uint8(rightA)
	cfg.RightBPin = uint8(rightB)
	cfg.RightPWMNum = uint8(rightPWM)
	cfg.Step = s
	return
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	leftA := rpio.Pin(cfg.LeftAPin)
	leftB := rpio.Pin(cfg.LeftBPin)
	rightA := rpio.Pin(cfg.RightAPin)
	rightB := rpio.Pin(cfg.RightBPin)
	leftDriver, err := driver.NewDriver(leftA, leftB, cfg.LeftPWMNum, cfg.Max)
	if err != nil {
		panic(err)
	}
	rightDriver, err := driver.NewDriver(rightA, rightB, cfg.RightPWMNum, cfg.Max)
	if err != nil {
		panic(err)
	}
	car := car.NewCar(
		controller.NewController(
			cfg.Max,
			cfg.Base,
			cfg.Step,
			leftDriver,
			rightDriver,
		),
		cfg.Addr,
	)
	car.Run()
}
