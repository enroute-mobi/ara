package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"

	e "bitbucket.org/enroute-mobi/ara/core/apierrs"
	"bitbucket.org/enroute-mobi/ara/core/partners"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
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
	partners.SlugAndCredentialsHandler

	New(partners.Slug) *PartnerTemplate
	Find(partners.Id) *PartnerTemplate
	FindBySlug(partners.Slug) *PartnerTemplate
	FindByRawCredential(string, string) *PartnerTemplate
	FindByCredential(string) (*PartnerTemplate, string)
	FindAll() []*PartnerTemplate
	Save(*PartnerTemplate) bool
	Delete(*PartnerTemplate) bool
	Referential() *Referential
	IsEmpty() bool
	SaveToDatabase() (int, error)
}

type PartnerTemplate struct {
	manager PartnerTemplates

	id               partners.Id
	credentialRegexp *regexp.Regexp

	Slug             partners.Slug
	Name             string `json:",omitempty"`
	CredentialType   string
	LocalCredential  string
	RemoteCredential string
	MaxPartners      int
	Settings         map[string]string
	ConnectorTypes   []string

	Errors e.Errors `json:",omitempty"`
}

type PartnerTemplateManager struct {
	uuid.UUIDConsumer

	mutex *sync.RWMutex

	byId        map[partners.Id]*PartnerTemplate
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
		Id partners.Id
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
		byId:        make(map[partners.Id]*PartnerTemplate),
		referential: referential,
	}
	return manager
}

func (manager *PartnerTemplateManager) UniqSlug(id partners.Id, s partners.Slug) bool {
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

func (manager *PartnerTemplateManager) UniqCredentials(id partners.Id, c string, params ...string) bool {
	pt := manager.FindByRawCredential(c, params[0]) // params = CredentialType

	if pt != nil && pt.id != id {
		return false
	}
	return true
}

func (manager *PartnerTemplateManager) New(slug partners.Slug) *PartnerTemplate {
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

func (manager *PartnerTemplateManager) Find(id partners.Id) (partner *PartnerTemplate) {
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

func (manager *PartnerTemplateManager) FindBySlug(slug partners.Slug) (pt *PartnerTemplate) {
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
		pt.id = partners.Id(manager.NewUUID())
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

func (manager *PartnerTemplateManager) SaveToDatabase() (int, error) {
	// Check presence of Referential
	selectReferentials := []model.SelectReferential{}
	sqlQuery := fmt.Sprintf("select * from referentials where referential_id = '%s'", manager.referential.Id())
	_, err := model.Database.Select(&selectReferentials, sqlQuery)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}
	if len(selectReferentials) == 0 {
		return http.StatusNotAcceptable, errors.New("can't save Partner templates without Referential in Database")
	}

	// Begin transaction
	tx, err := model.Database.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	// Delete partner templates
	sqlQuery = fmt.Sprintf("delete from partner_templates where referential_id = '%s';", manager.referential.Id())
	_, err = tx.Exec(sqlQuery)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	// Insert partner templates
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	for _, pt := range manager.byId {
		dbPT, err := manager.newDbPartnerTemplate(pt)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, fmt.Errorf("internal error: %v", err)
		}
		err = tx.Insert(dbPT)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, fmt.Errorf("internal error: %v", err)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	return http.StatusOK, nil
}

func (manager *PartnerTemplateManager) newDbPartnerTemplate(pt *PartnerTemplate) (*model.DatabasePartnerTemplate, error) {
	settings, err := json.Marshal(pt.Settings)
	if err != nil {
		return nil, err
	}
	connectors, err := json.Marshal(pt.ConnectorTypes)
	if err != nil {
		return nil, err
	}
	return &model.DatabasePartnerTemplate{
		Id:               string(pt.id),
		ReferentialId:    string(manager.referential.id),
		Slug:             string(pt.Slug),
		Name:             pt.Name,
		CredentialType:   pt.CredentialType,
		LocalCredential:  pt.LocalCredential,
		RemoteCredential: pt.RemoteCredential,
		MaxPartners:      pt.MaxPartners,
		Settings:         string(settings),
		ConnectorTypes:   string(connectors),
	}, nil
}
