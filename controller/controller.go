package controller

import "buxiong/car/model"

type controller struct {
	left  pwmer
	right pwmer
}

func NewController(left, right pwmer) *controller {
	return &controller{
		left,
		right,
	}
}

func (c *controller) Run() {
	c.left.Run()
	c.right.Run()
}

func (c *controller) Stop() {
	c.left.Stop()
	c.right.Stop()
}

func (c *controller) Accelerate() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	if leftStat.Proportion+leftStat.Step > 100 || rightStat.Proportion+rightStat.Step > 100 {
		return
	}
	c.left.Increment()
	c.right.Increment()
}

func (c *controller) Brake() {
	c.left.Decrement()
	c.right.Decrement()
}

func (c *controller) TurnLeft() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	if leftStat.Proportion <= rightStat.Proportion {
		c.left.Decrement()
	} else {
		c.right.Increment()
	}
}

func (c *controller) TurnRight() {
	leftStat := c.left.Status()
	rightStat := c.right.Status()
	if rightStat.Proportion <= leftStat.Proportion {
		c.right.Decrement()
	} else {
		c.left.Increment()
	}
}

func (c *controller) Status() model.ControllerStatus {
	return model.ControllerStatus{
		Left:  c.left.Status(),
		Right: c.right.Status(),
	}
}
