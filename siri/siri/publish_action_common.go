package siri

import (
	"bytes"
	"fmt"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type SIRIPublishActionCommon struct {
	HasAffects bool
	SIRIAffects

	Name               string                `json:",omitempty"`
	ActionType         string                `json:",omitempty"`
	Value              string                `json:",omitempty"`
	Prompt             *SIRITranslatedString `json:",omitempty"`
	HasPublishAtScope bool
	ScopeType          model.SituationScopeType `json:",omitempty"`
	ActionStatus       model.SituationActionStatus `json:",omitempty"`
	Description        *SIRITranslatedString       `json:",omitempty"`
	PublicationWindows []*model.TimeRange          `json:",omitempty"`
}

func (ac *SIRIPublishActionCommon) BuildPublishActionCommonXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("publish_action_common%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, ac); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}
