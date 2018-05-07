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

stop_area,Id,ParentId,ReferentId,ModelName,Name,ObjectIDs,LineIds,Attributes,References,CollectedAlways,CollectChildren,CollectGeneralMessages
line,Id,ModelName,Name,ObjectIDs,Attributes,References,CollectGeneralMessages
vehicle_journey,Id,ModelName,Name,ObjectIDs,LineId,OriginName,DestinationName,Attributes,References
stop_visit,Id,ModelName,ObjectIDs,StopAreaId,VehicleJourneyId,PassageOrder,Schedules,Attributes,References

Comments are '#'
Separators are ',' leading spaces are trimed

Escape quotes with another quote ex: "[""1234"",""5678""]"
*/
type Loader struct {
	filePath        string
	referentialSlug string
	force           bool
}

func LoadFromCSV(filePath string, referentialSlug string, force bool) error {
	return newLoader(filePath, referentialSlug, force).load()
}

func newLoader(filePath string, referentialSlug string, force bool) *Loader {
	return &Loader{
		filePath:        filePath,
		referentialSlug: referentialSlug,
		force:           force,
	}
}

func (loader Loader) load() error {
	prepareDatabase()
	var errors int

	file, err := os.Open(loader.filePath)
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
			err := loader.handleStopArea(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedStopAreas++
			}
		case "line":
			err := loader.handleLine(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedLines++
			}
		case "vehicle_journey":
			err := loader.handleVehicleJourney(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedVehicleJourneys++
			}
		case "stop_visit":
			err := loader.handleStopVisit(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Errorf("Error on line %d: %v", i, err)
				errors++
			} else {
				importedStopVisits++
			}
		case "operator":
			err := loader.handleOperator(record)
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

func (loader Loader) handleStopArea(record []string) error {
	if len(record) != 13 {
		return fmt.Errorf("Wrong number of entries, expected 13 got %v", len(record))
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

	var referent sql.NullString
	if record[3] != "" {
		referent = sql.NullString{
			String: record[3],
			Valid:  true,
		}
	}

	var collectedAlways bool
	if record[10] != "" {
		collectedAlways, err = strconv.ParseBool(record[10])
		if err != nil {
			parseErrors["CollectedAlways"] = err.Error()
		}
	}

	var collectChildren bool
	if record[11] != "" {
		collectChildren, err = strconv.ParseBool(record[11])
		if err != nil {
			parseErrors["CollectChildren"] = err.Error()
		}
	}

	var collectGeneralMessages bool
	if record[12] != "" {
		collectGeneralMessages, err = strconv.ParseBool(record[12])
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
		ReferentialSlug:        loader.referentialSlug,
		ParentId:               parent,
		ReferentId:             referent,
		ModelName:              record[4],
		Name:                   record[5],
		ObjectIDs:              record[6],
		LineIds:                record[7],
		Attributes:             record[8],
		References:             record[9],
		CollectedAlways:        collectedAlways,
		CollectChildren:        collectChildren,
		CollectGeneralMessages: collectGeneralMessages,
	}

	if loader.force {
		query := fmt.Sprintf("delete from stop_areas where id='%v' and model_name='%v'", stopArea.Id, stopArea.ModelName)
		_, err := Database.Exec(query)
		if err != nil {
			return err
		}
	}

	err = Database.Insert(&stopArea)
	if err != nil {
		return err
	}

	return nil
}

func (loader Loader) handleOperator(record []string) error {
	if len(record) != 5 {
		return fmt.Errorf("Wrong number of entries, expected 5 got %v", len(record))
	}

	operator := DatabaseOperator{
		Id:              record[1],
		ReferentialSlug: loader.referentialSlug,
		ModelName:       record[2],
		Name:            record[3],
		ObjectIDs:       record[4],
	}

	if loader.force {
		query := fmt.Sprintf("delete from operators where id='%v' and model_name='%v'", operator.Id, operator.ModelName)
		_, err := Database.Exec(query)
		if err != nil {
			return err
		}
	}

	err := Database.Insert(&operator)
	if err != nil {
		return err
	}

	return nil
}

func (loader Loader) handleLine(record []string) error {
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
		ReferentialSlug:        loader.referentialSlug,
		ModelName:              record[2],
		Name:                   record[3],
		ObjectIDs:              record[4],
		Attributes:             record[5],
		References:             record[6],
		CollectGeneralMessages: collectGeneralMessages,
	}

	if loader.force {
		query := fmt.Sprintf("delete from lines where id='%v' and model_name='%v'", line.Id, line.ModelName)
		_, err := Database.Exec(query)
		if err != nil {
			return err
		}
	}

	err = Database.Insert(&line)
	if err != nil {
		return err
	}

	return nil
}

func (loader Loader) handleVehicleJourney(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("Wrong number of entries, expected 10 got %v", len(record))
	}

	vehicleJourney := DatabaseVehicleJourney{
		Id:              record[1],
		ReferentialSlug: loader.referentialSlug,
		ModelName:       record[2],
		Name:            record[3],
		ObjectIDs:       record[4],
		LineId:          record[5],
		OriginName:      record[6],
		DestinationName: record[7],
		Attributes:      record[8],
		References:      record[9],
	}

	if loader.force {
		query := fmt.Sprintf("delete from vehicle_journeys where id='%v' and model_name='%v'", vehicleJourney.Id, vehicleJourney.ModelName)
		_, err := Database.Exec(query)
		if err != nil {
			return err
		}
	}

	err := Database.Insert(&vehicleJourney)
	if err != nil {
		return err
	}

	return nil
}

func (loader Loader) handleStopVisit(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("Wrong number of entries, expected 10 got %v", len(record))
	}

	var err error
	parseErrors := make(map[string]string)

	var passageOrder int
	if record[6] != "" {
		passageOrder, err = strconv.Atoi(record[6])
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
		ReferentialSlug:  loader.referentialSlug,
		ModelName:        record[2],
		ObjectIDs:        record[3],
		StopAreaId:       record[4],
		VehicleJourneyId: record[5],
		PassageOrder:     passageOrder,
		Schedules:        record[7],
		Attributes:       record[8],
		References:       record[9],
	}

	if loader.force {
		query := fmt.Sprintf("delete from stop_visits where id='%v' and model_name='%v'", stopVisit.Id, stopVisit.ModelName)
		_, err := Database.Exec(query)
		if err != nil {
			return err
		}
	}

	err = Database.Insert(&stopVisit)
	if err != nil {
		return err
	}

	return nil
}
