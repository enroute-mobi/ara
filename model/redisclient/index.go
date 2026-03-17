package redisclient

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/redis/go-redis/v9"
)

func (rc *Client) InitIndexes() (err error) {
	logger.Log.Debugf("Creating Line index")
	err = rc.createIndex("line", "referent_id")
	if err != nil {
		return err
	}
	err = rc.createIndex("line", "code", true)
	if err != nil {
		return err
	}
	logger.Log.Debugf("Line index done")

	return err
}

func (rc *Client) createIndex(model, index string, array ...bool) (err error) {
	indexKey := rc.prefixIndex(model, index)
	prefix := rc.prefix(model, ":")
	fieldName := fmt.Sprintf("$.%s", index)
	if len(array) != 0 {
		fieldName += "[*]"
	}

	_, err = rc.c.FTCreate(
		rc.ctx,
		indexKey,
		// Options:
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []any{prefix},
		},
		// Index schema fields:
		&redis.FieldSchema{
			FieldName: fmt.Sprintf("$.%s", index),
			As:        index,
			FieldType: redis.SearchFieldTypeTag,
		},
	).Result()
	if err != nil {
		logger.Log.Debugf("Error while creating %s index '%s': %v", model, index, err)
		return err
	}
	logger.Log.Debugf("%v Index '%v' created", model, index)

	return err
}
