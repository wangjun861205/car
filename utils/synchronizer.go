package utils

import "sync"

// Synchronizer single commander multiple actor
type Synchronizer struct {
	ready  *sync.WaitGroup
	active chan interface{}
	done   *sync.WaitGroup
}

// NewSynchronizer NewSynchronizer
func NewSynchronizer(numOfActor int) *Synchronizer {
	var ready, done sync.WaitGroup
	ready.Add(numOfActor)
	done.Add(numOfActor)
	return &Synchronizer{
		ready:  &ready,
		active: make(chan interface{}),
		done:   &done,
	}
}

// Ready Ready
func (s *Synchronizer) Ready() {
	s.ready.Done()
}

// WaitReady WaitReady
func (s *Synchronizer) WaitReady() {
	s.ready.Wait()
}

// Active Active
func (s *Synchronizer) Active() {
	close(s.active)
}

// WaitActive WaitActive
func (s *Synchronizer) WaitActive() {
	<-s.active
}

// Done Done
func (s *Synchronizer) Done() {
	s.done.Done()
}

// WaitDone WaitDone
func (s *Synchronizer) WaitDone() {
	s.done.Wait()
}
