package main

import (
	"buxiong/car/car"
	"buxiong/car/controller"
	"buxiong/car/driver"
	"buxiong/car/model"
	"os"
	"strconv"
	"strings"

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
	period, err := strconv.ParseUint(os.Getenv(model.Period), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse right moto pin num failed")
		return
	}
	leftSteps, rightSteps := make([]uint64, 0, 10), make([]uint64, 0, 10)
	lss := strings.Split(os.Getenv(model.LeftSteps), ",")
	for _, s := range lss {
		i, e := strconv.ParseUint(s, 10, 64)
		if err != nil {
			err = errors.Wrap(e, "parse left steps failed")
			return
		}
		leftSteps = append(leftSteps, i)
	}
	rss := strings.Split(os.Getenv(model.RightSteps), ",")
	for _, s := range rss {
		i, e := strconv.ParseUint(s, 10, 64)
		if err != nil {
			err = errors.Wrap(e, "parse left steps failed")
			return
		}
		rightSteps = append(rightSteps, i)
	}
	cfg.Period = period
	cfg.LeftSteps, cfg.RightSteps = leftSteps, rightSteps
	a := os.Getenv(model.ListenAddr)
	cfg.Addr = a
	cfg.LeftAPin = uint8(leftA)
	cfg.LeftBPin = uint8(leftB)
	cfg.LeftPWMNum = uint8(leftPWM)
	cfg.RightAPin = uint8(rightA)
	cfg.RightBPin = uint8(rightB)
	cfg.RightPWMNum = uint8(rightPWM)
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
	leftDriver, err := driver.NewDriver(leftA, leftB, cfg.LeftPWMNum, cfg.Period)
	if err != nil {
		panic(err)
	}
	rightDriver, err := driver.NewDriver(rightA, rightB, cfg.RightPWMNum, cfg.Period)
	if err != nil {
		panic(err)
	}
	car := car.NewCar(
		controller.NewController(
			cfg.LeftSteps,
			cfg.RightSteps,
			leftDriver,
			rightDriver,
		),
		cfg.Addr,
	)
	car.Run()
}
