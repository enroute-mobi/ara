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

/*
	CSV Structure

operator,Id,ModelName,Name,Codes
stop_area,Id,ParentId,ReferentId,ModelName,Name,Codes,LineIds,Attributes,References,CollectedAlways,CollectChildren,CollectSituations
line,Id,ModelName,Name,Codes,Attributes,References,CollectSituations
vehicle_journey,Id,ModelName,Name,Codes,LineId,OriginName,DestinationName,Attributes,References,DirectionType, Number
stop_visit,Id,ModelName,Codes,StopAreaId,VehicleJourneyId,PassageOrder,Schedules,Attributes,References
stop_area_group,Id,ModelName,Name,ShortName,StopAreaIds
line_group,Id,ModelName,Name,ShortName,LineIds
facility,Id,ModelName,Codes

Comments are '#'
Separators are ',' leading spaces are trimed

Escape quotes with another quote ex: "[""1234"",""5678""]"
*/
const (
	STOP_AREA       = "stop_area"
	STOP_AREA_GROUP = "stop_area_group"
	LINE            = "line"
	LINE_GROUP      = "line_group"
	VEHICLE_JOURNEY = "vehicle_journey"
	STOP_VISIT      = "stop_visit"
	OPERATOR        = "operator"
	FACILITY        = "facility"
	TOTAL_INSERTS   = "Total"
	ERRORS          = "Errors"
)

type Loader struct {
	result          Result
	deletedModels   map[string]map[string]struct{}
	bulkCounter     map[string]int
	referentialSlug string
	vehicleJourneys []byte
	stopVisits      []byte
	stopAreas       []byte
	stopAreaGroups  []byte
	lines           []byte
	lineGroups      []byte
	operators       []byte
	facilities      []byte
	force           bool
	printErrors     bool
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

	logger.Log.Debug(result.PrintResult())
	fmt.Println(result.PrintResult())

	return nil
}

func NewLoader(referentialSlug string, force, printErrors bool) *Loader {
	d := make(map[string]map[string]struct{})
	for _, m := range [8]string{
		STOP_AREA,
		STOP_AREA_GROUP,
		LINE,
		LINE_GROUP,
		VEHICLE_JOURNEY,
		STOP_VISIT,
		OPERATOR,
		FACILITY} {
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
		case STOP_AREA_GROUP:
			err := loader.handleStopAreaGroup(record)
			if err != nil {
				loader.err(i, err)
			}
		case LINE:
			err := loader.handleLine(record)
			if err != nil {
				loader.err(i, err)
			}
		case LINE_GROUP:
			err := loader.handleLineGroup(record)
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
		case FACILITY:
			err := loader.handleFacility(record)
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
	loader.insertStopAreaGroups()
	loader.insertLines()
	loader.insertLineGroups()
	loader.insertVehicleJourneys()
	loader.insertStopVisits()
	loader.insertFacilities()

	loader.result.setTotalInserts()

	logger.Log.Printf("Load operation done in %v", time.Since(startTime))
	logger.Log.Print(loader.result.PrintResult())

	return loader.result
}

func (loader *Loader) handleForce(klass, modelName string) error {
	if loader.force {
		if _, ok := loader.deletedModels[klass][modelName]; !ok {
			loader.deletedModels[klass][modelName] = struct{}{}
			var araModel string
			if klass == FACILITY {
				araModel = "facilities"
			} else {
				araModel = fmt.Sprintf("%ss", klass)
			}
			query := fmt.Sprintf("delete from %v where model_name='%v' and referential_slug='%v';", araModel, modelName, loader.referentialSlug)
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

	query := fmt.Sprintf("INSERT INTO operators(referential_slug,id,model_name,name,codes) VALUES %v;", string(loader.operators[:len(loader.operators)-1]))
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

	var collectSituations bool
	if record[12] != "" {
		collectSituations, err = strconv.ParseBool(record[12])
		if err != nil {
			parseErrors.Add("CollectSituations", err)
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
		collectSituations,
	)
	loader.stopAreas = append(loader.stopAreas, values...)
	loader.bulkCounter[STOP_AREA]++

	if loader.bulkCounter[STOP_AREA] >= config.Config.LoadMaxInsert {
		loader.insertStopAreas()
	}

	return nil
}

func (loader *Loader) handleStopAreaGroup(record []string) error {
	if len(record) != 6 {
		return fmt.Errorf("wrong number of entries, expected 6 got %v", len(record))
	}

	parseErrors := ComplexError{}
	if parseErrors.ErrorCount() != 0 {
		return parseErrors
	}

	err := loader.handleForce(STOP_AREA_GROUP, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$, $$%v$$, $$%v$$, $$%v$$, $$%v$$, $$%v$$),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
		record[4],
		record[5],
	)

	loader.stopAreaGroups = append(loader.stopAreaGroups, values...)
	loader.bulkCounter[STOP_AREA_GROUP]++

	if loader.bulkCounter[STOP_AREA_GROUP] >= config.Config.LoadMaxInsert {
		loader.insertStopAreaGroups()
	}

	return nil
}

func (loader *Loader) insertStopAreaGroups() {
	if len(loader.stopAreaGroups) == 0 {
		return
	}

	defer func() {
		loader.stopAreaGroups = []byte{}
		loader.bulkCounter[STOP_AREA_GROUP] = 0
	}()

	query := fmt.Sprintf("INSERT INTO stop_area_groups(referential_slug, id, model_name, name, short_name, stop_area_ids) VALUES %v;",
		string(loader.stopAreaGroups[:len(loader.stopAreaGroups)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		loader.errInsert("stopAreaGroups", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("stopAreaGroups", err)
		return
	}

	loader.result.Import[STOP_AREA_GROUP] += rows
}

func (loader *Loader) insertStopAreas() {
	if len(loader.stopAreas) == 0 {
		return
	}

	defer func() {
		loader.stopAreas = []byte{}
		loader.bulkCounter[STOP_AREA] = 0
	}()

	query := fmt.Sprintf("INSERT INTO stop_areas(referential_slug, id, parent_id, referent_id, model_name, name, codes, line_ids, attributes, siri_references, collected_always, collect_children, collect_situations) VALUES %v;",
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

	var collectSituations bool
	if record[7] != "" {
		collectSituations, err = strconv.ParseBool(record[7])
		if err != nil {
			parseErrors.Add("CollectSituations", err)
		}
	}

	if len(record) >= 9 {
		number = record[8]
	}

	var referent string
	if len(record) == 10 && record[9] != "" {
		referent = fmt.Sprintf("$$%v$$", record[9])
	} else {
		referent = "null"
	}

	if parseErrors.ErrorCount() != 0 {
		return parseErrors
	}

	err = loader.handleForce(LINE, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,%v,$$%v$$,%v),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
		record[4],
		record[5],
		record[6],
		collectSituations,
		number,
		referent,
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

	query := fmt.Sprintf("INSERT INTO lines(referential_slug,id,model_name,name,codes,attributes,siri_references,collect_situations, number, referent_id) VALUES %v;", string(loader.lines[:len(loader.lines)-1]))
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

func (loader *Loader) handleLineGroup(record []string) error {
	if len(record) != 6 {
		return fmt.Errorf("wrong number of entries, expected 6 got %v", len(record))
	}

	parseErrors := ComplexError{}
	if parseErrors.ErrorCount() != 0 {
		return parseErrors
	}

	err := loader.handleForce(LINE_GROUP, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$, $$%v$$, $$%v$$, $$%v$$, $$%v$$, $$%v$$),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
		record[4],
		record[5],
	)

	loader.lineGroups = append(loader.lineGroups, values...)
	loader.bulkCounter[LINE_GROUP]++

	if loader.bulkCounter[LINE_GROUP] >= config.Config.LoadMaxInsert {
		loader.insertLineGroups()
	}

	return nil
}

func (loader *Loader) insertLineGroups() {
	if len(loader.lineGroups) == 0 {
		return
	}

	defer func() {
		loader.lineGroups = []byte{}
		loader.bulkCounter[LINE_GROUP] = 0
	}()

	query := fmt.Sprintf("INSERT INTO line_groups(referential_slug, id, model_name, name, short_name, line_ids) VALUES %v;",
		string(loader.lineGroups[:len(loader.lineGroups)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		loader.errInsert("lineGroups", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("lineGroups", err)
		return
	}

	loader.result.Import[LINE_GROUP] += rows
}

func (loader *Loader) handleVehicleJourney(record []string) error {
	if len(record) < 11 {
		return fmt.Errorf("wrong number of entries, expected 11 minimun got %v", len(record))
	}

	err := loader.handleForce(VEHICLE_JOURNEY, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$,$$%v$$),",
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
		record[10],
		record[11],
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

	query := fmt.Sprintf("INSERT INTO vehicle_journeys(referential_slug,id,model_name,name,codes,line_id,origin_name,destination_name,attributes,siri_references, direction_type, aimed_stop_visit_count) VALUES %v;", string(loader.vehicleJourneys[:len(loader.vehicleJourneys)-1]))

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

	query := fmt.Sprintf("INSERT INTO stop_visits(referential_slug,id,model_name,codes,stop_area_id,vehicle_journey_id,passage_order,schedules,attributes,siri_references) VALUES %v;", string(loader.stopVisits[:len(loader.stopVisits)-1]))
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

func (loader *Loader) handleFacility(record []string) error {
	if len(record) < 4 {
		return fmt.Errorf("wrong number of entries, expected at least 4 got %v", len(record))
	}

	var err error
	parseErrors := ComplexError{}

	if parseErrors.ErrorCount() != 0 {
		return parseErrors
	}

	err = loader.handleForce(FACILITY, record[2])
	if err != nil {
		return err
	}

	values := fmt.Sprintf("($$%v$$,$$%v$$,$$%v$$,$$%v$$),",
		loader.referentialSlug,
		record[1],
		record[2],
		record[3],
	)
	loader.facilities = append(loader.facilities, values...)
	loader.bulkCounter[FACILITY]++

	if loader.bulkCounter[FACILITY] >= config.Config.LoadMaxInsert {
		loader.insertFacilities()
	}

	return nil
}

func (loader *Loader) insertFacilities() {

	if len(loader.facilities) == 0 {
		return
	}

	defer func() {
		loader.facilities = []byte{}
		loader.bulkCounter[FACILITY] = 0
	}()

	query := fmt.Sprintf("INSERT INTO facilities(referential_slug,id,model_name,codes) VALUES %v;", string(loader.facilities[:len(loader.facilities)-1]))
	result, err := Database.Exec(query)
	if err != nil {
		loader.errInsert("facilities", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil { // should not happen
		loader.errInsert("facilities", err)
		return
	}
	loader.result.Import[FACILITY] += rows
}

func (loader *Loader) err(i int, e error) {
	if loader.printErrors {
		logger.Log.Debugf("Error on line %v: %v", i, e)
		fmt.Printf("Error on line %v: %v\n", i, e)
	}
	loader.result.Import[ERRORS]++

	if cerr, ok := e.(ComplexError); ok {
		loader.result.Errors[fmt.Sprint("Error on line ", i)] = append(loader.result.Errors[fmt.Sprint("Error on line ", i)], cerr.Errors...)
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
	for _, model := range [8]string{
		STOP_AREA,
		STOP_AREA_GROUP,
		LINE,
		LINE_GROUP,
		VEHICLE_JOURNEY,
		STOP_VISIT,
		OPERATOR,
		FACILITY} {
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
  %v StopAreaGroups
  %v Lines
  %v LineGroups
  %v VehicleJourneys
  %v StopVisits
  %v Facilities`,
		r.Import[ERRORS],
		r.Import[OPERATOR],
		r.Import[STOP_AREA],
		r.Import[STOP_AREA_GROUP],
		r.Import[LINE],
		r.Import[LINE_GROUP],
		r.Import[VEHICLE_JOURNEY],
		r.Import[STOP_VISIT],
		r.Import[FACILITY])
}
