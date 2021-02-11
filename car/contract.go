package car

import (
	"buxiong/car/model"
)

type controller interface {
	Run()
	Stop()
	Accelerate()
	Brake()
	TurnLeft()
	TurnRight()
	Status() model.ControllerStatus
}

type server interface {
	Run()
	Stop()
	RegisterHandler(func(b []byte) []byte)
}
