package redisclient

import (
	"strings"

	"github.com/redis/go-redis/v9"
)

// The model interface ensure we can get an Id from a record to use it as a key
type model interface {
	ModelId() string
}

// Set the JSON of a record in a key:
// "[referential slug]:[referential startedTime]:[model name]:[model id]"
func (rc *client) Set(modelName string, record model) error {
	_, err := rc.c.JSONSet(rc.ctx, rc.prefix(modelName, ":", record.ModelId()), "$", record).Result()
	return err
}

// Get finds a record by its Id. If no attrs are given, it returns a string
// containing an array of json marshaled objects. Otherwise, it returns a string
// containing a map[string][]any (map[attribute name][]value)
func (rc *client) Get(modelName, id string, attrs ...string) (string, error) {
	if len(attrs) == 0 {
		return rc.c.JSONGet(rc.ctx, rc.prefix(modelName, ":", id), "$").Result()
	}
	return rc.c.JSONGet(rc.ctx, rc.prefix(modelName, ":", id), attrs...).Result()
}

// FindAll returns all records for a given model in the form of an array of redis.Document
// We should only ever get one document, and the json marshaled object is in Document.Fields["$"]
func (rc *client) FindAll(modelName string) ([]redis.Document, error) {
	return rc.findBy(modelName, "*", &redis.FTSearchOptions{})
}

// FindByCode finds a record by a codeSpace/value pair. It returns an array of redis.Document
// We should only ever get one document, and the json marshaled object is in Document.Fields["$"]
func (rc *client) FindByCode(modelName, codespace, id string) ([]redis.Document, error) {
	return rc.findBy(modelName, codeQuery(codespace, id), &redis.FTSearchOptions{})
}

// FindBy finds one or more record by an attribute on which we have an index.
// It returns an array of redis.Document.
// If no attrs are given, the json marshaled object is in Document.Fields["$"]
// otherwise all attributes will be in Document.Fields, a map[string]string (map[attribute]value)
func (rc *client) FindBy(modelName, fieldName, value string, attrs ...string) ([]redis.Document, error) {
	if len(attrs) == 0 {
		return rc.findBy(modelName, query(fieldName, value), &redis.FTSearchOptions{})
	}

	opts := &redis.FTSearchOptions{
		Return: []redis.FTSearchReturn{},
	}
	b := strings.Builder{}
	for i := range attrs {
		b.Reset()
		b.WriteString("$.")
		b.WriteString(attrs[i])
		opts.Return = append(opts.Return, redis.FTSearchReturn{FieldName: b.String(), As: attrs[i]})
	}
	return rc.findBy(modelName, query(fieldName, value), opts)
}

// findBy performs a Redis FT.SEARCH. It returns an array of redis.Document.
// We can get the requested attributes in Document.Fields, a map[string]string
// Warning: If we don't add an FTSearchOption, even empty, we aren't in Dialect 2
// And we need it for the query syntax
func (rc *client) findBy(modelName, query string, opts *redis.FTSearchOptions) ([]redis.Document, error) {
	cmd := rc.c.FTSearchWithArgs(
		rc.ctx,
		rc.prefixIndex(modelName),
		query,
		opts,
	)

	r, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return r.Docs, nil
}

func (rc *client) Delete(modelName, id string) error {
	_, err := rc.c.JSONDel(rc.ctx, rc.prefix(modelName, ":", id), "$").Result()
	return err
}

// If we have the '"' char in tags, we will need to use a replacer for queries.
// I don't think that's currently possible, but we should ensure that we can't do it in ara-front
// var replacer = strings.NewReplacer(
// 	"$", "\\$",
// 	"{", "\\{",
// 	"}", "\\}",
// 	"\\", "\\\\",
// 	"|", "\\|",
// 	":", "\\:",
// )

// query builds a query for FindBy requests
func query(fieldName, id string) string {
	b := strings.Builder{}
	b.WriteString("@")
	b.WriteString(fieldName)
	b.WriteString(":{\"")
	// replacer.WriteString(&b, id)
	b.WriteString(id)
	b.WriteString("\"}")
	return b.String()
}

// codeQuery builds a query for FindByCode requests
func codeQuery(codespace, id string) string {
	b := strings.Builder{}
	b.WriteString("@codespace_")
	b.WriteString(codespace)
	b.WriteString(":{\"")
	// replacer.WriteString(&b, id)
	b.WriteString(id)
	b.WriteString("\"}")
	return b.String()
}

// ------------------------------------------------------------------------------------------- \\

// It isn't used for now but could prove mandatory in the future.
// Warning: not properly tested

const (
	maxBatch = 200_000
)

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
