package model

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/af83/edwig/logger"
)

/* CSV Structure

stop_area,Id,ParentId,ModelName,Name,ObjectIDs,LineIds,Attributes,References,CollectedAlways,CollectChildren,CollectGeneralMessages
line,Id,ModelName,Name,ObjectIDs,Attributes,References,CollectGeneralMessages
vehicle_journey,Id,ModelName,Name,ObjectIDs,LineId,Attributes,References
stop_visit,Id,ModelName,ObjectIDs,StopAreaId,VehicleJourneyId,ArrivalStatus,DepartureStatus,Schedules,Attributes,References,Collected,VehicleAtStop,PassageOrder

Comments are '#'
Separators are ',' leading spaces are trimed

Escape quotes with another quote ex: "[""1234"",""5678""]"
*/

func LoadFromCSV(filePath string, referentialSlug string) error {
	prepareDatabase()
	var errors int

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

	importedStopAreas := 0
	importedLines := 0
	importedVehicleJourneys := 0
	importedStopVisits := 0
	importedOperators := 0

	var i int
	for {
		i++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Debugf("Error while reading %v", err)
			fmt.Errorf("Error while reading %v", err)
			errors++
			continue
		}

		switch record[0] {
		case "stop_area":
			err := handleStopArea(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedStopAreas++
			}
		case "line":
			err := handleLine(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedLines++
			}
		case "vehicle_journey":
			err := handleVehicleJourney(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedVehicleJourneys++
			}
		case "stop_visit":
			err := handleStopVisit(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedStopVisits++
			}
		case "operator":
			err := handleOperator(record, referentialSlug)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedOperators++
			}
		default:
			logger.Log.Debugf("Unknown record type: %v", record[0])
			fmt.Errorf("Unknown record type: %v", record[0])
			errors++
			continue
		}
	}

	if (importedStopAreas + importedLines + importedVehicleJourneys + importedStopVisits + importedOperators) == 0 {
		logger.Log.Debugf("Nothing imported, import raised %v errors", errors)
		fmt.Printf("Nothing imported, import raised %v errors\n", errors)
	} else {
		logger.Log.Debugf("Import successful, import raised %v errors", errors)
		logger.Log.Debugf("  %v StopAreas", importedStopAreas)
		logger.Log.Debugf("  %v Lines", importedLines)
		logger.Log.Debugf("  %v VehicleJourneys", importedVehicleJourneys)
		logger.Log.Debugf("  %v StopVisits", importedStopVisits)
		logger.Log.Debugf("  %v Operators", importedOperators)

		fmt.Printf("Import successful, import raised %v errors\n", errors)
		fmt.Printf("  %v StopAreas\n", importedStopAreas)
		fmt.Printf("  %v Lines\n", importedLines)
		fmt.Printf("  %v VehicleJourneys\n", importedVehicleJourneys)
		fmt.Printf("  %v StopVisits\n", importedStopVisits)
		fmt.Printf("  %v Operators\n", importedOperators)
	}

	return nil
}

func prepareDatabase() {
	Database.AddTableWithName(DatabaseStopArea{}, "stop_areas")
	Database.AddTableWithName(DatabaseLine{}, "lines")
	Database.AddTableWithName(DatabaseVehicleJourney{}, "vehicle_journeys")
	Database.AddTableWithName(DatabaseStopVisit{}, "stop_visits")
	Database.AddTableWithName(DatabaseOperator{}, "operators")
}

func handleStopArea(record []string, referentialSlug string) error {
	if len(record) != 12 {
		return fmt.Errorf("Wrong number of entries, expected 11 got %v", len(record))
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

	var collectedAlways bool
	if record[9] != "" {
		collectedAlways, err = strconv.ParseBool(record[9])
		if err != nil {
			parseErrors["CollectedAlways"] = err.Error()
		}
	}

	var collectChildren bool
	if record[10] != "" {
		collectChildren, err = strconv.ParseBool(record[10])
		if err != nil {
			parseErrors["CollectChildren"] = err.Error()
		}
	}

	var collectGeneralMessages bool
	if record[11] != "" {
		collectGeneralMessages, err = strconv.ParseBool(record[11])
		if err != nil {
			parseErrors["CollectGeneralMessages"] = err.Error()
		}
	}

	if len(parseErrors) != 0 {
		json, _ := json.Marshal(parseErrors)
		return fmt.Errorf(string(json))
	}

	stopArea := DatabaseStopArea{
		Id:                     record[1],
		ReferentialSlug:        referentialSlug,
		ParentId:               parent,
		ModelName:              record[3],
		Name:                   record[4],
		ObjectIDs:              record[5],
		LineIds:                record[6],
		Attributes:             record[7],
		References:             record[8],
		CollectedAlways:        collectedAlways,
		CollectChildren:        collectChildren,
		CollectGeneralMessages: collectGeneralMessages,
	}

	err = Database.Insert(&stopArea)
	if err != nil {
		return err
	}

	return nil
}

func handleOperator(record []string, referentialSlug string) error {
	if len(record) != 5 {
		return fmt.Errorf("Wrong number of entries, expected 5 got %v", len(record))
	}

	operator := DatabaseOperator{
		Id:              record[1],
		ReferentialSlug: referentialSlug,
		ModelName:       record[2],
		Name:            record[3],
		ObjectIDs:       record[4],
	}

	err := Database.Insert(&operator)
	if err != nil {
		return err
	}

	return nil
}

func handleLine(record []string, referentialSlug string) error {
	if len(record) != 8 {
		return fmt.Errorf("Wrong number of entries, expected 8 got %v", len(record))
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

	if len(parseErrors) != 0 {
		json, _ := json.Marshal(parseErrors)
		return fmt.Errorf(string(json))
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
	}

	err = Database.Insert(&line)
	if err != nil {
		return err
	}

	return nil
}

func handleVehicleJourney(record []string, referentialSlug string) error {
	if len(record) != 8 {
		return fmt.Errorf("Wrong number of entries, expected 8 got %v", len(record))
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
	if len(record) != 14 {
		return fmt.Errorf("Wrong number of entries, expected 14 got %v", len(record))
	}

	var err error
	parseErrors := make(map[string]string)

	var collected bool
	if record[11] != "" {
		collected, err = strconv.ParseBool(record[11])
		if err != nil {
			parseErrors["Collected"] = err.Error()
		}
	}

	var vehicleAtStop bool
	if record[12] != "" {
		vehicleAtStop, err = strconv.ParseBool(record[12])
		if err != nil {
			parseErrors["VehicleAtStop"] = err.Error()
		}
	}

	var passageOrder int
	if record[13] != "" {
		passageOrder, err = strconv.Atoi(record[13])
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
