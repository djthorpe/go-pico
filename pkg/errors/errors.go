package errors

type Error uint

const (
	ErrSuccess Error = iota
	ErrBadParameter
	ErrUnexpectedValue
)

func (e Error) Error() string {
	switch e {
	case ErrBadParameter:
		return "ErrBadParameter"
	case ErrUnexpectedValue:
		return "ErrUnexpectedValue"
	default:
		return "Undefined error"
	}
}
