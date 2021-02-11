package model

// PWMActionType PWM action type
type PWMActionType string

const (
	// PWMActionTypeStatus status
	PWMActionTypeStatus PWMActionType = "PWM_STATUS"
	// PWMActionTypeIncrement increment proportion
	PWMActionTypeIncrement PWMActionType = "PWM_INCREMENT"
	// PWMActionTypeDecrement decrement proportion
	PWMActionTypeDecrement PWMActionType = "PWM_DECREMENT"
)

// PWMAction  PWM action
type PWMAction struct {
	Type PWMActionType
}

// type pwmInstruction string

// const (
// 	pwmInstructionHigh pwmInstruction = "HIGH"
// 	pwmInstructionLow  pwmInstruction = "LOW"
// )

// type pwmer interface {
// 	run()
// 	stop()
// 	status() pwmStatus
// 	incProportion()
// 	decProportion()
// }

// PWMStatus PWM status
type PWMStatus struct {
	Base       uint8
	Step       uint8
	Proportion uint8
}
