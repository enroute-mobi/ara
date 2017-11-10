package model

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/af83/edwig/logger"
)

/* CSV Structure

stop_area,Id,ParentId,ModelName,Name,ObjectIDs,LineIds,Attributes,References,NextCollectAt,CollectedAt,CollectedUntil,CollectedAlways,CollectChildren
line,Id,ModelName,Name,ObjectIDs,Attributes,References,CollectGeneralMessages,CollectedAt
vehicle_journey,Id,ModelName,Name,ObjectIDs,LineId,Attributes,References
stop_visit,Id,ModelName,ObjectIDs,StopAreaId,VehicleJourneyId,ArrivalStatus,DepartureStatus,Schedules,Attributes,References,CollectedAt,RecordedAt,Collected,VehicleAtStop,PassageOrder

Comments are '#'
Separators are ',' leading spaces are trimed

Escape quotes with another quote ex: "[""1234"",""5678""]"
*/

func LoadFromCSV(filePath string, referentialSlug string) error {
	prepareDatabase()

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Error while opening file: %v", err)
	}

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	var i int
	for {
		i++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Debugf("Error while reading %v", err)
			continue
		}

		switch record[0] {
		case "stop_area":
			err := handleStopArea(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
			}
		case "line":
			err := handleLine(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
			}
		case "vehicle_journey":
			err := handleVehicleJourney(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
			}
		case "stop_visit":
			err := handleStopVisit(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
			}
		default:
			logger.Log.Debugf("Unknown record type: %v", record[0])
			continue
		}
	}
	return nil
}

func prepareDatabase() {
	Database.AddTableWithName(DatabaseStopArea{}, "stop_areas")
	Database.AddTableWithName(DatabaseLine{}, "lines")
	Database.AddTableWithName(DatabaseVehicleJourney{}, "vehicle_journeys")
	Database.AddTableWithName(DatabaseStopVisit{}, "stop_visits")
}

func handleStopArea(record []string, referentialSlug string) error {
	if len(record) != 14 {
		return fmt.Errorf("Wrong number of entries")
	}

	var err error
	parseErrors := make(map[string]string)
	var parent sql.NullString
	if record[2] != "" {
		parent = sql.NullString{
			String: record[2],
			Valid:  true,
		}
	}

	var nextCollectAt time.Time
	if record[9] != "" {
		nextCollectAt, err = time.Parse(time.RFC3339, record[9])
		if err != nil {
			parseErrors["NextCollectAt"] = err.Error()
		}
	}

	var collectedAt time.Time
	if record[10] != "" {
		collectedAt, err = time.Parse(time.RFC3339, record[10])
		if err != nil {
			parseErrors["CollectedAt"] = err.Error()
		}
	}

	var collectedUntil time.Time
	if record[11] != "" {
		collectedUntil, err = time.Parse(time.RFC3339, record[11])
		if err != nil {
			parseErrors["CollectedUntil"] = err.Error()
		}
	}

	var collectedAlways bool
	if record[12] != "" {
		collectedAlways, err = strconv.ParseBool(record[12])
		if err != nil {
			parseErrors["CollectedAlways"] = err.Error()
		}
	}

	var collectChildren bool
	if record[13] != "" {
		collectChildren, err = strconv.ParseBool(record[13])
		if err != nil {
			parseErrors["CollectChildren"] = err.Error()
		}
	}

	if len(parseErrors) != 0 {
		json, _ := json.Marshal(parseErrors)
		return fmt.Errorf(string(json))
	}

	stopArea := DatabaseStopArea{
		Id:              record[1],
		ReferentialSlug: referentialSlug,
		ParentId:        parent,
		ModelName:       record[3],
		Name:            record[4],
		ObjectIDs:       record[5],
		LineIds:         record[6],
		Attributes:      record[7],
		References:      record[8],
		NextCollectAt:   nextCollectAt,
		CollectedAt:     collectedAt,
		CollectedUntil:  collectedUntil,
		CollectedAlways: collectedAlways,
		CollectChildren: collectChildren,
	}

	err = Database.Insert(&stopArea)
	if err != nil {
		return err
	}

	return nil
}

func handleLine(record []string, referentialSlug string) error {
	if len(record) != 9 {
		return fmt.Errorf("Wrong number of entries")
	}

	var err error
	parseErrors := make(map[string]string)

	var collectGeneralMessages bool
	if record[7] != "" {
		collectGeneralMessages, err = strconv.ParseBool(record[7])
		if err != nil {
			parseErrors["CollectGeneralMessages"] = err.Error()
		}
	}

	var collectedAt time.Time
	if record[8] != "" {
		collectedAt, err = time.Parse(time.RFC3339, record[8])
		if err != nil {
			parseErrors["CollectedAt"] = err.Error()
		}
	}

	line := DatabaseLine{
		Id:                     record[1],
		ReferentialSlug:        referentialSlug,
		ModelName:              record[2],
		Name:                   record[3],
		ObjectIDs:              record[4],
		Attributes:             record[5],
		References:             record[6],
		CollectGeneralMessages: collectGeneralMessages,
		CollectedAt:            collectedAt,
	}

	err = Database.Insert(&line)
	if err != nil {
		return err
	}

	return nil
}

func handleVehicleJourney(record []string, referentialSlug string) error {
	if len(record) != 8 {
		return fmt.Errorf("Wrong number of entries")
	}

	vehicleJourney := DatabaseVehicleJourney{
		Id:              record[1],
		ReferentialSlug: referentialSlug,
		ModelName:       record[2],
		Name:            record[3],
		ObjectIDs:       record[4],
		LineId:          record[5],
		Attributes:      record[6],
		References:      record[7],
	}

	err := Database.Insert(&vehicleJourney)
	if err != nil {
		return err
	}

	return nil
}

func handleStopVisit(record []string, referentialSlug string) error {
	if len(record) != 16 {
		return fmt.Errorf("Wrong number of entries")
	}

	var err error
	parseErrors := make(map[string]string)

	var collectedAt time.Time
	if record[11] != "" {
		collectedAt, err = time.Parse(time.RFC3339, record[11])
		if err != nil {
			parseErrors["CollectedAt"] = err.Error()
		}
	}

	var recordedAt time.Time
	if record[12] != "" {
		recordedAt, err = time.Parse(time.RFC3339, record[12])
		if err != nil {
			parseErrors["RecordedAt"] = err.Error()
		}
	}

	var collected bool
	if record[13] != "" {
		collected, err = strconv.ParseBool(record[13])
		if err != nil {
			parseErrors["Collected"] = err.Error()
		}
	}

	var vehicleAtStop bool
	if record[14] != "" {
		vehicleAtStop, err = strconv.ParseBool(record[14])
		if err != nil {
			parseErrors["VehicleAtStop"] = err.Error()
		}
	}

	var passageOrder int
	if record[15] != "" {
		passageOrder, err = strconv.Atoi(record[15])
		if err != nil {
			parseErrors["PassageOrder"] = err.Error()
		}
	}

	if len(parseErrors) != 0 {
		json, _ := json.Marshal(parseErrors)
		return fmt.Errorf(string(json))
	}

	stopVisit := DatabaseStopVisit{
		Id:               record[1],
		ReferentialSlug:  referentialSlug,
		ModelName:        record[2],
		ObjectIDs:        record[3],
		StopAreaId:       record[4],
		VehicleJourneyId: record[5],
		ArrivalStatus:    record[6],
		DepartureStatus:  record[7],
		Schedules:        record[8],
		Attributes:       record[9],
		References:       record[10],
		CollectedAt:      collectedAt,
		RecordedAt:       recordedAt,
		Collected:        collected,
		VehicleAtStop:    vehicleAtStop,
		PassageOrder:     passageOrder,
	}

	err = Database.Insert(&stopVisit)
	if err != nil {
		return err
	}

	return nil
}
