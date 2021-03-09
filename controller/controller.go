package controller

import (
	"buxiong/car/model"
	"buxiong/car/utils"
)

type command string

const (
	moveForward  command = "MOVE_FORWARD"
	moveBackward command = "MOVE_BACKWARD"
	turnLeft     command = "TURN_LEFT"
	turnRight    command = "TURN_RIGHT"
	brake        command = "BRAKE"
)

// Controller Controller
type Controller struct {
	speeds     []float64
	leftIndex  int
	rightIndex int
	leftPID    PID
	rightPID   PID
	commands   chan command
	close      chan interface{}
	done       chan interface{}
}

// NewController NewController
func NewController(speeds []float64, leftPID, rightPID PID) (*Controller, error) {
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
		leftPID,
		rightPID,
		make(chan command),
		make(chan interface{}),
		make(chan interface{}),
	}, nil
}

func (c *Controller) forward() {
	speedLength := len(c.speeds)
	natureIndex := speedLength / 2
	topIndex := speedLength - 1
	switch {
	case c.leftIndex == natureIndex && c.rightIndex == natureIndex:
		c.leftIndex++
		c.rightIndex++
	case c.leftIndex+c.rightIndex-2*natureIndex > 0:
		if c.leftIndex != topIndex && c.rightIndex != topIndex {
			c.leftIndex++
			c.rightIndex++
		}
	case c.leftIndex+c.rightIndex-2*natureIndex < 0:
		if c.leftIndex != natureIndex && c.rightIndex != natureIndex {
			c.leftIndex++
			c.rightIndex++
		}
	}
	sync := utils.NewSynchronizer(2)
	c.leftPID.SetTarget(sync, c.speeds[c.leftIndex])
	c.rightPID.SetTarget(sync, c.speeds[c.rightIndex])
	sync.WaitReady()
	sync.Active()
	sync.WaitDone()
}

func (c *Controller) backward() {
	speedLength := len(c.speeds)
	natureIndex := speedLength / 2
	switch {
	case c.leftIndex == natureIndex && c.rightIndex == natureIndex:
		c.leftIndex--
		c.rightIndex--
	case c.leftIndex+c.rightIndex-2*natureIndex > 0:
		if c.leftIndex != natureIndex && c.rightIndex != natureIndex {
			c.leftIndex--
			c.rightIndex--
		}
	case c.leftIndex+c.rightIndex-2*natureIndex < 0:
		if c.leftIndex != 0 && c.rightIndex != 0 {
			c.leftIndex--
			c.rightIndex--
		}
	}
	sync := utils.NewSynchronizer(2)
	c.leftPID.SetTarget(sync, c.speeds[c.leftIndex])
	c.rightPID.SetTarget(sync, c.speeds[c.rightIndex])
	sync.WaitReady()
	sync.Active()
	sync.WaitDone()
}

func (c *Controller) turnLeft() {
	speedLength := len(c.speeds)
	natureIndex := speedLength / 2
	topIndex := speedLength - 1
	switch {
	case c.leftIndex+c.rightIndex-2*natureIndex == 0:
		if c.leftIndex != 0 && c.rightIndex != topIndex {
			c.leftIndex--
			c.rightIndex++
		}
	case c.leftIndex+c.rightIndex-2*natureIndex > 0:
		if c.leftIndex <= c.rightIndex {
			if c.leftIndex != natureIndex {
				c.leftIndex--
			}
		} else {
			c.rightIndex++
		}
	case c.leftIndex+c.rightIndex-2*natureIndex < 0:
		if c.leftIndex >= c.rightIndex && c.leftIndex != natureIndex {
			c.leftIndex++
		} else {
			c.rightIndex--
		}
	}
	sync := utils.NewSynchronizer(2)
	c.leftPID.SetTarget(sync, c.speeds[c.leftIndex])
	c.rightPID.SetTarget(sync, c.speeds[c.rightIndex])
	sync.WaitReady()
	sync.Active()
	sync.WaitDone()
}

func (c *Controller) turnRight() {
	speedLength := len(c.speeds)
	natureIndex := speedLength / 2
	topIndex := speedLength - 1
	switch {
	case c.leftIndex+c.rightIndex-2*natureIndex == 0:
		if c.leftIndex != topIndex && c.rightIndex != 0 {
			c.leftIndex++
			c.rightIndex--
		}
	case c.leftIndex+c.rightIndex-2*natureIndex > 0:
		if c.leftIndex >= c.rightIndex {
			if c.rightIndex != natureIndex {
				c.rightIndex--
			}
		} else {
			c.leftIndex++
		}
	case c.leftIndex+c.rightIndex-2*natureIndex < 0:
		if c.leftIndex <= c.rightIndex {
			if c.rightIndex != natureIndex {
				c.rightIndex++
			}
		} else {
			c.leftIndex--
		}
	}
	sync := utils.NewSynchronizer(2)
	c.leftPID.SetTarget(sync, c.speeds[c.leftIndex])
	c.rightPID.SetTarget(sync, c.speeds[c.rightIndex])
	sync.WaitReady()
	sync.Active()
	sync.WaitDone()
}

// Run Run
func (c *Controller) Run() {
	go c.leftPID.Run()
	go c.rightPID.Run()
	for {
		select {
		case <-c.close:
			c.leftPID.Close()
			c.rightPID.Close()
			close(c.done)
			return
		case command := <-c.commands:
			switch command {
			case moveForward:
				c.forward()
			case moveBackward:
				c.backward()
			case turnLeft:
				c.turnLeft()
			case turnRight:
				c.turnRight()
			case brake:
				c.leftIndex = len(c.speeds) / 2
				c.rightIndex = len(c.speeds) / 2
				sync := utils.NewSynchronizer(2)
				c.leftPID.Stop(sync)
				c.rightPID.Stop(sync)
				sync.WaitReady()
				sync.Active()
				sync.WaitDone()
			}
		}
	}
}

// Close Close
func (c *Controller) Close() error {
	close(c.close)
	return nil
}

// func (c *Controller) leftStepForward() {
// 	if c.leftIndex != len(c.speeds)-1 {
// 		c.leftIndex++
// 	}
// }

// func (c *Controller) rightStepForward() {
// 	if c.rightIndex != len(c.speeds)-1 {
// 		c.rightIndex++
// 	}
// }

// func (c *Controller) leftStepBackward() {
// 	if c.leftIndex != 0 {
// 		c.leftIndex--
// 	}
// }

// func (c *Controller) rightStepBackward() {
// 	if c.rightIndex != 0 {
// 		c.rightIndex--
// 	}
// }

// type direction string

// const (
// 	forward       direction = "FORWARD"
// 	backward      direction = "BACKWARD"
// 	left          direction = "LEFT"
// 	right         direction = "RIGHT"
// 	forwardLeft   direction = "FORWARD_LEFT"
// 	forwardRight  direction = "FORWARD_RIGHT"
// 	backwardLeft  direction = "BACKWARD_LEFT"
// 	backwardRight direction = "BACKWARD_RIGHT"
// 	stay          direction = "STAY"
// )

// func (c *Controller) determineCarDirection() direction {
// 	li, ri := c.leftIndex-len(c.speeds)/2, c.rightIndex-len(c.speeds)/2
// 	sum := li + ri
// 	if sum == 0 {
// 		if li == 0 && ri == 0 {
// 			return stay
// 		}
// 		if li < 0 {
// 			return left
// 		}
// 		return right
// 	} else if sum < 0 {
// 		switch {
// 		case li < 0 && ri < 0:
// 			switch {
// 			case li == ri:
// 				return backward
// 			case li < ri:
// 				return backwardRight
// 			default:
// 				return backwardLeft
// 			}
// 		case li < 0 && ri >= 0:
// 			return backwardRight
// 		default:
// 			return backwardLeft
// 		}
// 	} else {
// 		switch {
// 		case li > 0 && ri > 0:
// 			switch {
// 			case li == ri:
// 				return forward
// 			case li < ri:
// 				return forwardLeft
// 			default:
// 				return forwardRight
// 			}
// 		case li > 0 && ri <= 0:
// 			return forwardRight
// 		default:
// 			return forwardLeft
// 		}
// 	}
// }

// Forward Forward
func (c *Controller) Forward() {
	c.commands <- moveForward
}

// Backward Backward
func (c *Controller) Backward() {
	c.commands <- moveBackward
}

// TurnLeft TurnLeft
func (c *Controller) TurnLeft() {
	c.commands <- turnLeft
}

// TurnRight TurnRight
func (c *Controller) TurnRight() {
	c.commands <- turnRight
}

// Brake Brake
func (c *Controller) Brake() {
	c.commands <- brake
}

// Status Status
func (c *Controller) Status() *model.ControllerStatus {
	var leftDuty, rightDuty float64
	sync := utils.NewSynchronizer(2)
	c.leftPID.GetDuty(sync, &leftDuty)
	c.rightPID.GetDuty(sync, &rightDuty)
	sync.WaitReady()
	sync.Active()
	sync.WaitDone()
	var leftDir, rightDir model.Direction
	leftTarget := c.speeds[c.leftIndex]
	rightTarget := c.speeds[c.rightIndex]
	switch {
	case leftTarget == 0:
		leftDir = model.DirectionBrake
	case leftTarget > 0:
		leftDir = model.DirectionForward
	default:
		leftDir = model.DirectionBackward
	}
	switch {
	case rightTarget == 0:
		rightDir = model.DirectionBrake
	case rightTarget > 0:
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
