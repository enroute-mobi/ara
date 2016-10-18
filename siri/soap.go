package siri

// Handle SIRI CRITICAL errors
type SiriError struct {
	message string
}

func (e *SiriError) Error() string {
	return e.message
}

func newSiriError(message string) error {
	return &SiriError{message: message}
}
