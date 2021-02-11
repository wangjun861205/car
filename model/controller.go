package model

import "fmt"

type ControllerStatus struct {
	Left  PWMStatus
	Right PWMStatus
}

func (s ControllerStatus) String() string {
	return fmt.Sprintf(`left: base(%d%%) step(%d%%) power(%d%%) right: base(%d%%) step(%d%%) power(%d%%)
	`, s.Left.Base, s.Left.Step, s.Left.Proportion, s.Right.Base, s.Right.Step, s.Right.Proportion)
}
