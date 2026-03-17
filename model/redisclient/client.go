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

type Client struct {
	c   *redis.Client
	ctx context.Context

	p string
}

func New(slug string, t time.Time, ctxs ...context.Context) (rc *Client, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:        config.Config.RedisAddr,
		Password:    config.Config.RedisPassword,
		DB:          config.Config.RedisDB,
		Protocol:    2,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 5 * time.Second,
	})
	var ctx context.Context
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}

	err = client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Connected to Redis on %s", client.Options().Addr)

	rc = &Client{
		c:   client,
		ctx: ctx,
		p:   fmt.Sprintf("%s:%v:", slug, t.UnixNano()),
	}

	rc.InitIndexes()

	return
}

func (rc *Client) Close() {
	rc.c.Close()
}

func (rc *Client) prefix(s ...string) string {
	b := strings.Builder{}
	b.Grow(60)
	b.WriteString(rc.p)
	for i := range s {
		b.WriteString(s[i])
	}
	return b.String()
}

func (rc *Client) prefixIndex(modelName, indexName string) string {
	b := strings.Builder{}
	b.Grow(60)
	b.WriteString(rc.p)
	b.WriteString("indexes:")
	b.WriteString(modelName)
	b.WriteString(":by_")
	b.WriteString(indexName)
	return b.String()
}
