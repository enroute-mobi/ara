package api

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"bitbucket.org/enroute-mobi/ara/core"
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
		model.Line |
		core.Partner
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

// newNormalizedMatcher validates a search string (at least 3 characters) and
// returns a matcher reporting whether a candidate string contains it. Both the
// search string and the candidates are compared accent-insensitively (combining
// marks stripped) and case-insensitively. label names the searched field and is
// only used to make the length error message explicit.
func newNormalizedMatcher(label, value string) (func(string) (bool, error), error) {
	if len(value) < 3 {
		return nil, fmt.Errorf("length of search %s must be at least 3 characters, got: %s", label, value)
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalizedValue, _, err := transform.String(t, value)
	if err != nil {
		return nil, fmt.Errorf("query parameter %q %s: cannot normalize:, %v", label, value, err.Error())
	}

	searchPattern, err := regexp.Compile("(?i)" + normalizedValue)
	if err != nil {
		return nil, fmt.Errorf("cannot create search pattern: %v", err.Error())
	}

	return func(candidate string) (bool, error) {
		normalizedCandidate, _, err := transform.String(t, candidate)
		if err != nil {
			return false, fmt.Errorf("cannot normalize %q: %v", candidate, err.Error())
		}
		return searchPattern.MatchString(normalizedCandidate), nil
	}, nil
}

func searchByName[S SearchableByName, M ModelForName[S]](s []M, params url.Values) ([]M, error) {
	searchName := params.Get("name")
	if searchName == "" {
		return s, nil
	}
	params.Del("name")

	matches, err := newNormalizedMatcher("name", searchName)
	if err != nil {
		return nil, err
	}

	possibleModels := []M{}
	for i := range s {
		ok, err := matches(s[i].GetName())
		if err != nil {
			return nil, err
		}
		if ok {
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

	searchCodeSpace, searchValue, found := strings.Cut(searchCode, ":")
	if !found {
		return nil, fmt.Errorf("invalid request: query parameter \"code\" : %s", searchCode)
	}
	if searchCodeSpace == "" || searchValue == "" {
		return nil, fmt.Errorf("code space or value should not be empty")
	}

	matches, err := newNormalizedMatcher("value", searchValue)
	if err != nil {
		return nil, err
	}

	possibleModels := []M{}
	for i := range s {
		code, ok := s[i].Code(searchCodeSpace)
		if !ok {
			continue
		}
		match, err := matches(code.Value())
		if err != nil {
			return nil, err
		}
		if match {
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

	matches, err := newNormalizedMatcher("text", searchText)
	if err != nil {
		return nil, err
	}

	possibleModels := []M{}
	for i := range s {
		for _, text := range s[i].SearchableText() {
			ok, err := matches(text)
			if err != nil {
				return nil, err
			}
			if ok {
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
