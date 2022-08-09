package model

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/logger"
)

/* CSV Structure

operator,Id,ModelName,Name,ObjectIDs
stop_area,Id,ParentId,ReferentId,ModelName,Name,ObjectIDs,LineIds,Attributes,References,CollectedAlways,CollectChildren,CollectGeneralMessages
line,Id,ModelName,Name,ObjectIDs,Attributes,References,CollectGeneralMessages
vehicle_journey,Id,ModelName,Name,ObjectIDs,LineId,OriginName,DestinationName,Attributes,References,DirectionType, Number
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
	TOTAL_INSERTS   = "Total"
	ERRORS          = "Errors"
)

type Loader struct {
	referentialSlug string
	force           bool
	printErrors     bool
	deletedModels   map[string]map[string]struct{}
	operators       []byte
	stopAreas       []byte
	lines           []byte
	vehicleJourneys []byte
	stopVisits      []byte
	bulkCounter     map[string]int
	result          Result
}

type Result struct {
	Import map[string]int64
	Errors map[string][]string
}

type ComplexError struct {
	Errors []string
}

func (c ComplexError) Error() string {
	return strings.Join(c.Errors, ", ")
}

func (c *ComplexError) Add(field string, err error) {
	c.Errors = append(c.Errors, fmt.Sprintf("%v: %v", field, err))
}

func (c ComplexError) ErrorCount() int {
	return len(c.Errors)
}

func LoadFromCSVFile(filePath string, referentialSlug string, force bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("loader error: error while opening file: %v", err)
	}
	defer file.Close()

	result := NewLoader(referentialSlug, force, true).Load(file)

	if result.TotalInserts() == 0 {
		if result.ErrorCount() == 0 {
			return fmt.Errorf("loader error: empty file")
		}
		return fmt.Errorf("loader error: couldn't import anything, import raised %v errors", result.ErrorCount())
	}

	logger.Log.Debugf(result.PrintResult())
	fmt.Println(result.PrintResult())

	return nil
}

func NewLoader(referentialSlug string, force, printErrors bool) *Loader {
	d := make(map[string]map[string]struct{})
	for _, m := range [5]string{STOP_AREA, LINE, VEHICLE_JOURNEY, STOP_VISIT, OPERATOR} {
		d[m] = make(map[string]struct{})
	}
	r := Result{
		Import: make(map[string]int64),
		Errors: make(map[string][]string),
	}
	return &Loader{
		referentialSlug: referentialSlug,
		force:           force,
		printErrors:     printErrors,
		deletedModels:   d,
		bulkCounter:     make(map[string]int),
		result:          r,
	}
}

func (loader Loader) Load(reader io.Reader) Result {
	// Config CSV reader
	csvReader := csv.NewReader(reader)
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = -1
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	startTime := time.Now()
	logger.Log.Printf("Load operation started at %v", startTime)

	var i int
	for {
		i++
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			loader.err(i, err)
			continue
		}

		switch record[0] {
		case OPERATOR:
			err := loader.handleOperator(record)
			if err != nil {
				loader.err(i, err)
			}
		case STOP_AREA:
			err := loader.handleStopArea(record)
			if err != nil {
				loader.err(i, err)
			}
		case LINE:
			err := loader.handleLine(record)
			if err != nil {
				loader.err(i, err)
			}
		case VEHICLE_JOURNEY:
			err := loader.handleVehicleJourney(record)
			if err != nil {
				loader.err(i, err)
			}
		case STOP_VISIT:
			err := loader.handleStopVisit(record)
			if err != nil {
				loader.err(i, err)
			}
		default:
			loader.err(i, fmt.Errorf("unknown record type %v", record[0]))
			continue
		}
	}

	loader.insertOperators()
	loader.insertStopAreas()
	loader.insertLines()
	loader.insertVehicleJourneys()
	loader.insertStopVisits()

	loader.result.setTotalInserts()

	logger.Log.Printf("Load operation done in %v", time.Since(startTime))
	logger.Log.Printf(loader.result.PrintResult())

	return loader.result
}

func (loader *Loader) handleForce(klass, modelName string) error {
	if loader.force {
		if _, ok := loader.deletedModels[klass][modelName]; !ok {
			loader.deletedModels[klass][modelName] = struct{}{}
			query := fmt.Sprintf("delete from %vs where model_name='%v' and referential_slug='%v';", klass, modelName, loader.referentialSlug)
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
		loader.errInsert("operators", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("operators", err)
		return
	}

	loader.result.Import[OPERATOR] += rows
}

func (loader *Loader) handleStopArea(record []string) error {
	if len(record) != 13 {
		return fmt.Errorf("wrong number of entries, expected 13 got %v", len(record))
	}

	var err error
	parseErrors := ComplexError{}

	var collectedAlways bool
	if record[10] != "" {
		collectedAlways, err = strconv.ParseBool(record[10])
		if err != nil {
			parseErrors.Add("CollectedAlways", err)
		}
	}

	var collectChildren bool
	if record[11] != "" {
		collectChildren, err = strconv.ParseBool(record[11])
		if err != nil {
			parseErrors.Add("CollectChildren", err)
		}
	}

	var collectGeneralMessages bool
	if record[12] != "" {
		collectGeneralMessages, err = strconv.ParseBool(record[12])
		if err != nil {
			parseErrors.Add("CollectGeneralMessages", err)
		}
	}

	if parseErrors.ErrorCount() != 0 {
		return parseErrors
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
		loader.errInsert("stopAreas", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("stopAreas", err)
		return
	}

	loader.result.Import[STOP_AREA] += rows
}

func (loader *Loader) handleLine(record []string) error {
	var number string

	if len(record) < 8 {
		return fmt.Errorf("wrong number of entries, expected at least 8 got %v", len(record))
	}

	var err error
	parseErrors := ComplexError{}

	var collectGeneralMessages bool
	if record[7] != "" {
		collectGeneralMessages, err = strconv.ParseBool(record[7])
		if err != nil {
			parseErrors.Add("CollectGeneralMessages", err)
		}
	}

	if len(record) == 9 {
		number = record[8]
	}

	if parseErrors.ErrorCount() != 0 {
		return parseErrors
	}

	err = loader.handleForce(LINE, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,%v,$$%v$$),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
		record[4],
		record[5],
		record[6],
		collectGeneralMessages,
		number,
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

	query := fmt.Sprintf("INSERT INTO lines(referential_slug,id,model_name,name,object_ids,attributes,siri_references,collect_general_messages, number) VALUES %v;", string(loader.lines[:len(loader.lines)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		loader.errInsert("lines", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("lines", err)
		return
	}

	loader.result.Import[LINE] += rows
}

func (loader *Loader) handleVehicleJourney(record []string) error {
	var directionType string

	if len(record) < 10 {
		return fmt.Errorf("wrong number of entries, expected 10 minimun got %v", len(record))
	}

	err := loader.handleForce(VEHICLE_JOURNEY, record[2])
	if err != nil {
		return err
	}

	if len(record) == 11 {
		directionType = record[10]
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$),",
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
		directionType,
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

	query := fmt.Sprintf("INSERT INTO vehicle_journeys(referential_slug,id,model_name,name,object_ids,line_id,origin_name,destination_name,attributes,siri_references, direction_type) VALUES %v;", string(loader.vehicleJourneys[:len(loader.vehicleJourneys)-1]))

	result, err := Database.Exec(query)
	if err != nil {
		loader.errInsert("vehicleJourneys", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("vehicleJourneys", err)
		return
	}

	loader.result.Import[VEHICLE_JOURNEY] += rows
}

func (loader *Loader) handleStopVisit(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("wrong number of entries, expected 10 got %v", len(record))
	}

	var err error
	parseErrors := ComplexError{}

	var passageOrder int
	if record[6] != "" {
		passageOrder, err = strconv.Atoi(record[6])
		if err != nil {
			parseErrors.Add("PassageOrder", err)
		}
	}

	if parseErrors.ErrorCount() != 0 {
		return parseErrors
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
		loader.errInsert("stopVisits", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("stopVisits", err)
		return
	}

	loader.result.Import[STOP_VISIT] += rows
}

func (loader *Loader) err(i int, e error) {
	if loader.printErrors {
		logger.Log.Debugf("Error on line %v: %v", i, e)
		fmt.Printf("Error on line %v: %v\n", i, e)
	}
	loader.result.Import[ERRORS]++

	if cerr, ok := e.(ComplexError); ok {
		for i := range cerr.Errors {
			loader.result.Errors[fmt.Sprint("Error on line ", i)] = append(loader.result.Errors[fmt.Sprint("Error on line ", i)], cerr.Errors[i])
		}
	} else {
		loader.result.Errors[fmt.Sprint("Error on line ", i)] = append(loader.result.Errors[fmt.Sprint("Error on line ", i)], e.Error())
	}
}

func (loader *Loader) errInsert(m string, e error) {
	if loader.printErrors {
		logger.Log.Debugf("Error while inserting %v: %v", m, e)
		fmt.Printf("Error while inserting %v: %v\n", m, e)
	}
	loader.result.Import[ERRORS]++
	loader.result.Errors[fmt.Sprint("Error while inserting ", m)] = append(loader.result.Errors[fmt.Sprint("Error while inserting ", m)], e.Error())
}

func (r *Result) setTotalInserts() {
	var c int64
	for _, model := range [5]string{STOP_AREA, LINE, VEHICLE_JOURNEY, STOP_VISIT, OPERATOR} {
		c += r.Import[model]
	}
	r.Import[TOTAL_INSERTS] = c
}

func (r Result) TotalInserts() int64 {
	return r.Import[TOTAL_INSERTS]
}

func (r Result) Inserted(m string) int64 {
	return r.Import[m]
}

func (r Result) ErrorCount() int64 {
	return r.Import[ERRORS]
}

func (r Result) PrintResult() string {
	return fmt.Sprintf(`Import successful. Import raised %v errors
  %v Operators
  %v StopAreas
  %v Lines
  %v VehicleJourneys
  %v StopVisits`, r.Import[ERRORS], r.Import[OPERATOR], r.Import[STOP_AREA], r.Import[LINE], r.Import[VEHICLE_JOURNEY], r.Import[STOP_VISIT])
}
