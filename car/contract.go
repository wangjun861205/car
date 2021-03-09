package car

import (
	"buxiong/car/model"
)

type controller interface {
	Run()
	Close() error
	Forward()
	Backward()
	TurnLeft()
	TurnRight()
	Brake()
	Status() *model.ControllerStatus
}

type server interface {
	Run()
	Requests() chan model.Request
	Errors() chan error
	Reply(model.Response) error
	Close() error
}
