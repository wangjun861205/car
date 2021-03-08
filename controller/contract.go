package controller

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
