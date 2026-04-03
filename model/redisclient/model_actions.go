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
	rc   *client
	docs []redis.JSONSetArgs
}

func (rc *client) NewBatch() Batch {
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

func (rc *client) Set(modelName string, value model) error {
	_, err := rc.c.JSONSet(rc.ctx, rc.prefix(modelName, ":", value.ModelId()), "$", value).Result()
	return err
}

func (rc *client) Get(modelName, id string) (string, error) {
	return rc.c.JSONGet(rc.ctx, rc.prefix(modelName, ":", id), "$").Result()
}

func (rc *client) GetPath(modelName, id, path string) (string, error) {
	return rc.c.JSONGet(rc.ctx, rc.prefix(modelName, ":", id), path).Result()
}

func (rc *client) FindAll(modelName string) ([]redis.Document, error) {
	return rc.findBy(modelName, "*")
}

func (rc *client) FindByCode(modelName, codespace, id string) ([]redis.Document, error) {
	// We suppose we can't have any '"' or ':' in codespaces or ids
	return rc.findBy(modelName, codeQuery(codespace, id))
}

// https://redis.io/docs/latest/develop/ai/search-and-query/query/exact-match/#em3?lang=Go
func (rc *client) FindBy(modelName, fieldName, id string) ([]redis.Document, error) {
	// We suppose we can't have any '"' in ids
	// id = strings.ReplaceAll(id, "\"", "\\\"")
	return rc.findBy(modelName, query(fieldName, id))
}

func (rc *client) findBy(modelName, query string) ([]redis.Document, error) {
	r, err := rc.c.FTSearch(
		rc.ctx,
		rc.prefixIndex(modelName),
		query,
	).Result()
	if err != nil {
		return nil, err
	}

	return r.Docs, nil
}

func (rc *client) Delete(modelName, id string) error {
	_, err := rc.c.JSONDel(rc.ctx, rc.prefix(modelName, ":", id), "$").Result()
	return err
}

func query(fieldName, id string) string {
	b := strings.Builder{}
	b.WriteString("@")
	b.WriteString(fieldName)
	b.WriteString(":{\"")
	b.WriteString(id)
	b.WriteString("\"}")
	return b.String()
}

func codeQuery(codespace, id string) string {
	b := strings.Builder{}
	b.WriteString("@codespace_")
	b.WriteString(codespace)
	b.WriteString(":{\"")
	b.WriteString(id)
	b.WriteString("\"}")
	return b.String()
}
