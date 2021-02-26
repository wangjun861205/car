package model

import "fmt"

// ControllerStatus ControllerStatus
type ControllerStatus struct {
	left  *DriverStatus
	right *DriverStatus
}

// NewControllerStatus NewControllerStatus
func NewControllerStatus(left, right *DriverStatus) *ControllerStatus {
	return &ControllerStatus{
		left,
		right,
	}
}

// String String
func (s ControllerStatus) String() string {
	return fmt.Sprintf("left: || direction: %s || power(%03d%%) || =========  right: || direction: %s || power(%03d%%) ||\n",
		s.right.direction,
		s.left.duty,
		s.right.direction,
		s.right.duty)
}

// Left left driver status
func (s ControllerStatus) Left() *DriverStatus {
	return s.left
}

// Right right driver status
func (s ControllerStatus) Right() *DriverStatus {
	return s.right
}
