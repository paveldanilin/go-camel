package observer

import (
	"fmt"
	"sync"
)

type State any

type Observer func(State)

type Subject struct {
	mu    sync.RWMutex
	slots []Observer
	free  int // free slots
}

func NewSubject() *Subject { return &Subject{} }

// Subscribe registers subscriber and returns function for unsubscribing.
func (s *Subject) Subscribe(ob Observer) (unsubscribe func()) {
	if ob == nil {
		return func() {}
	}

	s.mu.Lock()
	idx := -1
	if s.free > 0 {
		for i := range s.slots {
			if s.slots[i] == nil {
				s.slots[i] = ob
				idx = i
				s.free--
				break
			}
		}
	}
	if idx == -1 {
		s.slots = append(s.slots, ob)
		idx = len(s.slots) - 1
	}
	s.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			s.mu.Lock()
			if idx >= 0 && idx < len(s.slots) && s.slots[idx] != nil {
				s.slots[idx] = nil
				s.free++
				// Clean up tail
				for s.free > 0 && len(s.slots) > 0 {
					last := len(s.slots) - 1
					if s.slots[last] != nil {
						break
					}
					s.slots = s.slots[:last]
					s.free--
				}
			}
			s.mu.Unlock()
		})
	}
}

// Notify makes sync notification of all subscribers about new state.
// Does not handle panic.
func (s *Subject) Notify(st State) {
	s.mu.RLock()
	if len(s.slots) == 0 {
		s.mu.RUnlock()
		return
	}

	activeObservers := s.makeActiveObserversSnapshot()
	s.mu.RUnlock()

	for _, ob := range activeObservers {
		ob(st)
	}
}

// NotifySafe makes sync notification of all subscribers about new state.
// Handles panic and returns first error.
func (s *Subject) NotifySafe(st State) (err error) {
	s.mu.RLock()
	if len(s.slots) == 0 {
		s.mu.RUnlock()
		return nil
	}

	activeObservers := s.makeActiveObserversSnapshot()
	s.mu.RUnlock()

	// fail fast
	for i, ob := range activeObservers {

		func(idx int, fn Observer) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("observer panic at index %d: %v", idx, r)
				}
			}()
			if err != nil {
				return
			}
			fn(st)
		}(i, ob)

		if err != nil {
			return err
		}
	}
	return nil
}

// makeActiveObserversSnapshot returns a list of active observers.
func (s *Subject) makeActiveObserversSnapshot() []Observer {
	list := make([]Observer, 0, len(s.slots))
	for _, ob := range s.slots {
		if ob != nil {
			list = append(list, ob)
		}
	}
	return list
}
