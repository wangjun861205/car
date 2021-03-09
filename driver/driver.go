package driver

import (
	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio/v4"
)

// Driver motor driver
type Driver struct {
	period    uint64
	dutyCycle int64
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

// NewDriver NewDriver
func NewDriver(aPin, bPin rpio.Pin, pwmNum uint8, period uint64) (*Driver, error) {
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
	return &Driver{
		period:    period,
		dutyCycle: 0,
		aPin:      aPin,
		bPin:      bPin,
		pwm:       pwm,
	}, nil
}

// Stop Stop
func (d *Driver) Stop() {
	d.aPin.Low()
	d.bPin.Low()
	d.pwm.SetDutyCycle(0)
	d.dutyCycle = 0
}

// SetDuty SetDuty
func (d *Driver) SetDuty(duty float64) {
	if duty > 0 {
		d.aPin.Low()
		d.bPin.High()
		if int64(duty) > int64(d.period) {
			d.pwm.SetDutyCycle(d.period)
			d.dutyCycle = int64(d.period)
		} else {
			d.pwm.SetDutyCycle(uint64(duty))
			d.dutyCycle = int64(duty)
		}
	} else {
		d.aPin.High()
		d.bPin.Low()
		if -int64(duty) > int64(d.period) {
			d.pwm.SetDutyCycle(d.period)
			d.dutyCycle = -int64(d.period)
		} else {
			d.pwm.SetDutyCycle(uint64(-duty))
			d.dutyCycle = -int64(duty)
		}
	}
}

// GetDuty GetDuty
func (d *Driver) GetDuty() float64 {
	return float64(d.dutyCycle)
}

// Close Close
func (d *Driver) Close() {
	d.pwm.Close()
}
