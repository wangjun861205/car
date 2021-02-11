package controller

import (
	"buxiong/car/model"
)

type pwmer interface {
	Run()
	Stop()
	Increment()
	Decrement()
	Status() model.PWMStatus
}
