package model

import "database/sql"

type DatabaseStructureWithContext interface {
	GetHook() sql.NullString
	GetModelType() sql.NullString
	GetType() string
	GetAttributes() sql.NullString
}

type DatabaseReferential struct {
	ReferentialId  string         `db:"referential_id"`
	Slug           string         `db:"slug"`
	Name           string         `db:"name"`
	Settings       string         `db:"settings"`
	Tokens         string         `db:"tokens"`
	ImportTokens   string         `db:"import_tokens"`
	OrganisationId sql.NullString `db:"organisation_id"`
}

type SelectReferential struct {
	ReferentialId  string         `db:"referential_id"`
	OrganisationId sql.NullString `db:"organisation_id"`
	Slug           string
	Name           sql.NullString
	Settings       sql.NullString
	Tokens         sql.NullString
	ImportTokens   sql.NullString `db:"import_tokens"`
}

type DatabasePartner struct {
	Id             string `db:"id"`
	ReferentialId  string `db:"referential_id"`
	Slug           string `db:"slug"`
	Name           string `db:"name"`
	Settings       string `db:"settings"`
	ConnectorTypes string `db:"connector_types"`
}

type SelectPartner struct {
	Id             string
	ReferentialId  string `db:"referential_id"`
	Slug           string
	Name           sql.NullString
	Settings       sql.NullString
	ConnectorTypes sql.NullString `db:"connector_types"`
}

type DatabasePartnerTemplate struct {
	Id               string `db:"id"`
	ReferentialId    string `db:"referential_id"`
	Slug             string `db:"slug"`
	Name             string `db:"name"`
	Settings         string `db:"settings"`
	ConnectorTypes   string `db:"connector_types"`
	CredentialType   string `db:"credential_type"`
	LocalCredential  string `db:"local_credential"`
	RemoteCredential string `db:"remote_credential"`
	MaxPartners      int    `db:"max_partners"`
}

type SelectPartnerTemplate struct {
	Id               string
	ReferentialId    string `db:"referential_id"`
	Slug             string
	CredentialType   string        `db:"credential_type"`
	LocalCredential  string        `db:"local_credential"`
	RemoteCredential string        `db:"remote_credential"`
	MaxPartners      sql.NullInt64 `db:"max_partners"`
	Name             sql.NullString
	Settings         sql.NullString
	ConnectorTypes   sql.NullString `db:"connector_types"`
}

type DatabaseOperator struct {
	Id              string `db:"id"`
	ReferentialSlug string `db:"referential_slug"`
	Name            string `db:"name"`
	Codes           string `db:"codes"`
	ModelName       string `db:"model_name"`
}

type SelectOperator struct {
	Id              string
	ReferentialSlug string `db:"referential_slug"`
	ModelName       string `db:"model_name"`
	Name            sql.NullString
	Codes           sql.NullString `db:"codes"`
}

type DatabaseStopArea struct {
	LineIds           string         `db:"line_ids"`
	ReferentialSlug   string         `db:"referential_slug"`
	References        string         `db:"siri_references"`
	Attributes        string         `db:"attributes"`
	ModelName         string         `db:"model_name"`
	Name              string         `db:"name"`
	Codes             string         `db:"codes"`
	Id                string         `db:"id"`
	ReferentId        sql.NullString `db:"referent_id"`
	ParentId          sql.NullString `db:"parent_id"`
	CollectedAlways   bool           `db:"collected_always"`
	CollectChildren   bool           `db:"collect_children"`
	CollectSituations bool           `db:"collect_situations"`
}

type SelectStopArea struct {
	Id                string
	ReferentialSlug   string `db:"referential_slug"`
	ModelName         string `db:"model_name"`
	Name              sql.NullString
	Codes             sql.NullString `db:"codes"`
	ParentId          sql.NullString `db:"parent_id"`
	ReferentId        sql.NullString `db:"referent_id"`
	Attributes        sql.NullString
	References        sql.NullString `db:"siri_references"`
	LineIds           sql.NullString `db:"line_ids"`
	CollectedAlways   sql.NullBool   `db:"collected_always"`
	CollectChildren   sql.NullBool   `db:"collect_children"`
	CollectSituations sql.NullBool   `db:"collect_situations"`
}

type DatabaseStopAreaGroup struct {
	Id              string `db:"id"`
	ReferentialSlug string `db:"referential_slug"`
	ModelName       string `db:"model_name"`
	Name            string `db:"name"`
	ShortName       string `db:"short_name"`
	StopAreaIds     string `db:"stop_area_ids"`
}

type SelectStopAreaGroup struct {
	Id              string
	ReferentialSlug string         `db:"referential_slug"`
	ModelName       string         `db:"model_name"`
	Name            sql.NullString `db:"name"`
	ShortName       sql.NullString `db:"short_name"`
	StopAreaIds     sql.NullString `db:"stop_area_ids"`
}

type DatabaseLine struct {
	Id                string         `db:"id"`
	ReferentialSlug   string         `db:"referential_slug"`
	ModelName         string         `db:"model_name"`
	Name              string         `db:"name"`
	Codes             string         `db:"codes"`
	Attributes        string         `db:"attributes"`
	References        string         `db:"siri_references"`
	ReferentId        sql.NullString `db:"referent_id"`
	CollectSituations bool           `db:"collect_situations"`
}

type SelectLine struct {
	Id                string
	ReferentialSlug   string         `db:"referential_slug"`
	ModelName         string         `db:"model_name"`
	ReferentId        sql.NullString `db:"referent_id"`
	Name              sql.NullString
	Number            sql.NullString `db:"number"`
	Codes             sql.NullString `db:"codes"`
	Attributes        sql.NullString
	References        sql.NullString `db:"siri_references"`
	CollectSituations sql.NullBool   `db:"collect_situations"`
}

type DatabaseLineGroup struct {
	Id              string `db:"id"`
	ReferentialSlug string `db:"referential_slug"`
	ModelName       string `db:"model_name"`
	Name            string `db:"name"`
	ShortName       string `db:"short_name"`
	LineIds         string `db:"line_ids"`
}

type SelectLineGroup struct {
	Id              string
	ReferentialSlug string         `db:"referential_slug"`
	ModelName       string         `db:"model_name"`
	Name            sql.NullString `db:"name"`
	ShortName       sql.NullString `db:"short_name"`
	LineIds         sql.NullString `db:"line_ids"`
}

type DatabaseVehicleJourney struct {
	Id              string `db:"id"`
	ReferentialSlug string `db:"referential_slug"`
	ModelName       string `db:"model_name"`
	Name            string `db:"name"`
	Codes           string `db:"codes"`
	LineId          string `db:"line_id"`
	OriginName      string `db:"origin_name"`
	DestinationName string `db:"destination_name"`
	DirectionType   string `db:"direction_type"`
	Attributes      string `db:"attributes"`
	References      string `db:"siri_references"`
}

type SelectVehicleJourney struct {
	Id                  string
	ReferentialSlug     string `db:"referential_slug"`
	ModelName           string `db:"model_name"`
	Name                sql.NullString
	Codes               sql.NullString `db:"codes"`
	LineId              sql.NullString `db:"line_id"`
	OriginName          sql.NullString `db:"origin_name"`
	DestinationName     sql.NullString `db:"destination_name"`
	DirectionType       sql.NullString `db:"direction_type"`
	Attributes          sql.NullString
	References          sql.NullString `db:"siri_references"`
	AimedStopVisitCount sql.NullInt64  `db:"aimed_stop_visit_count"`
}

type DatabaseStopVisit struct {
	Id               string
	ReferentialSlug  string `db:"referential_slug"`
	ModelName        string `db:"model_name"`
	Codes            string `db:"codes"`
	StopAreaId       string `db:"stop_area_id"`
	VehicleJourneyId string `db:"vehicle_journey_id"`
	Schedules        string `db:"schedules"`
	Attributes       string `db:"attributes"`
	References       string `db:"siri_references"`
	PassageOrder     int    `db:"passage_order"`
}

type SelectStopVisit struct {
	Id               string
	ReferentialSlug  string         `db:"referential_slug"`
	ModelName        string         `db:"model_name"`
	Codes            sql.NullString `db:"codes"`
	StopAreaId       sql.NullString `db:"stop_area_id"`
	VehicleJourneyId sql.NullString `db:"vehicle_journey_id"`
	Schedules        sql.NullString `db:"schedules"`
	Attributes       sql.NullString `db:"attributes"`
	References       sql.NullString `db:"siri_references"`
	PassageOrder     sql.NullInt64  `db:"passage_order"`
}

type DatabaseFacility struct {
	Id              string `db:"id"`
	ReferentialSlug string `db:"referential_slug"`
	ModelName       string `db:"model_name"`
	Codes           string `db:"codes"`
}

type SelectFacility struct {
	Id              string
	ReferentialSlug string         `db:"referential_slug"`
	ModelName       string         `db:"model_name"`
	Codes           sql.NullString `db:"codes"`
}

// type DatabaseMacro struct {
// 	ReferentialId  string         `db:"referential_id"`
// 	Slug           string         `db:"slug"`
// 	Name           string         `db:"name"`
// 	Settings       string         `db:"settings"`
// 	Tokens         string         `db:"tokens"`
// 	ImportTokens   string         `db:"import_tokens"`
// 	OrganisationId sql.NullString `db:"organisation_id"`
// }

type SelectMacro struct {
	Id              string
	ReferentialSlug string         `db:"referential_slug"`
	ContextId       sql.NullString `db:"context_id"`
	Position        int
	Type            string
	ModelType       sql.NullString `db:"model_type"`
	Hook            sql.NullString
	Attributes      sql.NullString
}

func (m *SelectMacro) GetHook() sql.NullString {
	return m.Hook
}

func (m *SelectMacro) GetModelType() sql.NullString {
	return m.ModelType
}

func (m *SelectMacro) GetType() string {
	return m.Type
}

func (m *SelectMacro) GetAttributes() sql.NullString {
	return m.Attributes
}

type SelectControl struct {
	Id              string
	ReferentialSlug string         `db:"referential_slug"`
	ContextId       sql.NullString `db:"context_id"`
	Position        int
	Type            string
	ModelType       sql.NullString `db:"model_type"`
	Hook            sql.NullString
	Criticity       sql.NullString
	InternalCode    sql.NullString `db:"internal_code"`
	Attributes      sql.NullString
}

func (c *SelectControl) GetHook() sql.NullString {
	return c.Hook
}

func (c *SelectControl) GetModelType() sql.NullString {
	return c.ModelType
}

func (c *SelectControl) GetType() string {
	return c.Type
}

func (m *SelectControl) GetAttributes() sql.NullString {
	return m.Attributes
}
