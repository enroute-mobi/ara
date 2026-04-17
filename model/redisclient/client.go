package redisclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/redis/go-redis/v9"
)

var TestClient Client

type Client interface {
	SetSlug(slug string)
	SetCodespaces(codeSpaces []string)
	Start(startTime time.Time) error

	Set(modelName string, record model) error
	Get(modelName, id string, attrs ...string) (string, error)
	FindAll(modelName string) ([]redis.Document, error)
	FindBy(modelName, fieldName, value string, attrs ...string) ([]redis.Document, error)
	FindByCode(modelName, codeSpace, id string) ([]redis.Document, error)
	Delete(modelName, id string) error
	FlushAll()
}

type client struct {
	c   *redis.Client
	ctx context.Context

	slug       string
	p          string
	codespaces []string
}

// New returns a new Client configured with the referential Slug and Codespaces
// to create its prefix. We can pass a context (optional).
// The Client is started with the default configuration, then tries to ping the
// Redis database and Panic if it can't do it
func New(slug string, codespaces []string, ctxs ...context.Context) (newClient Client, err error) {
	var ctx context.Context
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}

	c := redis.NewClient(&redis.Options{
		Addr:        config.Config.RedisAddr,
		Password:    config.Config.RedisPassword,
		DB:          config.Config.RedisDB,
		Protocol:    2,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 5 * time.Second,
	})

	err = c.Ping(ctx).Err()
	if err != nil {
		logger.Log.Panicf("Can't Ping Redis database: %v", err)
		// Return is useless but maybe we won't want to Panic in the future
		return
	}

	logger.Log.Debugf("Connected to Redis on %s", c.Options().Addr)

	newClient = &client{
		c:          c,
		ctx:        ctx,
		slug:       slug,
		codespaces: codespaces,
	}

	return
}

func (rc *client) SetSlug(slug string) {
	rc.slug = slug
}

func (rc *client) SetCodespaces(codeSpaces []string) {
	rc.codespaces = codeSpaces
}

// Start will build the new Client prefix with its Slug and the given StartTime.
// If the prefix hasn't changed we don't do anything, otherwise we initialize
// all the indexes.
func (rc *client) Start(startTime time.Time) error {
	newPrefix := fmt.Sprintf("%s:%v:", rc.slug, startTime.UnixNano())
	if rc.p == newPrefix {
		return nil
	}
	rc.p = newPrefix
	return rc.initIndexes()
}

// FlushAll is a method meant for specs that flush all the Database
func (rc *client) FlushAll() {
	rc.c.FlushAll(rc.ctx)
}

// prefix is used to prefix all indexes and records with:
// "[referential slug]:[referential startedTime]:""
//
// An idiomatic way to optimize builders is to reuse them. But as we can use them in a
// lot of threads, we would need to use a pool and maybe that's overkill for now.
func (rc *client) prefix(s ...string) string {
	b := strings.Builder{}
	b.Grow(60)
	b.WriteString(rc.p)
	for i := range s {
		b.WriteString(s[i])
	}
	return b.String()
}
