package core

import (
	"testing"

	e "bitbucket.org/enroute-mobi/ara/core/apierrs"
	"github.com/stretchr/testify/assert"
)

func createTestPartnerTemplateManager() *PartnerTemplateManager {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.collectManager = NewTestCollectManager()
	referentials.Save(referential)

	return NewPartnerTemplateManager(referential)
}

func TestFindByCredential(t *testing.T) {
	assert := assert.New(t)

	m := createTestPartnerTemplateManager()
	pt := m.New("test")
	pt.CredentialType = FormatMatching

	okTests := map[string][]string{
		"%{value}":          []string{"1234", "abcneizp", "ab_A12-32"},
		"test:%{value}:LOC": []string{"test:1234:LOC", "test:test-abcd:LOC"},
	}

	for k, s := range okTests {
		pt.LocalCredential = k
		pt.CompileLocalCredentials()
		m.Save(pt)

		for _, v := range s {
			assert.NotNil(m.FindByCredential(v))
		}
	}

	nokTests := map[string][]string{
		"%{value}":          []string{"12\\34", "&abcneizp"},
		"test:%{value}:LOC": []string{"nok:test:1234:LOC", "test:1234:LOC123"},
	}

	for k, s := range nokTests {
		pt.LocalCredential = k
		pt.CompileLocalCredentials()
		m.Save(pt)

		for _, v := range s {
			assert.Nil(m.FindByCredential(v))
		}
	}
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	m := createTestPartnerTemplateManager()
	pt := m.New("test")
	pt.CredentialType = FormatMatching
	pt.LocalCredential = "%{value}"
	pt.RemoteCredential = "%{value}"
	pt.Save()

	var TestCases = []struct {
		pt          *PartnerTemplate
		errLen      int
		errs        map[string][]string
		settingErrs map[string][]string
	}{
		{
			pt: &PartnerTemplate{
				manager:        m,
				Settings:       make(map[string]string),
				ConnectorTypes: []string{},
			},
			errs: map[string][]string{
				"CredentialType":   []string{e.ERROR_BLANK},
				"LocalCredential":  []string{e.ERROR_BLANK},
				"RemoteCredential": []string{e.ERROR_BLANK},
				"Slug":             []string{e.ERROR_BLANK},
			},
			settingErrs: map[string][]string{},
		},
		{
			pt: &PartnerTemplate{
				manager:          m,
				Settings:         make(map[string]string),
				ConnectorTypes:   []string{},
				CredentialType:   "wrong",
				LocalCredential:  "wrong",
				RemoteCredential: "wrong",
				Slug:             "WRONG",
			},
			errs: map[string][]string{
				"CredentialType":   []string{e.ERROR_FORMAT},
				"LocalCredential":  []string{e.ERROR_FORMAT},
				"RemoteCredential": []string{e.ERROR_FORMAT},
				"Slug":             []string{e.ERROR_SLUG_FORMAT},
			},
			settingErrs: map[string][]string{},
		},
		{
			pt: &PartnerTemplate{
				manager:          m,
				Settings:         make(map[string]string),
				ConnectorTypes:   []string{},
				CredentialType:   FormatMatching,
				LocalCredential:  "%{value}",
				RemoteCredential: "%{value}",
				Slug:             "test",
			},
			errs: map[string][]string{
				"LocalCredential": []string{e.ERROR_UNIQUE},
				"Slug":            []string{e.ERROR_UNIQUE},
			},
			settingErrs: map[string][]string{},
		},
		{
			pt: &PartnerTemplate{
				manager:          m,
				Settings:         make(map[string]string),
				ConnectorTypes:   []string{},
				CredentialType:   FormatMatching,
				LocalCredential:  "test:%{value}",
				RemoteCredential: "%{value}",
				Slug:             "test2",
			},
			errs:        map[string][]string{},
			settingErrs: map[string][]string{},
		},
		{
			pt: &PartnerTemplate{
				manager:  m,
				Settings: make(map[string]string),
				ConnectorTypes: []string{
					SIRI_CHECK_STATUS_CLIENT_TYPE,
					SIRI_CHECK_STATUS_SERVER_TYPE,
					SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR,
					SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER,
				},
				CredentialType:   FormatMatching,
				LocalCredential:  "test:%{value}",
				RemoteCredential: "%{value}",
				Slug:             "test2",
			},
			errs: map[string][]string{},
			settingErrs: map[string][]string{
				"remote_code_space": []string{e.ERROR_BLANK},
			},
		},
		{
			pt: &PartnerTemplate{
				manager: m,
				Settings: map[string]string{
					"remote_code_space": "test",
					"remote_url":        "test",
				},
				ConnectorTypes: []string{
					SIRI_CHECK_STATUS_CLIENT_TYPE,
					SIRI_CHECK_STATUS_SERVER_TYPE,
					SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR,
					SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER,
				},
				CredentialType:   FormatMatching,
				LocalCredential:  "test:%{value}",
				RemoteCredential: "%{value}",
				Slug:             "test2",
			},
			errs:        map[string][]string{},
			settingErrs: map[string][]string{},
		},
	}

	for _, tc := range TestCases {
		tc.pt.Validate()
		errLen := len(tc.errs)
		if len(tc.settingErrs) != 0 {
			errLen += 1
		}
		assert.Len(tc.pt.Errors, errLen)
		for k, _ := range tc.errs {
			assert.Equal(tc.errs[k], tc.pt.Errors.Get(k))
		}
		for k, _ := range tc.settingErrs {
			assert.Equal(tc.settingErrs[k], tc.pt.Errors.GetSettingError(k))
		}
	}
}
