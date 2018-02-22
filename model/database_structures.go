package model

import "database/sql"

type DatabaseReferential struct {
	ReferentialId string `db:"referential_id"`
	Slug          string `db:"slug"`
	Settings      string `db:"settings"`
	Tokens        string `db:"tokens"`
}

type SelectReferential struct {
	Referential_id string
	Slug           string
	Settings       sql.NullString
	Tokens         sql.NullString
}

type DatabasePartner struct {
	Id             string `db:"id"`
	ReferentialId  string `db:"referential_id"`
	Slug           string `db:"slug"`
	Settings       string `db:"settings"`
	ConnectorTypes string `db:"connector_types"`
}

type SelectPartner struct {
	Id             string
	ReferentialId  string `db:"referential_id"`
	Slug           string
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
	ArrivalStatus    string `db:"arrival_status"`
	DepartureStatus  string `db:"departure_status"`
	Schedules        string `db:"schedules"`
	Attributes       string `db:"attributes"`
	References       string `db:"siri_references"`
	Collected        bool   `db:"collected"`
	VehicleAtStop    bool   `db:"vehicle_at_stop"`
	PassageOrder     int    `db:"passage_order"`
}

type SelectStopVisit struct {
	Id               string
	ReferentialSlug  string         `db:"referential_slug"`
	ModelName        string         `db:"model_name"`
	ObjectIDs        sql.NullString `db:"object_ids"`
	StopAreaId       sql.NullString `db:"stop_area_id"`
	VehicleJourneyId sql.NullString `db:"vehicle_journey_id"`
	ArrivalStatus    sql.NullString `db:"arrival_status"`
	DepartureStatus  sql.NullString `db:"departure_status"`
	Schedules        sql.NullString `db:"schedules"`
	Attributes       sql.NullString `db:"attributes"`
	References       sql.NullString `db:"siri_references"`
	Collected        sql.NullBool   `db:"collected"`
	VehicleAtStop    sql.NullBool   `db:"vehicle_at_stop"`
	PassageOrder     sql.NullInt64  `db:"passage_order"`
}
