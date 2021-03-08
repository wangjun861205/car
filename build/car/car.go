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
	leftEncoderA, err := strconv.ParseInt(os.Getenv(model.LeftEncoderAPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse left encoder a pin num failed")
		return
	}
	leftEncoderB, err := strconv.ParseInt(os.Getenv(model.LeftEncoderBPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse left encoder b pin num failed")
		return
	}
	rightEncoderA, err := strconv.ParseInt(os.Getenv(model.RightEncoderAPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse right encoder a pin num failed")
		return
	}
	rightEncoderB, err := strconv.ParseInt(os.Getenv(model.RightEncoderBPin), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse right encoder b pin num failed")
		return
	}
	speeds := make([]float64, 0, 10)
	ss := strings.Split(os.Getenv(model.Speeds), ",")
	for _, s := range ss {
		v, e := strconv.ParseFloat(s, 64)
		if err != nil {
			err = errors.Wrap(e, "parse left steps failed")
			return
		}
		speeds = append(speeds, v)
	}
	cfg.Speeds = speeds
	a := os.Getenv(model.ListenAddr)
	cfg.Addr = a
	cfg.LeftAPin = uint8(leftA)
	cfg.LeftBPin = uint8(leftB)
	cfg.LeftPWMNum = uint8(leftPWM)
	cfg.LeftEncoderAPin = uint8(leftEncoderA)
	cfg.LeftEncoderBPin = uint8(leftEncoderB)
	cfg.RightAPin = uint8(rightA)
	cfg.RightBPin = uint8(rightB)
	cfg.RightPWMNum = uint8(rightPWM)
	cfg.RightEncoderAPin = uint8(rightEncoderA)
	cfg.RightEncoderBPin = uint8(rightEncoderB)
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
	leftDriver, err := driver.NewDriver(leftA, leftB, cfg.LeftPWMNum, 1000)
	if err != nil {
		panic(err)
	}
	rightDriver, err := driver.NewDriver(rightA, rightB, cfg.RightPWMNum, 1000)
	if err != nil {
		panic(err)
	}
	ctl, err := controller.NewController(
		cfg.Speeds,
		leftDriver,
		rightDriver,
		int(cfg.LeftEncoderAPin),
		int(cfg.LeftEncoderBPin),
		int(cfg.RightEncoderAPin),
		int(cfg.RightEncoderBPin),
	)
	if err != nil {
		panic(err)
	}
	car := car.NewCar(
		ctl,
		cfg.Addr,
	)
	car.Run()
}
