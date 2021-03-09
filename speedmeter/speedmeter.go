package speedmeter

import (
	"time"

	"github.com/pkg/errors"
	"github.com/warthog618/gpiod"
)

type eventType int

const (
	aRising eventType = iota
	aFalling
	bRising
	bFalling
)

type direction int8

const (
	clockwise        direction = 1
	counterClockwise direction = -1
)

type event struct {
	typ       eventType
	timestamp int64
}

type pin int

const (
	none pin = 0
	a    pin = 1
	b    pin = 2
)

// SpeedMeter motor balancer
type SpeedMeter struct {
	chip          *gpiod.Chip
	aLine         *gpiod.Line
	bLine         *gpiod.Line
	aLevel        uint8
	bLevel        uint8
	aRisingTime   int64
	aFallingTime  int64
	bRisingTime   int64
	bFallingTime  int64
	prevActivePin pin
	direction     direction
	speed         float64
	events        chan event
	close         chan interface{}
	done          chan interface{}
}

func genHandler(typ eventType, channal chan event) func(gpiod.LineEvent) {
	return func(e gpiod.LineEvent) {
		channal <- event{
			typ,
			time.Now().UnixNano(),
		}
	}
}

// NewSpeedMeter NewSpeedMeter
func NewSpeedMeter(aPin, bPin int) (*SpeedMeter, error) {
	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return nil, errors.Wrap(err, "NewSpeedMeter failed")
	}
	events := make(chan event)
	aLine, err := chip.RequestLine(aPin, gpiod.WithEventHandler(func(e gpiod.LineEvent) {
		if e.Type == gpiod.LineEventRisingEdge {
			events <- event{
				aRising,
				time.Now().UnixNano(),
			}
		} else {
			events <- event{
				aFalling,
				time.Now().UnixNano(),
			}
		}
	}), gpiod.WithBothEdges, gpiod.WithMonotonicEventClock)
	if err != nil {
		return nil, errors.Wrapf(err, "NewSpeedMeter failed(a line: %d)", aPin)
	}
	bLine, err := chip.RequestLine(bPin, gpiod.WithEventHandler(func(e gpiod.LineEvent) {
		if e.Type == gpiod.LineEventRisingEdge {
			events <- event{
				bRising,
				time.Now().UnixNano(),
			}
		} else {
			events <- event{
				bFalling,
				time.Now().UnixNano(),
			}
		}
	}), gpiod.WithBothEdges, gpiod.WithMonotonicEventClock)
	if err != nil {
		return nil, errors.Wrapf(err, "NewSpeedMeter failed(b line: %d)", bPin)
	}
	return &SpeedMeter{
		chip,
		aLine,
		bLine,
		0,
		0,
		0,
		0,
		0,
		0,
		none,
		0,
		0,
		events,
		make(chan interface{}),
		make(chan interface{}),
	}, nil
}

// Run Run
func (s *SpeedMeter) Run() {
OUTER:
	for {
		select {
		case <-s.close:
			s.aLine.Close()
			s.bLine.Close()
			s.chip.Close()
			close(s.done)
			return
		case e := <-s.events:
			switch e.typ {
			case aRising:
				if e.timestamp < s.aRisingTime {
					continue OUTER
				}
				if s.aLevel == 0 {
					if s.aFallingTime != 0 {
						s.speed = (1000 * 1000 * 1000) / (float64(e.timestamp) - float64(s.aFallingTime)) / 390 * 60 / 2
					}
					if s.prevActivePin == b || s.prevActivePin == none {
						if s.bLevel == 0 {
							s.direction = counterClockwise
						} else {
							s.direction = clockwise
						}
					}
				}
				s.prevActivePin = a
				s.aLevel = 1
				s.aRisingTime = e.timestamp
			case aFalling:
				if e.timestamp < s.aFallingTime {
					continue OUTER
				}
				if s.aLevel == 1 {
					if s.aRisingTime != 0 {
						s.speed = (1000 * 1000 * 1000) / (float64(e.timestamp) - float64(s.aRisingTime)) / 390 * 60 / 2
					}
					if s.prevActivePin == b || s.prevActivePin == none {
						if s.bLevel == 1 {
							s.direction = counterClockwise
						} else {
							s.direction = clockwise
						}
					}
				}
				s.prevActivePin = a
				s.aLevel = 0
				s.aFallingTime = e.timestamp
			case bRising:
				if e.timestamp < s.bRisingTime {
					continue OUTER
				}
				if s.bLevel == 0 {
					if s.bRisingTime != 0 {
						s.speed = (1000 * 1000 * 1000) / (float64(e.timestamp) - float64(s.bRisingTime)) / 390 * 60 / 2
					}
					if s.prevActivePin == a || s.prevActivePin == none {
						if s.aLevel == 1 {
							s.direction = counterClockwise
						} else {
							s.direction = clockwise
						}
					}
				}
				s.prevActivePin = b
				s.bLevel = 1
				s.bRisingTime = e.timestamp
			case bFalling:
				if e.timestamp < s.bFallingTime {
					continue OUTER
				}
				if s.bLevel == 1 {
					if s.bFallingTime != 0 {
						s.speed = (1000 * 1000 * 1000) / (float64(e.timestamp) - float64(s.bFallingTime)) / 390 * 60 / 2
					}
					if s.prevActivePin == a || s.prevActivePin == none {
						if s.aLevel == 0 {
							s.direction = counterClockwise
						} else {
							s.direction = clockwise
						}
					}
				}
				s.prevActivePin = b
				s.bLevel = 0
				s.bFallingTime = e.timestamp
			}
		}

	}
}

// Close Close
func (s *SpeedMeter) Close() {
	close(s.close)
	// <-s.done
}

// Measure motor speed
func (s *SpeedMeter) Measure() float64 {
	if time.Now().UnixNano()-s.aRisingTime > int64(100*time.Millisecond) {
		return 0
	}
	return s.speed * float64(s.direction)
}
