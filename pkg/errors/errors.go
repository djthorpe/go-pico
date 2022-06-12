package errors

import "fmt"

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Error uint

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ErrSuccess Error = iota
	ErrBadParameter
	ErrUnexpectedValue
	ErrDuplicateValue
	ErrSampleSkipped
	ErrTimeout
	ErrNotImplemented
	ErrNotInitialised
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Error) Error() string {
	switch e {
	case ErrBadParameter:
		return "ErrBadParameter"
	case ErrUnexpectedValue:
		return "ErrUnexpectedValue"
	case ErrDuplicateValue:
		return "ErrDuplicateValue"
	case ErrSampleSkipped:
		return "ErrSampleSkipped"
	case ErrTimeout:
		return "ErrTimeout"
	case ErrNotImplemented:
		return "ErrNotImplemented"
	case ErrNotInitialised:
		return "ErrNotInitialised"
	default:
		return "Undefined error"
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e Error) With(args ...interface{}) error {
	return fmt.Errorf("%w: %s", e, fmt.Sprint(args...))
}

func (e Error) Withf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", e, fmt.Sprintf(format, args...))
}
