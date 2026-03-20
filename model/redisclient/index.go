package redisclient

import (
	"fmt"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/redis/go-redis/v9"
)

func (rc *client) initIndexes() (err error) {
	logger.Log.Debugf("Creating Line indexes")
	err = rc.createIndex(Line, []string{ByReferentID}, true)
	if err != nil {
		return err
	}
	logger.Log.Debugf("Line index done")

	return err
}

func (rc *client) createIndex(model string, indexes []string, indexcodespaces bool) (err error) {
	indexKey := rc.prefixIndex(model)

	schema := []*redis.FieldSchema{
		{
			FieldName: "$.id",
			As:        "id",
			FieldType: redis.SearchFieldTypeTag,
		},
	}

	for _, i := range indexes {
		schema = append(schema, &redis.FieldSchema{
			FieldName: fmt.Sprintf("$.%s", i),
			As:        i,
			FieldType: redis.SearchFieldTypeTag,
		})
	}

	if indexcodespaces {
		for _, i := range rc.codespaces {
			schema = append(schema, &redis.FieldSchema{
				FieldName: fmt.Sprintf("$.codes.%v.value", i),
				As:        fmt.Sprintf("codespace_%v", i),
				FieldType: redis.SearchFieldTypeTag,
			})
		}
	}

	_, err = rc.c.FTCreate(
		rc.ctx,
		indexKey,
		// Options:
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []any{rc.prefix(model, ":")},
		},
		schema...,
	).Result()
	if err != nil {
		logger.Log.Debugf("Error while creating %s indexes: %v", model, err)
	} else {
		logger.Log.Debugf("%v Indexes created", model)
	}

	return err
}

func (rc *client) prefixIndex(modelName string) string {
	b := strings.Builder{}
	b.Grow(60)
	b.WriteString(rc.p)
	b.WriteString("indexes:")
	b.WriteString(modelName)
	return b.String()
}
