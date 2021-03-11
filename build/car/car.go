package main

import (
	"buxiong/car/car"
	"buxiong/car/controller"
	"buxiong/car/driver"
	"buxiong/car/electric"
	"buxiong/car/model"
	"buxiong/car/pid"
	"buxiong/car/speedmeter"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio/v4"
	"github.com/warthog618/gpiod"
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
	pidKp, err := strconv.ParseFloat(os.Getenv(model.PIDKp), 64)
	if err != nil {
		err = errors.Wrap(err, "parse pid Kp failed")
		return
	}
	pidKi, err := strconv.ParseFloat(os.Getenv(model.PIDKi), 64)
	if err != nil {
		err = errors.Wrap(err, "parse pid Ki failed")
		return
	}
	pidKd, err := strconv.ParseFloat(os.Getenv(model.PIDKd), 64)
	if err != nil {
		err = errors.Wrap(err, "parse pid Kd failed")
		return
	}
	pidCycle, err := time.ParseDuration(os.Getenv(model.PIDCycle))
	if err != nil {
		err = errors.Wrap(err, "parse pid cycle failed")
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
	mainSwitch, err := strconv.ParseUint(os.Getenv("MAIN_SWITCH_PIN"), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse main switch pin failed")
		return
	}
	cfg.MainSwitchPin = uint8(mainSwitch)
	headLight, err := strconv.ParseUint(os.Getenv("HEAD_LIGHT_PIN"), 10, 64)
	if err != nil {
		err = errors.Wrap(err, "parse headlight pin failed")
		return
	}
	cfg.HeadLightPin = uint8(headLight)
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
	cfg.PIDKp = pidKp
	cfg.PIDKi = pidKi
	cfg.PIDKd = pidKd
	cfg.PIDCycle = pidCycle
	return
}

func newElectric(chip *gpiod.Chip, mainSwitch, headLight uint8) (*electric.Controller, error) {
	mainSwitchLine, err := chip.RequestLine(int(mainSwitch), gpiod.AsOutput())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create main switch")
	}
	headLightLine, err := chip.RequestLine(int(headLight), gpiod.AsOutput())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create headlight")
	}
	return electric.NewController(mainSwitchLine, headLightLine), nil
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
	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}
	leftMeter, err := speedmeter.NewSpeedMeter(chip, int(cfg.LeftEncoderAPin), int(cfg.LeftEncoderBPin))
	if err != nil {
		panic(err)
	}
	rightMeter, err := speedmeter.NewSpeedMeter(chip, int(cfg.RightEncoderAPin), int(cfg.RightEncoderBPin))
	if err != nil {
		panic(err)
	}
	leftPID := pid.NewPID(cfg.PIDKp, cfg.PIDKi, cfg.PIDKd, cfg.PIDCycle, leftDriver, leftMeter)
	rightPID := pid.NewPID(cfg.PIDKp, cfg.PIDKi, cfg.PIDKd, cfg.PIDCycle, rightDriver, rightMeter)
	ctl, err := controller.NewController(
		cfg.Speeds,
		leftPID,
		rightPID,
	)
	if err != nil {
		panic(err)
	}
	elec, err := newElectric(chip, cfg.MainSwitchPin, cfg.HeadLightPin)
	if err != nil {
		panic(err)
	}
	car := car.NewCar(
		ctl,
		elec,
		cfg.Addr,
	)
	car.Run()
}
