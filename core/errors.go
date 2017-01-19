package core

type Errors map[string][]string

func NewErrors() Errors {
	return make(Errors)
}

func (errors Errors) Empty() bool {
	return len(errors) == 0
}

const (
	ERROR_BLANK  = "Can't be empty"
	ERROR_ZERO   = "Can't be zero"
	ERROR_UNIQUE = "Is already in use"
)

func (errors Errors) Add(attribute string, message string) {
	errors[attribute] = append(errors[attribute], message)
}

func (errors Errors) Added(attribute string, message string) bool {
	messages, ok := errors[attribute]
	if !ok {
		return false
	}

	for _, m := range messages {
		if m == message {
			return true
		}
	}

	return false
}
