package api

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type SearchableByCode interface {
	model.Line |
		model.Vehicle |
		model.VehicleJourney |
		model.StopArea |
		model.Situation
}

type ModelForCode[S SearchableByCode] interface {
	Code(string) (model.Code, bool)
	*S
}

type SearchableByName interface {
	model.StopArea |
		model.VehicleJourney |
		model.Line
}

type ModelForName[S SearchableByName] interface {
	GetName() string
	*S
}

type SearchableByLineId interface {
	model.Vehicle
}

type ModelForLineIds[S SearchableByLineId] interface {
	GetLineId() model.LineId
	*S
}

type SearchableByText interface {
	model.Situation
}

type ModelForText[S SearchableByText] interface {
	SearchableText() []string
	*S
}

type SearchableByAffectedLineIds interface {
	model.Situation
}

type ModelForAffectedLineIds[S SearchableByAffectedLineIds] interface {
	GetAffectedLineIds() []model.LineId
	*S
}

type SearchableByAffectedStopAreaIds interface {
	model.Situation
}

type ModelForAffectedStopAreaIds[S SearchableByAffectedStopAreaIds] interface {
	GetAffectedStopAreaIds() []model.StopAreaId
	*S
}

func searchByName[S SearchableByName, M ModelForName[S]](s []M, params url.Values) ([]M, error) {
	searchName := params.Get("name")
	if searchName == "" {
		return s, nil
	}
	params.Del("name")

	possibleModels := []M{}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	if len(searchName) < 3 {
		return nil, fmt.Errorf("length of search name must be at least 3 characters, got: %s", searchName)
	}

	normalizedSearchPattern, _, err := transform.String(t, searchName)
	if err != nil {
		return nil, fmt.Errorf("query parameter \"name\" %s: cannot normalize:, %v", searchName, err.Error())
	}

	searchPattern, err := regexp.Compile("(?i)" + normalizedSearchPattern)
	if err != nil {
		return nil, fmt.Errorf("cannot create search pattern: %v", err.Error())
	}

	for i := range s {
		normalizedSaName, _, err := transform.String(t, s[i].GetName())
		if err != nil {
			return nil, fmt.Errorf("cannot normalize stopArea name: %v", err.Error())
		}
		if searchPattern.MatchString(normalizedSaName) {
			possibleModels = append(possibleModels, s[i])
		}
	}
	return possibleModels, nil
}

func searchByCode[S SearchableByCode, M ModelForCode[S]](s []M, params url.Values) ([]M, error) {
	searchCode := params.Get("code")
	if searchCode == "" {
		return s, nil
	}
	params.Del("code")

	possibleModels := []M{}

	searchCodeSpace, searchValue, found := strings.Cut(searchCode, ":")
	if !found {
		return nil, fmt.Errorf("invalid request: query parameter \"code\" : %s", searchCode)
	}
	if searchCodeSpace == "" || searchValue == "" {
		return nil, fmt.Errorf("code space or value should not be empty")
	}
	if len(searchValue) < 3 {
		return nil, fmt.Errorf("length of search value must be at least 3 characters, got: %s", searchValue)
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalizedSearchValue, _, err := transform.String(t, searchValue)
	if err != nil {
		return nil, fmt.Errorf("invalid request: query parameter \"code\" %s: cannot normalize value:, %v", searchValue, err.Error())
	}

	searchPattern, err := regexp.Compile("(?i)" + normalizedSearchValue)
	if err != nil {
		return nil, fmt.Errorf("cannot create search pattern: %v", err.Error())
	}

	for i := range s {
		code, ok := s[i].Code(searchCodeSpace)
		if !ok {
			continue
		}
		normalisedCodeValue, _, err := transform.String(t, code.Value())
		if err != nil {
			return nil, fmt.Errorf("cannot normalize code value: %v", err.Error())
		}
		if searchPattern.MatchString(normalisedCodeValue) {
			possibleModels = append(possibleModels, s[i])
		}
	}
	return possibleModels, nil
}

func searchByLineIds[S SearchableByLineId, M ModelForLineIds[S]](s []M, params url.Values) ([]M, error) {
	lineIds := params["line_ids[]"]
	if len(lineIds) == 0 {
		return s, nil
	}

	params.Del("line_ids[]")

	for i := range lineIds {
		err := uuid.Validate(lineIds[i])
		if err != nil {
			return nil, fmt.Errorf("line id is not a valid UUID: %s", lineIds[i])
		}
	}

	possibleModels := []M{}
	for i := range s {
		if slices.Contains(lineIds, string(s[i].GetLineId())) {
			possibleModels = append(possibleModels, s[i])
		}
	}
	return possibleModels, nil
}

func searchByText[S SearchableByText, M ModelForText[S]](s []M, params url.Values) ([]M, error) {
	searchText := params.Get("text")
	if searchText == "" {
		return s, nil
	}
	params.Del("text")

	if len(searchText) < 3 {
		return nil, fmt.Errorf("length of search text must be at least 3 characters, got: %s", searchText)
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalizedSearchText, _, err := transform.String(t, searchText)
	if err != nil {
		return nil, fmt.Errorf("query parameter \"text\" %s: cannot normalize:, %v", searchText, err.Error())
	}

	searchPattern, err := regexp.Compile("(?i)" + normalizedSearchText)
	if err != nil {
		return nil, fmt.Errorf("cannot create search pattern: %v", err.Error())
	}

	possibleModels := []M{}
	for i := range s {
		for _, text := range s[i].SearchableText() {
			normalizedText, _, err := transform.String(t, text)
			if err != nil {
				return nil, fmt.Errorf("cannot normalize text: %v", err.Error())
			}
			if searchPattern.MatchString(normalizedText) {
				possibleModels = append(possibleModels, s[i])
				break
			}
		}
	}
	return possibleModels, nil
}

func searchByAffectedLineIds[S SearchableByAffectedLineIds, M ModelForAffectedLineIds[S]](s []M, params url.Values) ([]M, error) {
	lineIds := params["line_ids[]"]
	if len(lineIds) == 0 {
		return s, nil
	}

	params.Del("line_ids[]")

	for i := range lineIds {
		err := uuid.Validate(lineIds[i])
		if err != nil {
			return nil, fmt.Errorf("line id is not a valid UUID: %s", lineIds[i])
		}
	}

	possibleModels := []M{}
	for i := range s {
		for _, lineId := range s[i].GetAffectedLineIds() {
			if slices.Contains(lineIds, string(lineId)) {
				possibleModels = append(possibleModels, s[i])
				break
			}
		}
	}
	return possibleModels, nil
}

func searchByAffectedStopAreaIds[S SearchableByAffectedStopAreaIds, M ModelForAffectedStopAreaIds[S]](s []M, params url.Values) ([]M, error) {
	stopAreaIds := params["stop_area_ids[]"]
	if len(stopAreaIds) == 0 {
		return s, nil
	}

	params.Del("stop_area_ids[]")

	for i := range stopAreaIds {
		err := uuid.Validate(stopAreaIds[i])
		if err != nil {
			return nil, fmt.Errorf("stop area id is not a valid UUID: %s", stopAreaIds[i])
		}
	}

	possibleModels := []M{}
	for i := range s {
		for _, stopAreaId := range s[i].GetAffectedStopAreaIds() {
			if slices.Contains(stopAreaIds, string(stopAreaId)) {
				possibleModels = append(possibleModels, s[i])
				break
			}
		}
	}
	return possibleModels, nil
}
