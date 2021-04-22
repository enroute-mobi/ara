package remote

// To handle Partner Status after a Gtfs request
type GtfsError struct {
	message string
}

func (e GtfsError) Error() string {
	return e.message
}

func NewGtfsError(message string) error {
	return &GtfsError{message: message}
}
