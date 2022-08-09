package model

import "database/sql"

type DatabaseReferential struct {
	ReferentialId  string         `db:"referential_id"`
	OrganisationId sql.NullString `db:"organisation_id"`
	Slug           string         `db:"slug"`
	Name           string         `db:"name"`
	Settings       string         `db:"settings"`
	Tokens         string         `db:"tokens"`
	ImportTokens   string         `db:"import_tokens"`
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

type DatabaseOperator struct {
	Id              string `db:"id"`
	ReferentialSlug string `db:"referential_slug"`
	Name            string `db:"name"`
	ObjectIDs       string `db:"object_ids"`
	ModelName       string `db:"model_name"`
}

type SelectOperator struct {
	Id              string
	ReferentialSlug string `db:"referential_slug"`
	Name            sql.NullString
	ObjectIDs       sql.NullString `db:"object_ids"`
	ModelName       string         `db:"model_name"`
}

type DatabaseStopArea struct {
	Id                     string         `db:"id"`
	ReferentialSlug        string         `db:"referential_slug"`
	ParentId               sql.NullString `db:"parent_id"`
	ReferentId             sql.NullString `db:"referent_id"`
	ModelName              string         `db:"model_name"`
	Name                   string         `db:"name"`
	ObjectIDs              string         `db:"object_ids"`
	LineIds                string         `db:"line_ids"`
	Attributes             string         `db:"attributes"`
	References             string         `db:"siri_references"`
	CollectedAlways        bool           `db:"collected_always"`
	CollectChildren        bool           `db:"collect_children"`
	CollectGeneralMessages bool           `db:"collect_general_messages"`
}

type SelectStopArea struct {
	Id                     string
	ReferentialSlug        string `db:"referential_slug"`
	ModelName              string `db:"model_name"`
	Name                   sql.NullString
	ObjectIDs              sql.NullString `db:"object_ids"`
	ParentId               sql.NullString `db:"parent_id"`
	ReferentId             sql.NullString `db:"referent_id"`
	Attributes             sql.NullString
	References             sql.NullString `db:"siri_references"`
	LineIds                sql.NullString `db:"line_ids"`
	CollectedAlways        sql.NullBool   `db:"collected_always"`
	CollectChildren        sql.NullBool   `db:"collect_children"`
	CollectGeneralMessages sql.NullBool   `db:"collect_general_messages"`
}

type DatabaseLine struct {
	Id                     string `db:"id"`
	ReferentialSlug        string `db:"referential_slug"`
	ModelName              string `db:"model_name"`
	Name                   string `db:"name"`
	ObjectIDs              string `db:"object_ids"`
	Attributes             string `db:"attributes"`
	References             string `db:"siri_references"`
	CollectGeneralMessages bool   `db:"collect_general_messages"`
}

type SelectLine struct {
	Id                     string
	ReferentialSlug        string `db:"referential_slug"`
	ModelName              string `db:"model_name"`
	Name                   sql.NullString
	Number                 sql.NullString `db:"number"`
	ObjectIDs              sql.NullString `db:"object_ids"`
	Attributes             sql.NullString
	References             sql.NullString `db:"siri_references"`
	CollectGeneralMessages sql.NullBool   `db:"collect_general_messages"`
}

type DatabaseVehicleJourney struct {
	Id              string `db:"id"`
	ReferentialSlug string `db:"referential_slug"`
	ModelName       string `db:"model_name"`
	Name            string `db:"name"`
	ObjectIDs       string `db:"object_ids"`
	LineId          string `db:"line_id"`
	OriginName      string `db:"origin_name"`
	DestinationName string `db:"destination_name"`
	DirectionType   string `db:"direction_type"`
	Attributes      string `db:"attributes"`
	References      string `db:"siri_references"`
}

type SelectVehicleJourney struct {
	Id              string
	ReferentialSlug string `db:"referential_slug"`
	ModelName       string `db:"model_name"`
	Name            sql.NullString
	ObjectIDs       sql.NullString `db:"object_ids"`
	LineId          sql.NullString `db:"line_id"`
	OriginName      sql.NullString `db:"origin_name"`
	DestinationName sql.NullString `db:"destination_name"`
	DirectionType   sql.NullString `db:"direction_type"`
	Attributes      sql.NullString
	References      sql.NullString `db:"siri_references"`
}

type DatabaseStopVisit struct {
	Id               string
	ReferentialSlug  string `db:"referential_slug"`
	ModelName        string `db:"model_name"`
	ObjectIDs        string `db:"object_ids"`
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
	ObjectIDs        sql.NullString `db:"object_ids"`
	StopAreaId       sql.NullString `db:"stop_area_id"`
	VehicleJourneyId sql.NullString `db:"vehicle_journey_id"`
	Schedules        sql.NullString `db:"schedules"`
	Attributes       sql.NullString `db:"attributes"`
	References       sql.NullString `db:"siri_references"`
	PassageOrder     sql.NullInt64  `db:"passage_order"`
}
