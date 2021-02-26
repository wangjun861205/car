package driver

import (
	"buxiong/car/model"

	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio/v4"
)

type driver struct {
	direction model.Direction
	period    uint64
	dutyCycle uint64
	aPin      rpio.Pin
	bPin      rpio.Pin
	pwm       *PWM
}

func initPWM(pwm *PWM, period uint64) error {
	dutyCycle, err := pwm.GetDutyCycle()
	if err != nil {
		return errors.Wrap(err, "failed to init pwm")
	}
	switch {
	case period < dutyCycle:
		if err := pwm.SetDutyCycle(0); err != nil {
			return errors.Wrap(err, "failed to init pwm")
		}
		if err := pwm.SetPeriod(period); err != nil {
			return errors.Wrap(err, "failed to init pwm")
		}
	default:
		if err := pwm.SetPeriod(period); err != nil {
			return errors.Wrap(err, "failed to init pwm")
		}
		if err := pwm.SetDutyCycle(0); err != nil {
			return errors.Wrap(err, "failed to init pwm")
		}
	}
	if err := pwm.Enable(); err != nil {
		return errors.Wrap(err, "failed to init pwm")
	}
	return nil
}

func NewDriver(aPin, bPin rpio.Pin, pwmNum uint8, period uint64) (*driver, error) {
	aPin.Output()
	aPin.High()
	bPin.Output()
	bPin.High()
	pwm, err := NewPWM(pwmNum)
	if err != nil {
		return nil, err
	}
	if err := initPWM(pwm, period); err != nil {
		pwm.Close()
		return nil, err
	}
	return &driver{
		direction: model.DirectionGlide,
		period:    period,
		dutyCycle: 0,
		aPin:      aPin,
		bPin:      bPin,
		pwm:       pwm,
	}, nil
}

func (d *driver) Brake() {
	d.aPin.Low()
	d.bPin.Low()
	d.pwm.SetDutyCycle(0)
	d.direction = model.DirectionBrake
	d.dutyCycle = 0
}

func (d *driver) Glide() {
	d.aPin.High()
	d.bPin.High()
	d.pwm.SetDutyCycle(0)
	d.direction = model.DirectionGlide
	d.dutyCycle = 0
}

func (d *driver) Forward(duty uint64) {
	d.aPin.Low()
	d.bPin.High()
	d.pwm.SetDutyCycle(duty)
	d.direction = model.DirectionForward
	d.dutyCycle = duty
}

func (d *driver) Backward(duty uint64) {
	d.aPin.High()
	d.bPin.Low()
	d.pwm.SetDutyCycle(duty)
	d.direction = model.DirectionBackward
	d.dutyCycle = duty
}

func (d *driver) Status() *model.DriverStatus {
	return model.NewDriverStatus(d.direction, d.dutyCycle)
}

func (d *driver) Close() error {
	return d.pwm.Close()
}
