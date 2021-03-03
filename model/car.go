package model

import "fmt"

// CarAction Action
type CarAction string

const (
	// Forward Accelerate
	Forward CarAction = "FORWARD"
	// Backward Brake
	Backward CarAction = "BACKWARD"
	// TurnLeft Turn Left
	TurnLeft CarAction = "TURN_LEFT"
	// TurnRight Turn Right
	TurnRight CarAction = "TURN_RIGHT"
	// Stop Stop
	Stop CarAction = "Stop"
)

// Request Request
type Request struct {
	Action CarAction
}

// Response response
type Response struct {
	LeftDirection  Direction
	RightDirection Direction
	LeftDuty       uint64
	RightDuty      uint64
	Error          error
}

func (r Response) String() string {
	return fmt.Sprintf("Left Direction: %s || Left Power: %d ========== Right Direction: %s || Right Power: %d\n",
		r.LeftDirection, r.LeftDuty, r.RightDirection, r.RightDuty)
}

// Config config
type Config struct {
	LeftAPin    uint8
	LeftBPin    uint8
	LeftPWMNum  uint8
	RightAPin   uint8
	RightBPin   uint8
	RightPWMNum uint8
	Period      uint64
	LeftSteps   []uint64
	RightSteps  []uint64
	Addr        string
}
