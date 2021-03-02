package balancer

import (
	"github.com/pkg/errors"
	"github.com/warthog618/gpiod"
)

const numOfSample = 39000

// Balancer motor balancer
type Balancer struct {
	aPin      int
	bPin      int
	aLevel    uint8
	bLevel    uint8
	direction uint8
	speed     uint64
	isReset   bool
	samples   []int64
	close     chan interface{}
	done      chan interface{}
}

// NewBalancer NewBalancer
func NewBalancer(aPin, bPin int) *Balancer {
	return &Balancer{
		aPin,
		bPin,
		0,
		0,
		0,
		0,
		false,
		make([]int64, 0, numOfSample),
		make(chan interface{}),
		make(chan interface{}),
	}
}

// Run Run
func (b *Balancer) Run() error {
	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return errors.Wrap(err, "Balancer.Run() failed")
	}
	aline, err := chip.RequestLine(b.aPin, gpiod.WithEventHandler(func(event gpiod.LineEvent) {
		switch event.Type {
		case gpiod.LineEventRisingEdge:
			if b.isReset {
				return
			}
			b.aLevel = 1
			b.direction = b.aLevel & b.bLevel
			if len(b.samples) < numOfSample {
				b.samples = append(b.samples, event.Timestamp.Microseconds())
			} else {
				b.samples = append(b.samples[1:], event.Timestamp.Microseconds())
			}
			if len(b.samples) > 1 {
				b.speed = uint64(len(b.samples)) * 1000 * 1000 * 60 / uint64(b.samples[len(b.samples)-1]-b.samples[0]) / 390
			}
		case gpiod.LineEventFallingEdge:
			if b.isReset {
				return
			}
			b.aLevel = 0
		}
	}), gpiod.WithBothEdges, gpiod.WithMonotonicEventClock)
	if err != nil {
		return errors.Wrap(err, "Balancer.Run() failed")
	}
	bline, err := chip.RequestLine(b.bPin, gpiod.WithEventHandler(func(event gpiod.LineEvent) {
		switch event.Type {
		case gpiod.LineEventRisingEdge:
			if b.isReset {
				return
			}
			b.bLevel = 1
			b.direction = b.aLevel ^ b.bLevel
		case gpiod.LineEventFallingEdge:
			if b.isReset {
				return
			}
			b.bLevel = 0
		}
	}), gpiod.WithBothEdges, gpiod.WithMonotonicEventClock)
	if err != nil {
		return errors.Wrap(err, "Balancer.Run() failed")
	}
	<-b.close
	aline.Close()
	bline.Close()
	chip.Close()
	close(b.done)
	return nil
}

// Close Close
func (b *Balancer) Close() {
	close(b.close)
	<-b.done
}

// Reset Reset
func (b *Balancer) Reset() {
	b.isReset = true
	b.aLevel = 0
	b.bLevel = 0
	b.direction = 0
	b.speed = 0
	b.samples = make([]int64, 0, 390)
	b.isReset = false
}

// GetSpeed GetSpeed
func (b *Balancer) GetSpeed() uint64 {
	return b.speed
}
