package pico

import (
	"fmt"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	EventSource
	v map[EventField]value
}

type value struct {
	v interface{}
	u EventUnit
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func New(source EventSource) *event {
	return &event{
		EventSource: source,
		v:           make(map[EventField]value),
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e *event) Source() EventSource {
	return e.EventSource
}

func (e *event) Is(f EventField) bool {
	_, exists := e.v[f]
	return exists
}

func (e *event) Set(f EventField, u EventUnit, v interface{}) Event {
	e.v[f] = value{v, u}
	return e
}

func (e *event) Value(f EventField) (interface{}, EventUnit) {
	if v, exists := e.v[f]; exists {
		return v.v, v.u
	} else {
		return nil, 0
	}
}

func (e *event) Emit(C chan<- Event) error {
	for i := 0; i < 3; i++ {
		select {
		case C <- e:
			return nil
		default:
			// Wait
			time.Sleep(time.Duration(i) * time.Millisecond)
		}
	}
	return ErrTimeout
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e *event) String() string {
	str := "<event"
	for k, v := range e.v {
		str += fmt.Sprintf(" %v=%v", k, v)
	}
	return str + ">"
}

func (v value) String() string {
	return fmt.Sprint(v.v, v.u)
}
