package siri

import "fmt"

// Handle SIRI CRITICAL errors
type SiriError struct {
	errCode string
	message string
}

func (e *SiriError) ErrCode() string {
	return e.errCode
}

func (e *SiriError) Error() string {
	return e.message
}

func (e *SiriError) FullMessage() string {
	return fmt.Sprintf("%v: %v", e.errCode, e.message)
}

func NewSiriError(message string) error {
	return &SiriError{message: message}
}

func NewSiriErrorWithCode(errCode, message string) error {
	return &SiriError{
		errCode: errCode,
		message: message,
	}
}
