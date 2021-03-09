package controller

import "buxiong/car/pid"

type driver interface {
	Close() error
	Brake()
	SetDuty(duty int64)
	GetDuty() int64
}

type speedMeter interface {
	Run()
	Close()
	Speed() float64
}

// PID PID controller
type PID interface {
	Run()
	Close()
	SetTarget(sync pid.Synchronizer, target float64)
	GetDuty(sync pid.Synchronizer, target *float64)
	Stop(sync pid.Synchronizer)
}
