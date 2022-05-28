package errors

type Error uint

const (
	ErrSuccess Error = iota
	ErrBadParameter
	ErrUnexpectedValue
	ErrSampleSkipped
	ErrTimeout
)

func (e Error) Error() string {
	switch e {
	case ErrBadParameter:
		return "ErrBadParameter"
	case ErrUnexpectedValue:
		return "ErrUnexpectedValue"
	case ErrSampleSkipped:
		return "ErrSampleSkipped"
	case ErrTimeout:
		return "ErrTimeout"
	default:
		return "Undefined error"
	}
}
