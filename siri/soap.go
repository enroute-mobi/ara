package siri

import "strings"

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

// Temp
func WrapSoap(s string) string {
	soap := strings.Join([]string{
		"<?xml version='1.0' encoding='utf-8'?>",
		"<S:Envelope xmlns:S=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\">",
		"<S:Body>",
		s,
		"</S:Body>",
		"</S:Envelope>"}, "\n")
	return soap
}
