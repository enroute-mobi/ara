package model

import (
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

stop_area,Id,ReferentialId,ParentId,ModelName,Name,ObjectIDs,LineIds,Attributes,References,NextCollectAt,CollectedAt,CollectedUntil,CollectedAlways,CollectChildren
line,Id,ReferentialId,ModelName,Name,ObjectIDs,Attributes,References
vehicle_journey,Id,ReferentialId,ModelName,Name,ObjectIDs,LineId,Attributes,References
stop_visit,Id,ReferentialId,ModelName,ObjectIDs,StopAreaId,VehicleJourneyId,ArrivalStatus,DepartureStatus,Schedules,Attributes,References,CollectedAt,RecordedAt,Collected,VehicleAtStop,PassageOrder

Comments are '#'
Separators are ',' leading spaces are trimed

Escape quotes with another quote ex: "[""1234"",""5678""]"
*/

func LoadFromCSV(filePath string) error {
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
			err := handleStopArea(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
			}
		case "line":
			err := handleLine(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
			}
		case "vehicle_journey":
			err := handleVehicleJourney(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
			}
		case "stop_visit":
			err := handleStopVisit(record)
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

func handleStopArea(record []string) error {
	if len(record) != 15 {
		return fmt.Errorf("Wrong number of entries")
	}

	var err error
	parseErrors := make(map[string]string)

	var nextCollectAt time.Time
	if record[10] != "" {
		nextCollectAt, err = time.Parse(time.RFC3339, record[10])
		if err != nil {
			parseErrors["NextCollectAt"] = err.Error()
		}
	}

	var collectedAt time.Time
	if record[11] != "" {
		collectedAt, err = time.Parse(time.RFC3339, record[11])
		if err != nil {
			parseErrors["CollectedAt"] = err.Error()
		}
	}

	var collectedUntil time.Time
	if record[12] != "" {
		collectedUntil, err = time.Parse(time.RFC3339, record[12])
		if err != nil {
			parseErrors["CollectedUntil"] = err.Error()
		}
	}

	var collectedAlways bool
	if record[13] != "" {
		collectedAlways, err = strconv.ParseBool(record[13])
		if err != nil {
			parseErrors["CollectedAlways"] = err.Error()
		}
	}

	var collectChildren bool
	if record[14] != "" {
		collectChildren, err = strconv.ParseBool(record[14])
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
		ReferentialId:   record[2],
		ParentId:        record[3],
		ModelName:       record[4],
		Name:            record[5],
		ObjectIDs:       record[6],
		LineIds:         record[7],
		Attributes:      record[8],
		References:      record[9],
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

func handleLine(record []string) error {
	if len(record) != 8 {
		return fmt.Errorf("Wrong number of entries")
	}

	line := DatabaseLine{
		Id:            record[1],
		ReferentialId: record[2],
		ModelName:     record[3],
		Name:          record[4],
		ObjectIDs:     record[5],
		Attributes:    record[6],
		References:    record[7],
	}

	err := Database.Insert(&line)
	if err != nil {
		return err
	}

	return nil
}

func handleVehicleJourney(record []string) error {
	if len(record) != 9 {
		return fmt.Errorf("Wrong number of entries")
	}

	vehicleJourney := DatabaseVehicleJourney{
		Id:            record[1],
		ReferentialId: record[2],
		ModelName:     record[3],
		Name:          record[4],
		ObjectIDs:     record[5],
		LineId:        record[6],
		Attributes:    record[7],
		References:    record[8],
	}

	err := Database.Insert(&vehicleJourney)
	if err != nil {
		return err
	}

	return nil
}

func handleStopVisit(record []string) error {
	if len(record) != 17 {
		return fmt.Errorf("Wrong number of entries")
	}

	var err error
	parseErrors := make(map[string]string)

	var collectedAt time.Time
	if record[12] != "" {
		collectedAt, err = time.Parse(time.RFC3339, record[12])
		if err != nil {
			parseErrors["CollectedAt"] = err.Error()
		}
	}

	var recordedAt time.Time
	if record[13] != "" {
		recordedAt, err = time.Parse(time.RFC3339, record[13])
		if err != nil {
			parseErrors["RecordedAt"] = err.Error()
		}
	}

	var collected bool
	if record[14] != "" {
		collected, err = strconv.ParseBool(record[14])
		if err != nil {
			parseErrors["Collected"] = err.Error()
		}
	}

	var vehicleAtStop bool
	if record[15] != "" {
		vehicleAtStop, err = strconv.ParseBool(record[15])
		if err != nil {
			parseErrors["VehicleAtStop"] = err.Error()
		}
	}

	var passageOrder int
	if record[16] != "" {
		passageOrder, err = strconv.Atoi(record[16])
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
		ReferentialId:    record[2],
		ModelName:        record[3],
		ObjectIDs:        record[4],
		StopAreaId:       record[5],
		VehicleJourneyId: record[6],
		ArrivalStatus:    record[7],
		DepartureStatus:  record[8],
		Schedules:        record[9],
		Attributes:       record[10],
		References:       record[11],
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
