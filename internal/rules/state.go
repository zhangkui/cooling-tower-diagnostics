package rules

import (
	"fmt"
	"time"
)

type State struct {
	Name      string
	EnteredAt time.Time
	Reason    string
}
type StateMachine struct {
	Current     State
	Transitions map[string]map[string]bool
}

func NewStateMachine(initial string, now time.Time) *StateMachine {
	return &StateMachine{Current: State{Name: initial, EnteredAt: now}, Transitions: map[string]map[string]bool{}}
}
func (s *StateMachine) Allow(from, to string) {
	if s.Transitions[from] == nil {
		s.Transitions[from] = map[string]bool{}
	}
	s.Transitions[from][to] = true
}
func (s *StateMachine) Move(to, reason string, now time.Time) error {
	if s.Current.Name == to {
		return fmt.Errorf("already in %s", to)
	}
	if allowed := s.Transitions[s.Current.Name][to]; !allowed {
		return fmt.Errorf("transition %s -> %s denied", s.Current.Name, to)
	}
	s.Current = State{Name: to, EnteredAt: now, Reason: reason}
	return nil
}
func (s *StateMachine) Age(now time.Time) time.Duration { return now.Sub(s.Current.EnteredAt) }
func (s *StateMachine) In(name string) bool             { return s.Current.Name == name }
func (s *StateMachine) Reset(name string, now time.Time) {
	s.Current = State{Name: name, EnteredAt: now}
}
