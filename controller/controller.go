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

func (c *controller) leftStepForward() {
	stat := c.left.Status()
	switch stat.Direction() {
	case model.DirectionForward:
		if stat.Duty()+c.step > c.max {
			return
		}
		c.left.Forward(stat.Duty() + c.step)
	case model.DirectionBackward:
		if stat.Duty()-c.step < c.base {
			c.left.Brake()
			return
		}
		c.left.Backward(stat.Duty() - c.step)
	default:
		c.left.Forward(c.base)
	}
}

func (c *controller) rightStepForward() {
	stat := c.right.Status()
	switch stat.Direction() {
	case model.DirectionForward:
		if stat.Duty()+c.step > c.max {
			return
		}
		c.right.Forward(stat.Duty() + c.step)
	case model.DirectionBackward:
		if stat.Duty()-c.step < c.base {
			c.right.Brake()
			return
		}
		c.right.Backward(stat.Duty() - c.step)
	default:
		c.right.Forward(c.base)
	}
}

func (c *controller) leftStepBackward() {
	stat := c.left.Status()
	switch stat.Direction() {
	case model.DirectionBackward:
		if stat.Duty()+c.step > c.max {
			return
		}
		c.left.Backward(stat.Duty() + c.step)
	case model.DirectionForward:
		if stat.Duty()-c.step < c.base {
			c.left.Brake()
			return
		}
		c.left.Forward(stat.Duty() - c.step)
	default:
		c.left.Backward(c.base)
	}
}

func (c *controller) rightStepBackward() {
	stat := c.right.Status()
	switch stat.Direction() {
	case model.DirectionBackward:
		if stat.Duty()+c.step > c.max {
			return
		}
		c.right.Backward(stat.Duty() + c.step)
	case model.DirectionForward:
		if stat.Duty()-c.step < c.base {
			c.right.Brake()
			return
		}
		c.right.Forward(stat.Duty() - c.step)
	default:
		c.right.Backward(c.base)
	}
}

type direction string

const (
	forward       direction = "FORWARD"
	backward      direction = "BACKWARD"
	left          direction = "LEFT"
	right         direction = "RIGHT"
	forwardLeft   direction = "FORWARD_LEFT"
	forwardRight  direction = "FORWARD_RIGHT"
	backwardLeft  direction = "BACKWARD_LEFT"
	backwardRight direction = "BACKWARD_RIGHT"
	stay          direction = "STAY"
)

func determineCarDirection(leftStat, rightStat *model.DriverStatus) direction {
	switch {
	case leftStat.Direction() == model.DirectionForward && rightStat.Direction() == model.DirectionForward:
		switch {
		case leftStat.Duty() == rightStat.Duty():
			return forward
		case leftStat.Duty() < rightStat.Duty():
			return forwardLeft
		default:
			return forwardRight
		}
	case leftStat.Direction() == model.DirectionForward && (rightStat.Direction() == model.DirectionBrake || rightStat.Direction() == model.DirectionGlide):
		return forwardRight
	case leftStat.Direction() == model.DirectionForward && rightStat.Direction() == model.DirectionBackward:
		switch {
		case leftStat.Duty() == rightStat.Duty():
			return right
		case leftStat.Duty() < rightStat.Duty():
			return backwardLeft
		default:
			return forwardRight
		}
	case (leftStat.Direction() == model.DirectionBrake || leftStat.Direction() == model.DirectionGlide) && rightStat.Direction() == model.DirectionForward:
		return forwardLeft
	case (leftStat.Direction() == model.DirectionBrake || leftStat.Direction() == model.DirectionGlide) && (rightStat.Direction() == model.DirectionBrake || rightStat.Direction() == model.DirectionGlide):
		return stay
	case (leftStat.Direction() == model.DirectionBrake || leftStat.Direction() == model.DirectionGlide) && rightStat.Direction() == model.DirectionBackward:
		return backwardLeft
	case leftStat.Direction() == model.DirectionBackward && rightStat.Direction() == model.DirectionForward:
		switch {
		case leftStat.Duty() == rightStat.Duty():
			return left
		case leftStat.Duty() < rightStat.Duty():
			return forwardLeft
		default:
			return backwardRight
		}
	case leftStat.Direction() == model.DirectionBackward && (rightStat.Direction() == model.DirectionBrake || rightStat.Direction() == model.DirectionGlide):
		return backwardRight
	default:
		switch {
		case leftStat.Duty() == rightStat.Duty():
			return backward
		case leftStat.Duty() < rightStat.Duty():
			return backwardLeft
		default:
			return backwardRight
		}
	}
}

func (c *controller) Forward() {
	leftStat, rightStat := c.left.Status(), c.right.Status()
	direction := determineCarDirection(leftStat, rightStat)
	switch direction {
	case stay, backward, backwardLeft, backwardRight:
		c.leftStepForward()
		c.rightStepForward()
	case forward, forwardLeft, forwardRight:
		if leftStat.Duty()+c.step > c.max || rightStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepForward()
		c.rightStepForward()
	case left:
		if rightStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepForward()
		c.rightStepForward()
	case right:
		if leftStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepForward()
		c.rightStepForward()
	}
}

func (c *controller) Backward() {
	leftStat, rightStat := c.left.Status(), c.right.Status()
	direction := determineCarDirection(leftStat, rightStat)
	switch direction {
	case stay, forward, forwardLeft, forwardRight:
		c.leftStepBackward()
		c.rightStepBackward()
	case backward, backwardLeft, backwardRight:
		if leftStat.Duty()+c.step > c.max || rightStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepBackward()
		c.rightStepBackward()
	case left:
		if leftStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepBackward()
		c.rightStepBackward()
	case right:
		if rightStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepBackward()
		c.rightStepBackward()
	}
}

func (c *controller) TurnLeft() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	direction := determineCarDirection(leftStat, rightStat)
	switch direction {
	case stay:
		c.rightStepForward()
		c.leftStepBackward()
	case forward, forwardLeft:
		if rightStat.Duty()+c.step > c.max {
			return
		}
		c.rightStepForward()
		c.leftStepBackward()
	case forwardRight:
		c.rightStepForward()
		c.leftStepBackward()
	case backward, backwardLeft:
		if rightStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepForward()
		c.rightStepBackward()
	case backwardRight:
		c.leftStepForward()
		c.rightStepBackward()
	case left, right:
		c.leftStepBackward()
		c.rightStepForward()

	}
}

func (c *controller) TurnRight() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	direction := determineCarDirection(leftStat, rightStat)
	switch direction {
	case stay, left, right, forwardLeft:
		c.leftStepForward()
		c.rightStepBackward()
	case forward, forwardRight:
		if leftStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepForward()
		c.rightStepBackward()
	case backward, backwardRight:
		if leftStat.Duty()+c.step > c.max {
			return
		}
		c.leftStepBackward()
		c.rightStepForward()
	case backwardLeft:
		c.leftStepBackward()
		c.rightStepForward()
	}
}

func (c *controller) Stop() {
	c.left.Brake()
	c.right.Brake()
}

func (c *controller) Status() *model.ControllerStatus {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	return model.NewControllerStatus(leftStat, rightStat)
}
