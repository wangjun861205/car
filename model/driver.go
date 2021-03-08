package model

// InstructionType Driver action type
type InstructionType string

const (
	// MoveForward forward
	MoveForward InstructionType = "MOVE_FORWARD"
	// MoveBackward backward
	MoveBackward InstructionType = "MOVE_BACKWARD"
	// Brake brake
	Brake InstructionType = "BRAKE"
	// Glide glide
	Glide InstructionType = "GLIDE"
	// Close close
	Close InstructionType = "CLOSE"
	// GetStatus get status
	GetStatus InstructionType = "GET_STATUS"
)

// Direction pwm direction
type Direction string

const (
	// DirectionBrake brake
	DirectionBrake Direction = "BRAKE"
	// DirectionGlide glide
	DirectionGlide Direction = "GLIDE"
	// DirectionForward forward
	DirectionForward Direction = "FORWARD"
	// DirectionBackward backward
	DirectionBackward Direction = "BACKWARD"
)

// DriverInstruction  Driver action
type DriverInstruction struct {
	*SyncGroup
	typ    InstructionType
	duty   uint64
	status *DriverStatus
}

// NewDriverInstruction NewDriverInstruction
func NewDriverInstruction(typ InstructionType, duty uint64, syncGroup *SyncGroup) *DriverInstruction {
	return &DriverInstruction{
		SyncGroup: syncGroup,
		typ:       typ,
		duty:      duty,
		status:    nil,
	}
}

// Type return action type
func (p *DriverInstruction) Type() InstructionType {
	return p.typ
}

// Duty return duty
func (p *DriverInstruction) Duty() uint64 {
	return p.duty
}

// GetStatus get status
func (p *DriverInstruction) GetStatus() *DriverStatus {
	return p.status
}

// SetStatus put status
func (p *DriverInstruction) SetStatus(status *DriverStatus) {
	p.status = status
}

// DriverStatus DriverStatus
type DriverStatus string

const (
	// Outputting Outputting
	Outputting DriverStatus = "OUTPUTTING"
	// Braking Braking
	Braking DriverStatus = "BRAKING"
	// Gliding Gliding
	Gliding DriverStatus = "GLIDING"
)
