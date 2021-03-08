package model

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
	LeftTarget     float64
	RightTarget    float64
	Error          error
}

// Config config
type Config struct {
	LeftAPin         uint8
	LeftBPin         uint8
	LeftPWMNum       uint8
	RightAPin        uint8
	RightBPin        uint8
	RightPWMNum      uint8
	LeftEncoderAPin  uint8
	LeftEncoderBPin  uint8
	RightEncoderAPin uint8
	RightEncoderBPin uint8
	Speeds           []float64
	Addr             string
}
