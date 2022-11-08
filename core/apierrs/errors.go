package apierrs

import "fmt"

type Errors map[string]interface{}

func NewErrors() Errors {
	return make(Errors)
}

func (errors Errors) Empty() bool {
	return len(errors) == 0
}

const (
	SETTINGS          = "Settings"
	CONNECTOR_TYPES   = "ConnectorTypes"
	ERROR_BLANK       = "Can't be empty"
	ERROR_SLUG_FORMAT = "Invalid format: only lowercase alphanumeric characters and _"
	ERROR_ZERO        = "Can't be zero"
	ERROR_UNIQUE      = "Is already in use"
)

func (errors Errors) Get(attribute string) []string {
	m, ok := errors[attribute]
	if !ok {
		return []string{}
	}
	s, ok := m.([]string)
	if !ok {
		return []string{}
	}
	return s
}

func (errors Errors) Add(attribute, message string) {
	if errors.added(attribute, message) {
		return
	}
	errors[attribute] = append(errors[attribute].([]string), message)
}

func (errors Errors) added(attribute, message string) bool {
	i, ok := errors[attribute]
	if !ok {
		errors[attribute] = []string{}
		return false
	}

	messageArray, ok := i.([]string)
	if !ok {
		return true // Shouldn't happen, but if it does we want to do nothing
	}
	for _, m := range messageArray {
		if m == message {
			return true
		}
	}

	return false
}

func (errors Errors) GetSettingError(attribute string) []string {
	m, ok := errors[SETTINGS]
	if !ok {
		return []string{}
	}
	s, ok := m.(Errors)
	if !ok {
		return []string{}
	}
	return s.Get(attribute)
}

func (errors Errors) AddSettingError(attribute, message string) {
	se, ok := errors[SETTINGS]
	if !ok {
		se = NewErrors()
		errors[SETTINGS] = se
	}
	se.(Errors).Add(attribute, message)
}

func (errors Errors) GetConnectorTypesError(connector string) []string {
	m, ok := errors[CONNECTOR_TYPES]
	if !ok {
		return []string{}
	}
	s, ok := m.(Errors)
	if !ok {
		return []string{}
	}
	return s.Get(connector)
}

func (errors Errors) AddConnectorTypesError(connector string) {
	cte, ok := errors[CONNECTOR_TYPES]
	if !ok {
		cte = NewErrors()
		errors[CONNECTOR_TYPES] = cte
	}
	cte.(Errors).Add(connector, fmt.Sprintf("Unknown connector %v", connector))
}

// Test method
func (errors Errors) GetSettings() Errors {
	s, ok := errors[SETTINGS]
	if ok {
		return s.(Errors)
	}
	return NewErrors()
}
