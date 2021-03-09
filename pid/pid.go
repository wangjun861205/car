package pid

import (
	"time"
)

type (
	cmdSetTarget struct {
		Synchronizer
		value float64
	}
	cmdGetTarget struct {
		Synchronizer
		value *float64
	}
	cmdGetDuty struct {
		Synchronizer
		value *float64
	}
	cmdGetMeasure struct {
		Synchronizer
		value *float64
	}
	cmdStopDuty struct {
		Synchronizer
	}
)

type status string

const (
	running status = "RUNNING"
	stopped status = "STOPPED"
)

// PID PID controller
type PID struct {
	kp       float64
	ki       float64
	kd       float64
	ticker   *time.Ticker
	target   float64
	accDiff  float64
	lastDiff float64
	dutier   Dutier
	measurer Measurer
	status   status
	commands chan Synchronizer
	close    chan interface{}
	done     chan interface{}
}

// NewPID NewPID
func NewPID(kp, ki, kd float64, cycle time.Duration, dutier Dutier, measurer Measurer) *PID {
	return &PID{
		kp:       kp,
		ki:       ki,
		kd:       kd,
		ticker:   time.NewTicker(cycle),
		target:   0,
		accDiff:  0,
		lastDiff: 0,
		dutier:   dutier,
		measurer: measurer,
		status:   stopped,
		commands: make(chan Synchronizer),
		close:    make(chan interface{}),
		done:     make(chan interface{}),
	}
}

// Run Run
func (p *PID) Run() {
	go p.measurer.Run()
	for {
		select {
		case <-p.close:
			p.dutier.Close()
			p.measurer.Close()
			close(p.done)
			return
		case cmd := <-p.commands:
			switch c := cmd.(type) {
			case cmdSetTarget:
				c.Ready()
				c.WaitActive()
				p.target = c.value
				p.status = running
				c.Done()
			case cmdGetTarget:
				c.Ready()
				c.WaitActive()
				*c.value = p.target
				c.Done()
			case cmdGetDuty:
				c.Ready()
				c.WaitActive()
				*c.value = p.dutier.GetDuty()
				c.Done()
			case cmdGetMeasure:
				c.Ready()
				c.WaitActive()
				*c.value = p.measurer.Measure()
				c.Done()
			case cmdStopDuty:
				c.Ready()
				c.WaitActive()
				p.dutier.Stop()
				p.accDiff, p.lastDiff = 0, 0
				p.status = stopped
				c.Done()
			}
		case <-p.ticker.C:
			if p.status == running {
				measure := p.measurer.Measure()
				diff := p.target - measure
				lastDiff := p.lastDiff
				p.lastDiff = diff
				accDiff := p.accDiff
				p.accDiff += diff
				duty := p.kp*diff + p.ki*accDiff + p.kd*(lastDiff-diff)
				p.dutier.SetDuty(duty)
			}
		}
	}
}

// SetTarget SetTarget
func (p *PID) SetTarget(sync Synchronizer, target float64) {
	p.commands <- cmdSetTarget{
		sync,
		target,
	}
}

// GetTarget GetTarget
func (p *PID) GetTarget(sync Synchronizer, target *float64) {
	p.commands <- cmdGetTarget{
		sync,
		target,
	}
}

// GetDuty GetDuty
func (p *PID) GetDuty(sync Synchronizer, duty *float64) {
	p.commands <- cmdGetDuty{
		sync,
		duty,
	}
}

// GetMeasure GetMeasure
func (p *PID) GetMeasure(sync Synchronizer, measure *float64) {
	p.commands <- cmdGetMeasure{
		sync,
		measure,
	}
}

// Stop Stop
func (p *PID) Stop(sync Synchronizer) {
	p.commands <- cmdStopDuty{
		sync,
	}
}

// Close Close
func (p *PID) Close() {
	close(p.close)
	<-p.done
}
