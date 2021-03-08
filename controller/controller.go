package controller

import (
	"buxiong/car/model"
	"buxiong/car/speedmeter"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Controller Controller
type Controller struct {
	speeds     []float64
	leftIndex  int
	rightIndex int
	left       driver
	right      driver
	leftMeter  speedMeter
	rightMeter speedMeter
	close      chan interface{}
	done       chan interface{}
}

// NewController NewController
func NewController(speeds []float64, left, right driver, leftAPin, leftBPin, rightAPin, rightBPin int) (*Controller, error) {
	leftMeter, err := speedmeter.NewSpeedMeter(leftAPin, leftBPin)
	if err != nil {
		return nil, errors.Wrap(err, "NewController failed")
	}
	rightMeter, err := speedmeter.NewSpeedMeter(rightAPin, rightBPin)
	if err != nil {
		return nil, errors.Wrap(err, "NewController failed")
	}
	numOfSpeed := len(speeds)
	rev := make([]float64, 0, len(speeds))
	for i := len(speeds) - 1; i >= 0; i-- {
		rev = append(rev, -speeds[i])
	}
	rev = append(rev, 0)
	speeds = append(rev, speeds...)
	return &Controller{
		speeds,
		numOfSpeed,
		numOfSpeed,
		left,
		right,
		leftMeter,
		rightMeter,
		make(chan interface{}),
		make(chan interface{}),
	}, nil
}

func (c *Controller) pid(target, speed float64, driver driver) {
	duty := driver.GetDuty()
	diff := target - speed
	power := float64(duty) + 0.02*diff
	driver.SetDuty(int64(power))
}

// Run Run
func (c *Controller) Run() {
	go c.leftMeter.Run()
	go c.rightMeter.Run()
	for {
		select {
		case <-c.close:
			c.left.Brake()
			c.right.Brake()
			c.left.Close()
			c.right.Close()
			c.leftMeter.Close()
			c.rightMeter.Close()
			close(c.done)
			return
		default:
			leftSpeed, rightSpeed := c.leftMeter.Speed(), c.rightMeter.Speed()
			fmt.Println(leftSpeed, rightSpeed)
			if leftTarget := c.speeds[c.leftIndex]; leftTarget == 0 {
				c.left.Brake()
			} else {
				c.pid(leftTarget, leftSpeed, c.left)
			}
			if rightTarget := c.speeds[c.rightIndex]; rightTarget == 0 {
				c.right.Brake()
			} else {
				c.pid(rightTarget, rightSpeed, c.right)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// Close Close
func (c *Controller) Close() error {
	close(c.close)
	return nil
}

func (c *Controller) leftStepForward() {
	if c.leftIndex != len(c.speeds)-1 {
		c.leftIndex++
	}
}

func (c *Controller) rightStepForward() {
	if c.rightIndex != len(c.speeds)-1 {
		c.rightIndex++
	}
}

func (c *Controller) leftStepBackward() {
	if c.leftIndex != 0 {
		c.leftIndex--
	}
}

func (c *Controller) rightStepBackward() {
	if c.rightIndex != 0 {
		c.rightIndex--
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

func (c *Controller) determineCarDirection() direction {
	li, ri := c.leftIndex-len(c.speeds)/2, c.rightIndex-len(c.speeds)/2
	sum := li + ri
	if sum == 0 {
		if li == 0 && ri == 0 {
			return stay
		}
		if li < 0 {
			return left
		}
		return right
	} else if sum < 0 {
		switch {
		case li < 0 && ri < 0:
			switch {
			case li == ri:
				return backward
			case li < ri:
				return backwardRight
			default:
				return backwardLeft
			}
		case li < 0 && ri >= 0:
			return backwardRight
		default:
			return backwardLeft
		}
	} else {
		switch {
		case li > 0 && ri > 0:
			switch {
			case li == ri:
				return forward
			case li < ri:
				return forwardLeft
			default:
				return forwardRight
			}
		case li > 0 && ri <= 0:
			return forwardRight
		default:
			return forwardLeft
		}
	}
}

// Forward Forward
func (c *Controller) Forward() {
	// direction := c.determineCarDirection()
	// switch direction {
	// case stay:
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case forward:
	// 	if c.leftIndex == len(c.speeds)-1 || c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case backward:
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case forwardLeft:
	// 	if c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case forwardRight:
	// 	if c.leftIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case backwardLeft:
	// 	if c.leftIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case backwardRight:
	// 	if c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case left:
	// 	if c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// case right:
	// 	if c.leftIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepForward()
	// }
	if c.leftIndex == len(c.speeds)-1 || c.rightIndex == len(c.speeds)-1 {
		return
	}
	c.leftIndex++
	c.rightIndex++
}

// Backward Backward
func (c *Controller) Backward() {
	// direction := c.determineCarDirection()
	// switch direction {
	// case stay, forward, forwardLeft, forwardRight:
	// 	c.leftStepBackward()
	// 	c.rightStepBackward()
	// case backward, backwardLeft, backwardRight:
	// 	if c.leftIndex == len(c.speeds)-1 || c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepBackward()
	// 	c.rightStepBackward()
	// case left:
	// 	if c.leftIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepBackward()
	// 	c.rightStepBackward()
	// case right:
	// 	if c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepBackward()
	// 	c.rightStepBackward()
	// }
	if c.leftIndex == 0 || c.rightIndex == 0 {
		return
	}
	c.leftIndex--
	c.rightIndex--
}

// TurnLeft TurnLeft
func (c *Controller) TurnLeft() {
	// direction := c.determineCarDirection()
	// switch direction {
	// case stay:
	// 	c.rightStepForward()
	// 	c.leftStepBackward()
	// case forward, forwardLeft:
	// 	if c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.rightStepForward()
	// 	c.leftStepBackward()
	// case forwardRight:
	// 	c.rightStepForward()
	// 	c.leftStepBackward()
	// case backward, backwardLeft:
	// 	if c.rightIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepBackward()
	// case backwardRight:
	// 	c.leftStepForward()
	// 	c.rightStepBackward()
	// case left, right:
	// 	c.leftStepBackward()
	// 	c.rightStepForward()

	// }
	if c.leftIndex == 0 || c.rightIndex == len(c.speeds)-1 {
		return
	}
	c.leftIndex--
	c.rightIndex++
}

// TurnRight TurnRight
func (c *Controller) TurnRight() {
	// direction := c.determineCarDirection()
	// switch direction {
	// case stay, left, right, forwardLeft:
	// 	c.leftStepForward()
	// 	c.rightStepBackward()
	// case forward, forwardRight:
	// 	if c.leftIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepForward()
	// 	c.rightStepBackward()
	// case backward, backwardRight:
	// 	if c.leftIndex == len(c.speeds)-1 {
	// 		return
	// 	}
	// 	c.leftStepBackward()
	// 	c.rightStepForward()
	// case backwardLeft:
	// 	c.leftStepBackward()
	// 	c.rightStepForward()
	// }
	if c.leftIndex == len(c.speeds)-1 || c.rightIndex == 0 {
		return
	}
	c.leftIndex++
	c.rightIndex--
}

// Stop Stop
func (c *Controller) Stop() {
	c.leftIndex = len(c.speeds) / 2
	c.rightIndex = len(c.speeds) / 2
	c.left.Brake()
	c.right.Brake()
}

// Status Status
func (c *Controller) Status() *model.ControllerStatus {
	leftDuty := c.left.GetDuty()
	rightDuty := c.right.GetDuty()
	var leftDir, rightDir model.Direction
	switch {
	case leftDuty == 0:
		leftDir = model.DirectionBrake
	case leftDuty > 0:
		leftDir = model.DirectionForward
	default:
		leftDir = model.DirectionBackward
	}
	switch {
	case rightDuty == 0:
		rightDir = model.DirectionBrake
	case rightDuty > 0:
		rightDir = model.DirectionForward
	default:
		rightDir = model.DirectionBackward
	}
	return &model.ControllerStatus{
		LeftDirection:  leftDir,
		RightDirection: rightDir,
		LeftTarget:     c.speeds[c.leftIndex],
		RightTarget:    c.speeds[c.rightIndex],
	}
}
