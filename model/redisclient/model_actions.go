package redisclient

import (
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	maxBatch = 200_000
)

type model interface {
	ModelId() string
}

// --------------------------------- \\

type Batch struct {
	rc   *Client
	docs []redis.JSONSetArgs
}

func (rc *Client) NewBatch() Batch {
	return Batch{rc: rc}
}

func (b *Batch) Add(value model) {
	b.docs = append(b.docs, redis.JSONSetArgs{Key: value.ModelId(), Path: "$", Value: value})
}

func (b *Batch) Save(doc redis.JSONSetArgs) error {
	for i := 0; i < len(b.docs); i += maxBatch {
		j := i + maxBatch
		if j > len(b.docs) {
			j = len(b.docs)
		}

		_, err := b.rc.c.JSONMSetArgs(b.rc.ctx, b.docs[i:j]).Result()
		if err != nil {
			return err
		}
	}
	return nil

}

// --------------------------------- \\

func (rc *Client) Save(value model) error {
	_, err := rc.c.JSONSet(rc.ctx, value.ModelId(), "$", value).Result()
	return err
}

func (rc *Client) Get(id string) (string, error) {
	return rc.c.JSONGet(rc.ctx, rc.prefix(id), "$").Result()
}

func (rc *Client) FindBy(modelName, indexName, id string, uuid ...bool) ([]redis.Document, error) {
	if len(uuid) != 0 {
		id = strings.ReplaceAll(id, "-", "\\-")
	}

	r, err := rc.c.FTSearch(
		rc.ctx,
		rc.prefixIndex(modelName, indexName),
		query(indexName, id),
	).Result()
	if err != nil {
		return nil, err
	}

	return r.Docs, nil
}

func query(indexName, id string) string {
	b := strings.Builder{}
	b.WriteString("@")
	b.WriteString(indexName)
	b.WriteString(":{")
	b.WriteString(id)
	b.WriteRune('}')
	return b.String()
}

func (rc *Client) Delete(id string) error {
	_, err := rc.c.JSONDel(rc.ctx, rc.prefix(id), "$").Result()
	return err
}
