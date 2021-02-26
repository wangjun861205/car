package controller

import (
	"buxiong/car/model"
)

type driver interface {
	Close() error
	Brake()
	Glide()
	Forward(duty uint64)
	Backward(duty uint64)
	Status() *model.DriverStatus
}
