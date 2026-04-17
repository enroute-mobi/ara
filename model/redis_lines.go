package model

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/enroute-mobi/ara/model/redisclient"
)

type redisLines struct {
	redisCodeHandlerManager[LineId, Line, *Line]
}

func NewRedisLines(client redisclient.Client) Lines {
	return &redisLines{
		redisCodeHandlerManager: redisCodeHandlerManager[LineId, Line, *Line]{
			redisManager: redisManager[LineId, Line, *Line]{
				new:       NewLine,
				modelName: redisclient.Line,
				client:    client,
			},
		},
	}
}

func (manager *redisLines) FindByReferentId(id LineId) (lines []*Line) {
	return manager.FindAllBy(redisclient.ByReferentID, string(id))
}

func (manager *redisLines) FindFamily(lineId LineId) (lineIds []LineId) {
	lineIds = []LineId{lineId}

	ids := manager.FindAllAttributesBy(redisclient.ByReferentID, string(lineId), redisclient.ModelID)
	for i := range ids {
		lineIds = append(lineIds, manager.FindFamily(LineId(ids[i][redisclient.ModelID]))...)
	}

	return
}

func (manager *redisLines) FindFamilyFromCode(code Code) (lineIds []LineId) {
	l, ok := manager.FindByCode(code)
	if ok {
		return manager.FindFamily(l.id)
	}

	return
}

func (manager *redisLines) Load(referentialSlug string) error {
	var selectLines []SelectLine
	modelDate := manager.model.Date()
	sqlQuery := fmt.Sprintf("select * from lines where referential_slug = '%s' and model_date = '%s'", referentialSlug, modelDate.String())
	_, err := Database.Select(&selectLines, sqlQuery)
	if err != nil {
		return err
	}
	for _, sl := range selectLines {
		line := manager.New()
		line.id = LineId(sl.Id)
		if sl.Name.Valid {
			line.Name = sl.Name.String
		}
		if sl.ReferentId.Valid {
			line.ReferentId = LineId(sl.ReferentId.String)
		}

		if sl.CollectSituations.Valid {
			line.CollectSituations = sl.CollectSituations.Bool
		}

		if sl.Number.Valid {
			line.Number = sl.Number.String
		}
		if sl.RawAttributes.Valid && len(sl.RawAttributes.String) > 0 {
			if err = json.Unmarshal([]byte(sl.RawAttributes.String), &line.RawAttributes); err != nil {
				return err
			}
		}

		if sl.References.Valid && len(sl.References.String) > 0 {
			references := make(map[string]Reference)
			if err = json.Unmarshal([]byte(sl.References.String), &references); err != nil {
				return err
			}
			line.References.SetReferences(references)
		}

		if sl.Codes.Valid && len(sl.Codes.String) > 0 {
			codeMap := make(map[string]string)
			if err = json.Unmarshal([]byte(sl.Codes.String), &codeMap); err != nil {
				return err
			}
			line.SetCodesFromMap(codeMap)
		}

		manager.Save(line)
	}
	return nil
}
