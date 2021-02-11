package pwm

import (
	"buxiong/car/model"
	"time"

	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio"
)

type pwm struct {
	base       uint8
	step       uint8
	proportion uint8
	inChan     chan model.PWMAction
	statusChan chan model.PWMStatus
	ticker     *time.Ticker
	doneIn     chan struct{}
	doneOut    chan struct{}
	pin        rpio.Pin
}

func NewPWM(base, step uint8, pin rpio.Pin) (*pwm, error) {
	if base > 100 {
		return nil, errors.Errorf("invalid base proportion(proportion: %d)", base)
	}
	if step > 100 {
		return nil, errors.Errorf("invalid step(step: %d)", step)
	}
	return &pwm{
		base,
		step,
		0,
		make(chan model.PWMAction),
		make(chan model.PWMStatus),
		time.NewTicker(time.Microsecond * 100),
		make(chan struct{}),
		make(chan struct{}),
		pin,
	}, nil
}

func (p *pwm) handleInc() {
	switch {
	case p.proportion == 0:
		p.proportion = p.base
	case p.proportion+p.step > 100:
		p.proportion = 100
	default:
		p.proportion += p.step
	}
}

func (p *pwm) handleDec() {
	switch {
	case p.proportion == 0:
		return
	case p.proportion-p.step < p.base:
		p.proportion = 0
	default:
		p.proportion -= p.step
	}
}

func (p *pwm) Run() {
	go func() {
	OUTER:
		for {
			select {
			case action := <-p.inChan:
				switch action.Type {
				case model.PWMActionTypeIncrement:
					p.handleInc()
				case model.PWMActionTypeDecrement:
					p.handleDec()
				case model.PWMActionTypeStatus:
					p.statusChan <- model.PWMStatus{p.base, p.step, p.proportion}
				}
			case <-p.doneIn:
				p.ticker.Stop()
				p.pin.Low()
				p.doneOut <- struct{}{}
				return
			default:
				<-p.ticker.C
				if p.proportion == 0 {
					continue OUTER
				}
				p.pin.High()
				timer := time.NewTimer(time.Microsecond * time.Duration(p.proportion))
				<-timer.C
				p.pin.Low()

			}
		}
	}()
}

func (p *pwm) Stop() {
	close(p.doneIn)
	<-p.doneOut
}

func (p *pwm) Status() model.PWMStatus {
	p.inChan <- model.PWMAction{
		Type: model.PWMActionTypeStatus,
	}
	return <-p.statusChan
}

func (p *pwm) Increment() {
	p.inChan <- model.PWMAction{
		Type: model.PWMActionTypeIncrement,
	}
}

func (p *pwm) Decrement() {
	p.inChan <- model.PWMAction{
		Type: model.PWMActionTypeDecrement,
	}
}
