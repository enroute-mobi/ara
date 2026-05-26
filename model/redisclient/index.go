package redisclient

import (
	"fmt"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/redis/go-redis/v9"
)

// initIndexes is meant to create Indexes for each model (Stops, Lines, Vehicles,...)
func (rc *client) initIndexes() (err error) {
	logger.Log.Debugf("Creating Line indexes")
	err = rc.createIndex(Line, []string{ByReferentID}, true)
	if err != nil {
		return err
	}
	logger.Log.Debugf("Line index done")

	return err
}

// createIndex creates a default index on Id, an index for each of the given parameters in the indexes slice,
// and an index for all the client code spaces if parameter indexCodespaces is true
func (rc *client) createIndex(modelName string, indexes []string, indexCodespaces bool) (err error) {
	indexName := rc.prefixIndex(modelName)

	schema := []*redis.FieldSchema{
		{
			FieldName: "$.Id",
			As:        "Id",
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

	if indexCodespaces {
		for _, i := range rc.codespaces {
			schema = append(schema, &redis.FieldSchema{
				FieldName: fmt.Sprintf("$.Codes.%v", i),
				// FieldName: fmt.Sprintf("$.codes.%v.Value", i),
				As:        fmt.Sprintf("codespace_%v", i),
				FieldType: redis.SearchFieldTypeTag,
			})
		}
	}

	_, err = rc.c.FTCreate(
		rc.ctx,
		indexName,
		// Options:
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []any{rc.prefix(modelName, ":")},
		},
		schema...,
	).Result()
	if err != nil {
		// For now we Panic as ara can't function at all without the indexes correctly set
		logger.Log.Panicf("Error while creating %s indexes: %v", modelName, err)
	} else {
		logger.Log.Debugf("%v Indexes created", modelName)
	}

	return err
}

// prefixIndex is used to prefix all indexes with:
// "[referential slug]:[referential startedTime]:indexes:[model name]"
func (rc *client) prefixIndex(modelName string) string {
	b := strings.Builder{}
	b.Grow(60)
	b.WriteString(rc.p)
	b.WriteString("indexes:")
	b.WriteString(modelName)
	return b.String()
}
