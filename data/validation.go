package data

type ValidationError struct {
	msg    string
	Fields map[string]string
}

func NewValidationError(msg string, fields map[string]string) *ValidationError {
	return &ValidationError{
		msg,
		fields,
	}
}

func (e *ValidationError) Error() string {
	return e.msg
}
