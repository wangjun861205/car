package car

import (
	"buxiong/car/model"
)

type controller interface {
	Close() error
	Forward()
	Backward()
	TurnLeft()
	TurnRight()
	Stop()
	Status() *model.ControllerStatus
}

type server interface {
	Run()
	Requests() chan model.Request
	Errors() chan error
	Reply(model.Response) error
	Close() error
}
