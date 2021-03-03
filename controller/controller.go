package controller

import (
	"buxiong/car/model"
)

type controller struct {
	leftSteps  []uint64
	rightSteps []uint64
	leftIndex  int
	rightIndex int
	left       driver
	right      driver
}

func NewController(leftSteps, rightSteps []uint64, left, right driver) *controller {
	return &controller{
		leftSteps,
		rightSteps,
		0,
		0,
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
		if c.leftIndex == len(c.leftSteps)-1 {
			return
		}
		c.leftIndex++
		c.left.Forward(c.leftSteps[c.leftIndex])
	case model.DirectionBackward:
		if c.leftIndex <= 1 {
			c.leftIndex = 0
			c.left.Brake()
			return
		}
		c.leftIndex--
		c.left.Backward(c.leftSteps[c.leftIndex])
	default:
		c.leftIndex = 1
		c.left.Forward(c.leftSteps[1])
	}
}

func (c *controller) rightStepForward() {
	stat := c.right.Status()
	switch stat.Direction() {
	case model.DirectionForward:
		if c.rightIndex == len(c.rightSteps)-1 {
			return
		}
		c.rightIndex++
		c.right.Forward(c.rightSteps[c.rightIndex])
	case model.DirectionBackward:
		if c.rightIndex <= 1 {
			c.rightIndex = 0
			c.right.Brake()
			return
		}
		c.rightIndex--
		c.right.Backward(c.rightSteps[c.rightIndex])
	default:
		c.rightIndex = 1
		c.right.Forward(c.rightSteps[1])
	}
}

func (c *controller) leftStepBackward() {
	stat := c.left.Status()
	switch stat.Direction() {
	case model.DirectionBackward:
		if c.leftIndex == len(c.leftSteps)-1 {
			return
		}
		c.leftIndex++
		c.left.Backward(c.leftSteps[c.leftIndex])
	case model.DirectionForward:
		if c.leftIndex <= 1 {
			c.leftIndex = 0
			c.left.Brake()
			return
		}
		c.leftIndex--
		c.left.Forward(c.leftSteps[c.leftIndex])
	default:
		c.leftIndex = 1
		c.left.Backward(c.leftSteps[1])
	}
}

func (c *controller) rightStepBackward() {
	stat := c.right.Status()
	switch stat.Direction() {
	case model.DirectionBackward:
		if c.rightIndex == len(c.rightSteps)-1 {
			return
		}
		c.rightIndex++
		c.right.Backward(c.rightSteps[c.rightIndex])
	case model.DirectionForward:
		if c.rightIndex <= 1 {
			c.rightIndex = 0
			c.right.Brake()
			return
		}
		c.rightIndex--
		c.right.Forward(c.rightSteps[c.rightIndex])
	default:
		c.rightIndex = 1
		c.right.Backward(c.rightSteps[1])
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

func determineCarDirection(leftStat, rightStat *model.DriverStatus, leftIndex, rightIndex int) direction {
	switch {
	case leftStat.Direction() == model.DirectionForward && rightStat.Direction() == model.DirectionForward:
		switch {
		case leftIndex == rightIndex:
			return forward
		case leftIndex < rightIndex:
			return forwardLeft
		default:
			return forwardRight
		}
	case leftStat.Direction() == model.DirectionForward && (rightStat.Direction() == model.DirectionBrake || rightStat.Direction() == model.DirectionGlide):
		return forwardRight
	case leftStat.Direction() == model.DirectionForward && rightStat.Direction() == model.DirectionBackward:
		switch {
		case leftIndex == rightIndex:
			return right
		case leftIndex < rightIndex:
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
		case leftIndex == rightIndex:
			return left
		case leftIndex < rightIndex:
			return forwardLeft
		default:
			return backwardRight
		}
	case leftStat.Direction() == model.DirectionBackward && (rightStat.Direction() == model.DirectionBrake || rightStat.Direction() == model.DirectionGlide):
		return backwardRight
	default:
		switch {
		case leftIndex == rightIndex:
			return backward
		case leftIndex < rightIndex:
			return backwardLeft
		default:
			return backwardRight
		}
	}
}

func (c *controller) Forward() {
	leftStat, rightStat := c.left.Status(), c.right.Status()
	direction := determineCarDirection(leftStat, rightStat, c.leftIndex, c.rightIndex)
	switch direction {
	case stay, backward, backwardLeft, backwardRight:
		c.leftStepForward()
		c.rightStepForward()
	case forward, forwardLeft, forwardRight:
		if c.leftIndex == len(c.leftSteps)-1 || c.rightIndex == len(c.rightSteps)-1 {
			return
		}
		c.leftStepForward()
		c.rightStepForward()
	case left:
		if c.rightIndex == len(c.rightSteps)-1 {
			return
		}
		c.leftStepForward()
		c.rightStepForward()
	case right:
		if c.leftIndex == len(c.leftSteps)-1 {
			return
		}
		c.leftStepForward()
		c.rightStepForward()
	}
}

func (c *controller) Backward() {
	leftStat, rightStat := c.left.Status(), c.right.Status()
	direction := determineCarDirection(leftStat, rightStat, c.leftIndex, c.rightIndex)
	switch direction {
	case stay, forward, forwardLeft, forwardRight:
		c.leftStepBackward()
		c.rightStepBackward()
	case backward, backwardLeft, backwardRight:
		if c.leftIndex == len(c.leftSteps)-1 || c.rightIndex == len(c.rightSteps)-1 {
			return
		}
		c.leftStepBackward()
		c.rightStepBackward()
	case left:
		if c.leftIndex == len(c.leftSteps)-1 {
			return
		}
		c.leftStepBackward()
		c.rightStepBackward()
	case right:
		if c.rightIndex == len(c.rightSteps)-1 {
			return
		}
		c.leftStepBackward()
		c.rightStepBackward()
	}
}

func (c *controller) TurnLeft() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	direction := determineCarDirection(leftStat, rightStat, c.leftIndex, c.rightIndex)
	switch direction {
	case stay:
		c.rightStepForward()
		c.leftStepBackward()
	case forward, forwardLeft:
		if c.rightIndex == len(c.rightSteps)-1 {
			return
		}
		c.rightStepForward()
		c.leftStepBackward()
	case forwardRight:
		c.rightStepForward()
		c.leftStepBackward()
	case backward, backwardLeft:
		if c.rightIndex == len(c.rightSteps)-1 {
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
	direction := determineCarDirection(leftStat, rightStat, c.leftIndex, c.rightIndex)
	switch direction {
	case stay, left, right, forwardLeft:
		c.leftStepForward()
		c.rightStepBackward()
	case forward, forwardRight:
		if c.leftIndex == len(c.leftSteps)-1 {
			return
		}
		c.leftStepForward()
		c.rightStepBackward()
	case backward, backwardRight:
		if c.leftIndex == len(c.leftSteps)-1 {
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
