package model

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/af83/edwig/config"
	"github.com/af83/edwig/logger"
)

/* CSV Structure

operator,Id,ModelName,Name,ObjectIDs
stop_area,Id,ParentId,ReferentId,ModelName,Name,ObjectIDs,LineIds,Attributes,References,CollectedAlways,CollectChildren,CollectGeneralMessages
line,Id,ModelName,Name,ObjectIDs,Attributes,References,CollectGeneralMessages
vehicle_journey,Id,ModelName,Name,ObjectIDs,LineId,OriginName,DestinationName,Attributes,References
stop_visit,Id,ModelName,ObjectIDs,StopAreaId,VehicleJourneyId,PassageOrder,Schedules,Attributes,References

Comments are '#'
Separators are ',' leading spaces are trimed

Escape quotes with another quote ex: "[""1234"",""5678""]"
*/
const (
	STOP_AREA       = "stop_area"
	LINE            = "line"
	VEHICLE_JOURNEY = "vehicle_journey"
	STOP_VISIT      = "stop_visit"
	OPERATOR        = "operator"
)

type Loader struct {
	filePath        string
	referentialSlug string
	force           bool
	deletedModels   map[string]map[string]struct{}
	operators       []byte
	stopAreas       []byte
	lines           []byte
	vehicleJourneys []byte
	stopVisits      []byte
	bulkCounter     map[string]int
	insertedCount   map[string]int64
	errors          int
}

func LoadFromCSV(filePath string, referentialSlug string, force bool) error {
	return newLoader(filePath, referentialSlug, force).load()
}

func newLoader(filePath string, referentialSlug string, force bool) *Loader {
	d := make(map[string]map[string]struct{})
	for _, m := range [5]string{STOP_AREA, LINE, VEHICLE_JOURNEY, STOP_VISIT, OPERATOR} {
		d[m] = make(map[string]struct{})
	}
	return &Loader{
		filePath:        filePath,
		referentialSlug: referentialSlug,
		force:           force,
		deletedModels:   d,
		bulkCounter:     make(map[string]int),
		insertedCount:   make(map[string]int64),
	}
}

func (loader Loader) load() error {
	file, err := os.Open(loader.filePath)
	if err != nil {
		return fmt.Errorf("loader error: error while opening file: %v", err)
	}
	defer file.Close()

	// Config CSV reader
	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	startTime := time.Now()
	logger.Log.Debugf("Load operation started at %v", startTime)

	var i int
	for {
		i++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Debugf("Error while reading: %v", err)
			fmt.Printf("Error while reading: %v\n", err)
			loader.errors++
			continue
		}

		switch record[0] {
		case OPERATOR:
			err := loader.handleOperator(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Printf("Error on line %d: %v\n", i, err)
				loader.errors++
			}
		case STOP_AREA:
			err := loader.handleStopArea(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Printf("Error on line %d: %v\n", i, err)
				loader.errors++
			}
		case LINE:
			err := loader.handleLine(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Printf("Error on line %d: %v\n", i, err)
				loader.errors++
			}
		case VEHICLE_JOURNEY:
			err := loader.handleVehicleJourney(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Printf("Error on line %d: %v\n", i, err)
				loader.errors++
			}
		case STOP_VISIT:
			err := loader.handleStopVisit(record)
			if err != nil {
				logger.Log.Debugf("Error on line %d: %v", i, err)
				fmt.Printf("Error on line %d: %v\n", i, err)
				loader.errors++
			}
		default:
			logger.Log.Debugf("Unknown record type: %v", record[0])
			fmt.Printf("Unknown record type: %v\n", record[0])
			loader.errors++
			continue
		}
	}

	loader.insertOperators()
	loader.insertStopAreas()
	loader.insertLines()
	loader.insertVehicleJourneys()
	loader.insertStopVisits()

	logger.Log.Debugf("Load operation done in %v", time.Since(startTime))

	if (loader.insertedCount[OPERATOR] + loader.insertedCount[STOP_AREA] + loader.insertedCount[LINE] + loader.insertedCount[VEHICLE_JOURNEY] + loader.insertedCount[STOP_VISIT]) == 0 {
		if loader.errors == 0 {
			return fmt.Errorf("loader error: empty file")
		}
		return fmt.Errorf("loader error: couldn't import anything, import raised %v errors", loader.errors)
	}

	logger.Log.Debugf("Import successful, import raised %v errors", loader.errors)
	logger.Log.Debugf("  %v Operators", loader.insertedCount[OPERATOR])
	logger.Log.Debugf("  %v StopAreas", loader.insertedCount[STOP_AREA])
	logger.Log.Debugf("  %v Lines", loader.insertedCount[LINE])
	logger.Log.Debugf("  %v VehicleJourneys", loader.insertedCount[VEHICLE_JOURNEY])
	logger.Log.Debugf("  %v StopVisits", loader.insertedCount[STOP_VISIT])

	fmt.Printf("Import successful, import raised %v errors\n", loader.errors)
	fmt.Printf("  %v Operators\n", loader.insertedCount[OPERATOR])
	fmt.Printf("  %v StopAreas\n", loader.insertedCount[STOP_AREA])
	fmt.Printf("  %v Lines\n", loader.insertedCount[LINE])
	fmt.Printf("  %v VehicleJourneys\n", loader.insertedCount[VEHICLE_JOURNEY])
	fmt.Printf("  %v StopVisits\n", loader.insertedCount[STOP_VISIT])

	return nil
}

func (loader *Loader) handleForce(klass, modelName string) error {
	if loader.force {
		if _, ok := loader.deletedModels[klass][modelName]; !ok {
			loader.deletedModels[klass][modelName] = struct{}{}
			query := fmt.Sprintf("delete from %vs where model_name='%v'", klass, modelName)
			_, err := Database.Exec(query)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (loader *Loader) handleOperator(record []string) error {
	if len(record) != 5 {
		return fmt.Errorf("wrong number of entries, expected 5 got %v", len(record))
	}

	err := loader.handleForce(OPERATOR, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$),", loader.referentialSlug, record[1], record[2], record[3], record[4])
	loader.operators = append(loader.operators, values...)
	loader.bulkCounter[OPERATOR]++

	if loader.bulkCounter[OPERATOR] >= config.Config.LoadMaxInsert {
		loader.insertOperators()
	}

	return nil
}

func (loader *Loader) insertOperators() {
	if len(loader.operators) == 0 {
		return
	}

	defer func() {
		loader.operators = []byte{}
		loader.bulkCounter[OPERATOR] = 0
	}()

	query := fmt.Sprintf("INSERT INTO operators(referential_slug,id,model_name,name,object_ids) VALUES %v;", string(loader.operators[:len(loader.operators)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		logger.Log.Debugf("Error while inserting operators: %v", err)
		fmt.Printf("Error while inserting operators: %v", err)
		loader.errors++
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		logger.Log.Debugf("Unexpected error while inserting operators: %v", err)
		fmt.Printf("Unexpected error while inserting operators: %v", err)
		loader.errors++
		return
	}

	loader.insertedCount[OPERATOR] += rows
}

func (loader *Loader) handleStopArea(record []string) error {
	if len(record) != 13 {
		return fmt.Errorf("wrong number of entries, expected 13 got %v", len(record))
	}

	var err error
	parseErrors := make(map[string]string)

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

	err = loader.handleForce(STOP_AREA, record[4])
	if err != nil {
		return err
	}

	var parent string
	if record[2] != "" {
		parent = fmt.Sprintf("$$%v$$", record[2])
	} else {
		parent = "null"
	}

	var referent string
	if record[3] != "" {
		referent = fmt.Sprintf("$$%v$$", record[3])
	} else {
		referent = "null"
	}

	values := fmt.Sprintf("($$%v$$, $$%v$$, %v, %v, $$%v$$, $$%v$$, $$%v$$, $$%v$$, $$%v$$, $$%v$$, %v, %v, %v),",
		loader.referentialSlug,
		record[1],
		parent,
		referent,
		record[4],
		record[5],
		record[6],
		record[7],
		record[8],
		record[9],
		collectedAlways,
		collectChildren,
		collectGeneralMessages,
	)
	loader.stopAreas = append(loader.stopAreas, values...)
	loader.bulkCounter[STOP_AREA]++

	if loader.bulkCounter[STOP_AREA] >= config.Config.LoadMaxInsert {
		loader.insertStopAreas()
	}

	return nil
}

func (loader *Loader) insertStopAreas() {
	if len(loader.stopAreas) == 0 {
		return
	}

	defer func() {
		loader.stopAreas = []byte{}
		loader.bulkCounter[STOP_AREA] = 0
	}()

	query := fmt.Sprintf("INSERT INTO stop_areas(referential_slug, id, parent_id, referent_id, model_name, name, object_ids, line_ids, attributes, siri_references, collected_always, collect_children, collect_general_messages) VALUES %v;",
		string(loader.stopAreas[:len(loader.stopAreas)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		logger.Log.Debugf("Error while inserting stopAreas: %v", err)
		fmt.Printf("Error while inserting stopAreas: %v", err)
		loader.errors++
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		logger.Log.Debugf("Unexpected error while inserting stopAreas: %v", err)
		fmt.Printf("Unexpected error while inserting stopAreas: %v", err)
		loader.errors++
		return
	}

	loader.insertedCount[STOP_AREA] += rows
}

func (loader *Loader) handleLine(record []string) error {
	if len(record) != 8 {
		return fmt.Errorf("wrong number of entries, expected 8 got %v", len(record))
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

	err = loader.handleForce(LINE, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,%v),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
		record[4],
		record[5],
		record[6],
		collectGeneralMessages,
	)
	loader.lines = append(loader.lines, values...)
	loader.bulkCounter[LINE]++

	if loader.bulkCounter[LINE] >= config.Config.LoadMaxInsert {
		loader.insertLines()
	}

	return nil
}

func (loader *Loader) insertLines() {
	if len(loader.lines) == 0 {
		return
	}

	defer func() {
		loader.lines = []byte{}
		loader.bulkCounter[LINE] = 0
	}()

	query := fmt.Sprintf("INSERT INTO lines(referential_slug,id,model_name,name,object_ids,attributes,siri_references,collect_general_messages) VALUES %v;", string(loader.lines[:len(loader.lines)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		logger.Log.Debugf("Error while inserting lines: %v", err)
		fmt.Printf("Error while inserting lines: %v", err)
		loader.errors++
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		logger.Log.Debugf("Unexpected error while inserting lines: %v", err)
		fmt.Printf("Unexpected error while inserting lines: %v", err)
		loader.errors++
		return
	}

	loader.insertedCount[LINE] += rows
}

func (loader *Loader) handleVehicleJourney(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("wrong number of entries, expected 10 got %v", len(record))
	}

	err := loader.handleForce(VEHICLE_JOURNEY, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
		record[4],
		record[5],
		record[6],
		record[7],
		record[8],
		record[9],
	)
	loader.vehicleJourneys = append(loader.vehicleJourneys, values...)
	loader.bulkCounter[VEHICLE_JOURNEY]++

	if loader.bulkCounter[VEHICLE_JOURNEY] >= config.Config.LoadMaxInsert {
		loader.insertVehicleJourneys()
	}

	return nil
}

func (loader *Loader) insertVehicleJourneys() {
	if len(loader.vehicleJourneys) == 0 {
		return
	}

	defer func() {
		loader.vehicleJourneys = []byte{}
		loader.bulkCounter[VEHICLE_JOURNEY] = 0
	}()

	query := fmt.Sprintf("INSERT INTO vehicle_journeys(referential_slug,id,model_name,name,object_ids,line_id,origin_name,destination_name,attributes,siri_references) VALUES %v;", string(loader.vehicleJourneys[:len(loader.vehicleJourneys)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		logger.Log.Debugf("Error while inserting vehicleJourneys: %v", err)
		fmt.Printf("Error while inserting vehicleJourneys: %v", err)
		loader.errors++
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		logger.Log.Debugf("Unexpected error while inserting vehicleJourneys: %v", err)
		fmt.Printf("Unexpected error while inserting vehicleJourneys: %v", err)
		loader.errors++
		return
	}

	loader.insertedCount[VEHICLE_JOURNEY] += rows
}

func (loader *Loader) handleStopVisit(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("wrong number of entries, expected 10 got %v", len(record))
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

	err = loader.handleForce(STOP_VISIT, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
		record[4],
		record[5],
		passageOrder,
		record[7],
		record[8],
		record[9],
	)
	loader.stopVisits = append(loader.stopVisits, values...)
	loader.bulkCounter[STOP_VISIT]++

	if loader.bulkCounter[STOP_VISIT] >= config.Config.LoadMaxInsert {
		loader.insertStopVisits()
	}

	return nil
}

func (loader *Loader) insertStopVisits() {
	if len(loader.stopVisits) == 0 {
		return
	}

	defer func() {
		loader.stopVisits = []byte{}
		loader.bulkCounter[STOP_VISIT] = 0
	}()

	query := fmt.Sprintf("INSERT INTO stop_visits(referential_slug,id,model_name,object_ids,stop_area_id,vehicle_journey_id,passage_order,schedules,attributes,siri_references) VALUES %v;", string(loader.stopVisits[:len(loader.stopVisits)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		logger.Log.Debugf("Error while inserting stopVisits: %v", err)
		fmt.Printf("Error while inserting stopVisits: %v", err)
		loader.errors++
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		logger.Log.Debugf("Unexpected error while inserting stopVisits: %v", err)
		fmt.Printf("Unexpected error while inserting stopVisits: %v", err)
		loader.errors++
		return
	}

	loader.insertedCount[STOP_VISIT] += rows
}
