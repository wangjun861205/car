package controller

import (
	"buxiong/car/model"
)

type controller struct {
	max   uint64
	base  uint64
	step  uint64
	left  driver
	right driver
}

func NewController(max, base, step uint64, left, right driver) *controller {
	return &controller{
		max,
		base,
		step,
		left,
		right,
	}
}

func (c *controller) Close() error {
	if err := c.left.Close(); err != nil {
		return err
	}
	return c.right.Close()
}

func (c *controller) forward(status *model.DriverStatus, driver driver) {
	switch status.Direction() {
	case model.DirectionForward:
		if status.Duty()+c.step > c.max {
			driver.Forward(c.max)
		} else {
			driver.Forward(status.Duty() + c.step)
		}
	case model.DirectionBackward:
		if status.Duty()-c.step < c.base {
			driver.Brake()
		} else {
			driver.Backward(status.Duty() - c.step)
		}
	default:
		driver.Forward(c.base)
	}
}

func (c *controller) Forward() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	c.forward(leftStat, c.left)
	c.forward(rightStat, c.right)
}

func (c *controller) backward(status *model.DriverStatus, driver driver) {
	switch status.Direction() {
	case model.DirectionForward:
		if status.Duty()-c.step < c.base {
			driver.Brake()
		} else {
			driver.Forward(status.Duty() - c.step)
		}
	case model.DirectionBackward:
		if status.Duty()+c.step > c.max {
			driver.Backward(c.max)
		} else {
			driver.Backward(status.Duty() + c.step)
		}
	default:
		driver.Backward(c.base)
	}
}

func (c *controller) Backward() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	c.backward(leftStat, c.left)
	c.backward(rightStat, c.right)
}

func (c *controller) TurnLeft() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	if leftStat.Duty()+c.step > c.max || rightStat.Duty()+c.step > c.max {
		return
	}
	c.backward(leftStat, c.left)
	c.forward(rightStat, c.right)
}

func (c *controller) TurnRight() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	if leftStat.Duty()+c.step > c.max || rightStat.Duty()+c.step > c.max {
		return
	}
	c.forward(leftStat, c.left)
	c.backward(rightStat, c.right)
}

func (c *controller) Status() *model.ControllerStatus {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	return model.NewControllerStatus(leftStat, rightStat)
}
