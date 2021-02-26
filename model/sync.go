package model

import "sync"

// SyncGroup SyncGroup
type SyncGroup struct {
	received sync.WaitGroup
	active   chan struct{}
	done     sync.WaitGroup
}

// NewSyncGroup NewSyncGroup
func NewSyncGroup(numOfActor int) *SyncGroup {
	g := &SyncGroup{
		active: make(chan struct{}),
	}
	g.received.Add(numOfActor)
	g.done.Add(numOfActor)
	return g
}

// NotifyReady notify ready to active
func (s *SyncGroup) NotifyReady() {
	s.received.Done()
}

// WaitActive wait to active
func (s *SyncGroup) WaitActive() {
	<-s.active
}

// WaitReady wait ready notification
func (s *SyncGroup) WaitReady() {
	s.received.Wait()
}

// Active notify to active
func (s *SyncGroup) Active() {
	close(s.active)
}

// WaitDone wait finish acting
func (s *SyncGroup) WaitDone() {
	s.done.Wait()
}

// Done done
func (s *SyncGroup) Done() {
	s.done.Done()
}
