package redis

import (
	"context"
	"strings"
	"time"

	errors "github.com/eclipse-xfsc/microservice-core-go/pkg/err"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb        redis.Cmdable
	defaultTTL time.Duration
}

func New(addr, user, pass string, db int, defaultTTL time.Duration, cluster bool) *Client {
	var rdb redis.Cmdable
	if cluster {
		nodes := strings.Split(addr, ";")
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          nodes,
			Username:       user,
			Password:       pass,
			RouteByLatency: true,
			DialTimeout:    10 * time.Second,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxRedirects:   10,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:         addr,
			Username:     user,
			Password:     pass,
			DB:           db,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
	}

	return &Client{
		rdb:        rdb,
		defaultTTL: defaultTTL,
	}
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	result := c.rdb.Get(ctx, key)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, errors.New(errors.NotFound)
		}
		return nil, result.Err()
	}
	return []byte(result.Val()), nil
}

func (c *Client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	return c.rdb.Set(ctx, key, value, ttl).Err()
}
