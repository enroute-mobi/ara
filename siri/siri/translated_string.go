package siri

import (
	"bytes"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type SIRITranslatedString struct {
	Tag    string
	Prefix string
	model.TranslatedString
}

func NewSIRISXTranslatedString() *SIRITranslatedString {
	return &SIRITranslatedString{
		Prefix: "siri:",
	}
}

func (t *SIRITranslatedString) BuildTranslatedStringXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "translated_string.template", t); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}
