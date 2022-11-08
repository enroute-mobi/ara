package core

import (
	"encoding/json"
	"fmt"

	e "bitbucket.org/enroute-mobi/ara/core/apierrs"
	ps "bitbucket.org/enroute-mobi/ara/core/psettings"
)

type APIPartner struct {
	Id             PartnerId `json:"Id,omitempty"`
	Slug           PartnerSlug
	Name           string            `json:"Name,omitempty"`
	Settings       map[string]string `json:"Settings,omitempty"`
	ConnectorTypes []string          `json:"ConnectorTypes,omitempty"`
	Errors         e.Errors          `json:"Errors,omitempty"`

	factories map[string]ConnectorFactory
	manager   Partners
}

func (partner *APIPartner) Validate() bool {
	partner.Errors = e.NewErrors()

	// Check if slug is non null
	if partner.Slug == "" {
		partner.Errors.Add("Slug", e.ERROR_BLANK)
	} else if !slugRegexp.MatchString(string(partner.Slug)) { // slugRegexp defined in Referential
		partner.Errors.Add("Slug", e.ERROR_SLUG_FORMAT)
	}

	// Check factories
	partner.setFactories()
	for _, factory := range partner.factories {
		factory.Validate(partner)
	}

	// Check Slug uniqueness
	for _, existingPartner := range partner.manager.FindAll() {
		if existingPartner.id != partner.Id && existingPartner.slug == partner.Slug {
			partner.Errors.Add("Slug", e.ERROR_UNIQUE)
		}
	}

	// Check Credentials uniqueness
	if !partner.manager.UniqCredentials(partner.Id, partner.credentials()) {
		if _, ok := partner.Settings[ps.LOCAL_CREDENTIAL]; ok {
			partner.Errors.AddSettingError(ps.LOCAL_CREDENTIAL, e.ERROR_UNIQUE)
		}
		if _, ok := partner.Settings[ps.LOCAL_CREDENTIALS]; ok {
			partner.Errors.AddSettingError(ps.LOCAL_CREDENTIALS, e.ERROR_UNIQUE)
		}
	}

	return len(partner.Errors) == 0
}

func (partner *APIPartner) credentials() string {
	return fmt.Sprintf("%v,%v", partner.Settings[ps.LOCAL_CREDENTIAL], partner.Settings[ps.LOCAL_CREDENTIALS])
}

func (partner *APIPartner) setFactories() {
	for _, connectorType := range partner.ConnectorTypes {
		factory := NewConnectorFactory(connectorType)
		if factory == nil {
			partner.Errors.AddConnectorTypesError(connectorType)
			continue
		}
		partner.factories[connectorType] = factory
	}
}

func (partner *APIPartner) IsSettingDefined(setting string) (ok bool) {
	_, ok = partner.Settings[setting]
	return
}

func (partner *APIPartner) ValidatePresenceOfSetting(setting string) bool {
	if !partner.IsSettingDefined(setting) {
		partner.Errors.AddSettingError(setting, e.ERROR_BLANK)
		return false
	}
	return true
}

func (partner *APIPartner) ValidatePresenceOfLocalCredentials() bool {
	if !partner.IsSettingDefined(ps.LOCAL_CREDENTIAL) && !partner.IsSettingDefined(ps.LOCAL_CREDENTIALS) {
		partner.Errors.AddSettingError(ps.LOCAL_CREDENTIAL, e.ERROR_BLANK)
		return false
	}
	return true
}

func (partner *APIPartner) ValidatePresenceOfRemoteObjectIdKind() bool {
	return partner.ValidatePresenceOfSetting(ps.REMOTE_OBJECTID_KIND)
}

func (partner *APIPartner) ValidatePresenceOfRemoteCredentials() bool {
	return partner.ValidatePresenceOfSetting(ps.REMOTE_URL) && partner.ValidatePresenceOfSetting(ps.REMOTE_CREDENTIAL)
}

func (partner *APIPartner) ValidatePresenceOfLightRemoteCredentials() bool {
	return partner.ValidatePresenceOfSetting(ps.REMOTE_URL)
}

func (partner *APIPartner) ValidatePresenceOfConnector(connector string) bool {
	for _, listedConnector := range partner.ConnectorTypes {
		if listedConnector == connector {
			return true
		}
	}
	partner.Errors.Add(fmt.Sprintf("Connector %s", connector), e.ERROR_BLANK)
	return false
}

func (partner *APIPartner) UnmarshalJSON(data []byte) error {
	type Alias APIPartner
	aux := &struct {
		Settings map[string]string
		*Alias
	}{
		Alias: (*Alias)(partner),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Settings != nil {
		partner.Settings = aux.Settings
	}

	return nil
}
