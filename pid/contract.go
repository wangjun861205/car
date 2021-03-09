package pid

// Dutier Dutier
type Dutier interface {
	SetDuty(float64)
	GetDuty() float64
	Stop()
	Close()
}

// Measurer Measurer
type Measurer interface {
	Measure() float64
	Run()
	Close()
}

// Synchronizer Synchronizer
type Synchronizer interface {
	Ready()
	WaitActive()
	Done()
}
