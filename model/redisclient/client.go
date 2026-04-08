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

// var replacer = strings.NewReplacer(
// 	"$", "\\$",
// 	"{", "\\{",
// 	"}", "\\}",
// 	"\\", "\\\\",
// 	"|", "\\|",
// 	":", "\\:",
// )

type Client interface {
	Start(t time.Time) error
	Stop()
	Started() bool

	Set(string, model) error
	Get(string, string) (string, error)
	GetPath(string, string, ...string) (string, error)
	FindAll(string) ([]redis.Document, error)
	FindBy(string, string, string) ([]redis.Document, error)
	FindByCode(string, string, string) ([]redis.Document, error)
	Delete(string, string) error
	FlushAll()
}

type client struct {
	c   *redis.Client
	ctx context.Context

	replacer   *strings.Replacer
	slug       string
	p          string
	codespaces []string
	started    bool
}

func New(slug string, codespaces []string, ctxs ...context.Context) (rc Client, err error) {
	var ctx context.Context
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}

	rc = &client{
		ctx:        ctx,
		slug:       slug,
		codespaces: codespaces,
	}

	return
}

func (rc *client) Start(t time.Time) error {
	rc.p = fmt.Sprintf("%s:%v:", rc.slug, t.UnixNano())

	var err error
	rc.c, err = newRedisclient(rc.ctx)
	if err != nil {
		return err
	}
	rc.started = true
	return rc.initIndexes()
}

func newRedisclient(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        config.Config.RedisAddr,
		Password:    config.Config.RedisPassword,
		DB:          config.Config.RedisDB,
		Protocol:    2,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 5 * time.Second,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		logger.Log.Panicf("Can't Ping Redis database: %v", err)
		return nil, err
	}

	logger.Log.Debugf("Connected to Redis on %s", client.Options().Addr)
	return client, nil

}

func (rc *client) Stop() {
	rc.started = false
	if rc.c != nil {
		rc.c.Close()
	}
}

func (rc *client) Started() bool {
	return rc.started
}

func (rc *client) FlushAll() {
	rc.c.FlushAll(rc.ctx)
}

func (rc *client) prefix(s ...string) string {
	b := strings.Builder{}
	b.Grow(60)
	b.WriteString(rc.p)
	for i := range s {
		b.WriteString(s[i])
	}
	return b.String()
}
