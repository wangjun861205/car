package model

import "fmt"

// CarAction Action
type CarAction string

const (
	// ActionAccelerate Accelerate
	ActionAccelerate CarAction = "ACCELERATE"
	// ActionBrake Brake
	ActionBrake CarAction = "BRAKE"
	// ActionTurnLeft Turn Left
	ActionTurnLeft CarAction = "TURN_LEFT"
	// ActionTurnRight Turn Right
	ActionTurnRight CarAction = "TURN_RIGHT"
	// ActionStop stop
	ActionStop CarAction = "STOP"
)

// CarInstruction instruction
type CarInstruction struct {
	Action CarAction
}

// CarResponse response
type CarResponse struct {
	LeftBase        uint8
	RightBase       uint8
	LeftStep        uint8
	RightStep       uint8
	LeftProportion  uint8
	RightProportion uint8
	Error           error
}

func (r CarResponse) String() string {
	return fmt.Sprintf("Left Base: %d, Right Base: %d, Left Step: %d, Right Step: %d, Left Proportion: %d, Right Proportion: %d, Error: %s",
		r.LeftBase, r.RightBase, r.LeftStep, r.RightStep, r.LeftProportion, r.RightProportion, r.Error)
}
