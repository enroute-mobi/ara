package core

import (
	"encoding/json"
	"regexp"
	"slices"
	"strings"
	"sync"

	e "bitbucket.org/enroute-mobi/ara/core/apierrs"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

var CredentialTypes = []string{FormatMatching}

const (
	FormatMatching = "FormatMatching"

	valuePattern         = "%{value}"
	escaptedValuePattern = "%\\{value\\}"
	valuePatternReplace  = "([0-9a-zA-Z-_]+)"
)

type PartnerTemplates interface {
	SlugAndCredentialsHandler

	New(PartnerSlug) *PartnerTemplate
	Find(PartnerId) *PartnerTemplate
	FindBySlug(PartnerSlug) *PartnerTemplate
	FindByRawCredential(string, string) *PartnerTemplate
	FindByCredential(string) (*PartnerTemplate, string)
	FindAll() []*PartnerTemplate
	Save(*PartnerTemplate) bool
	Delete(*PartnerTemplate) bool
	Referential() *Referential
	IsEmpty() bool
}

type PartnerTemplate struct {
	manager PartnerTemplates

	id               PartnerId
	credentialRegexp *regexp.Regexp

	Slug             PartnerSlug
	Name             string `json:",omitempty"`
	CredentialType   string
	LocalCredential  string
	RemoteCredential string
	Settings         map[string]string
	ConnectorTypes   []string

	Errors e.Errors `json:",omitempty"`
}

type PartnerTemplateManager struct {
	uuid.UUIDConsumer

	mutex *sync.RWMutex

	byId        map[PartnerId]*PartnerTemplate
	referential *Referential
}

/* Partner template */

func (pt *PartnerTemplate) Save() {
	pt.manager.Save(pt)
}

func (pt *PartnerTemplate) Validate() bool {
	pt.Errors = e.NewErrors()

	if pt.CredentialType == "" {
		pt.Errors.Add("CredentialType", e.ERROR_BLANK)
	} else if !slices.Contains(CredentialTypes, pt.CredentialType) {
		pt.Errors.Add("CredentialType", e.ERROR_FORMAT)
	}

	if pt.LocalCredential == "" {
		pt.Errors.Add("LocalCredential", e.ERROR_BLANK)
	} else if !strings.Contains(pt.LocalCredential, valuePattern) {
		pt.Errors.Add("LocalCredential", e.ERROR_FORMAT)
	} else if !pt.UniqCredentials() {
		pt.Errors.Add("LocalCredential", e.ERROR_UNIQUE)
	} else {
		switch pt.CredentialType {
		case FormatMatching:
			if err := pt.CompileLocalCredentials(); err != nil {
				pt.Errors.Add("LocalCredential", e.ERROR_FORMAT)
			}
		}
	}

	if pt.RemoteCredential == "" {
		pt.Errors.Add("RemoteCredential", e.ERROR_BLANK)
	} else if !strings.Contains(pt.RemoteCredential, valuePattern) {
		pt.Errors.Add("RemoteCredential", e.ERROR_FORMAT)
	}

	validationSettings := make(map[string]string, len(pt.Settings))
	for k := range pt.Settings {
		validationSettings[k] = pt.Settings[k]
	}

	validationSettings[s.LOCAL_CREDENTIAL] = pt.LocalCredential
	validationSettings[s.REMOTE_CREDENTIAL] = pt.RemoteCredential
	validationSettings[s.REMOTE_URL] = "xxx"
	apiPartner := &APIPartner{
		Id:             pt.id,
		Slug:           pt.Slug,
		Name:           pt.Name,
		Settings:       validationSettings,
		ConnectorTypes: pt.ConnectorTypes,
		factories:      make(map[string]ConnectorFactory),
		Errors:         e.NewErrors(),
		manager:        pt.manager,
	}
	apiPartner.ValidateSlug()
	apiPartner.ValidateFactories()
	pt.Errors.Merge(apiPartner.Errors)

	return len(pt.Errors) == 0
}

func (pt *PartnerTemplate) CompileLocalCredentials() (err error) {
	pt.credentialRegexp, err = regexp.Compile("^" + strings.Replace(regexp.QuoteMeta(pt.LocalCredential), escaptedValuePattern, valuePatternReplace, -1) + "$")
	return
}

func (pt *PartnerTemplate) BuildRemoteCredential(m string) string {
	return strings.Replace(pt.RemoteCredential, valuePattern, m, -1)
}

func (pt *PartnerTemplate) UniqCredentials() bool {
	return pt.manager.UniqCredentials(pt.id, pt.LocalCredential, pt.CredentialType)
}

func (pt *PartnerTemplate) Copy() (copy *PartnerTemplate) {
	return &PartnerTemplate{
		id:               pt.id,
		Slug:             pt.Slug,
		manager:          pt.manager,
		Name:             pt.Name,
		CredentialType:   pt.CredentialType,
		LocalCredential:  pt.LocalCredential,
		RemoteCredential: pt.RemoteCredential,
		Settings:         pt.Settings,
		ConnectorTypes:   pt.ConnectorTypes,
	}
}

func (pt *PartnerTemplate) MarshalJSON() ([]byte, error) {
	type Alias PartnerTemplate
	return json.Marshal(&struct {
		Id PartnerId
		*Alias
	}{
		Id:    pt.id,
		Alias: (*Alias)(pt),
	})
}

/* Partner template manager */

func NewPartnerTemplateManager(referential *Referential) *PartnerTemplateManager {
	manager := &PartnerTemplateManager{
		mutex:       &sync.RWMutex{},
		byId:        make(map[PartnerId]*PartnerTemplate),
		referential: referential,
	}
	return manager
}

func (manager *PartnerTemplateManager) UniqSlug(id PartnerId, s PartnerSlug) bool {
	manager.mutex.RLock()
	for _, pt := range manager.byId {
		if pt.Slug == s && pt.id != id {
			manager.mutex.RUnlock()
			return false
		}
	}

	manager.mutex.RUnlock()
	return true

}

func (manager *PartnerTemplateManager) UniqCredentials(id PartnerId, c string, params ...string) bool {
	pt := manager.FindByRawCredential(c, params[0]) // params = CredentialType

	if pt != nil && pt.id != id {
		return false
	}
	return true
}

func (manager *PartnerTemplateManager) New(slug PartnerSlug) *PartnerTemplate {
	return &PartnerTemplate{
		Slug:           slug,
		manager:        manager,
		Settings:       make(map[string]string),
		ConnectorTypes: []string{},
	}
}

func (manager *PartnerTemplateManager) MarshalJSON() ([]byte, error) {
	pts := []PartnerTemplate{}
	for _, pt := range manager.byId {
		pts = append(pts, *pt)
	}
	return json.Marshal(pts)
}

func (manager *PartnerTemplateManager) Find(id PartnerId) (partner *PartnerTemplate) {
	manager.mutex.RLock()
	partner = manager.byId[id]
	manager.mutex.RUnlock()
	return
}

func (manager *PartnerTemplateManager) FindByRawCredential(c, ct string) (pt *PartnerTemplate) {
	manager.mutex.RLock()
	for k := range manager.byId {
		if manager.byId[k].LocalCredential == c && manager.byId[k].CredentialType == ct {
			pt = manager.byId[k]
			manager.mutex.RUnlock()
			return
		}
	}

	manager.mutex.RUnlock()
	return
}

func (manager *PartnerTemplateManager) FindByCredential(c string) (pt *PartnerTemplate, match string) {
	manager.mutex.RLock()
	for k := range manager.byId {
		switch manager.byId[k].CredentialType {
		case FormatMatching:
			if r := manager.byId[k].credentialRegexp.FindStringSubmatch(c); r != nil {
				pt = manager.byId[k]
				manager.mutex.RUnlock()
				match = r[1]
				return
			}
		}
	}

	manager.mutex.RUnlock()
	return
}

func (manager *PartnerTemplateManager) FindBySlug(slug PartnerSlug) (pt *PartnerTemplate) {
	manager.mutex.RLock()
	for _, t := range manager.byId {
		if t.Slug == slug {
			pt = t
			manager.mutex.RUnlock()
			return
		}
	}

	manager.mutex.RUnlock()
	return
}

func (manager *PartnerTemplateManager) FindAll() (pts []*PartnerTemplate) {
	manager.mutex.RLock()
	for _, pt := range manager.byId {
		pts = append(pts, pt)
	}
	manager.mutex.RUnlock()
	return
}

// Warning: PartnerTemplate.Validate() must be called for the regexp to be compiled for now
func (manager *PartnerTemplateManager) Save(pt *PartnerTemplate) bool {
	if pt.id == "" {
		pt.id = PartnerId(manager.NewUUID())
	}
	pt.manager = manager

	manager.mutex.Lock()
	manager.byId[pt.id] = pt
	manager.mutex.Unlock()

	return true
}

func (manager *PartnerTemplateManager) Delete(pt *PartnerTemplate) bool {
	manager.referential.Partners().DeleteFromTemplate(pt.id)

	manager.mutex.Lock()
	delete(manager.byId, pt.id)
	manager.mutex.Unlock()

	return true
}

func (manager *PartnerTemplateManager) Referential() *Referential {
	return manager.referential
}

func (manager *PartnerTemplateManager) IsEmpty() bool {
	return len(manager.byId) == 0
}
