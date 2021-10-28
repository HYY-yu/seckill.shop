package cache

import (
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/HYY-yu/seckill/internal/service/config"
	"github.com/HYY-yu/seckill/pkg/errors"
	"github.com/HYY-yu/seckill/pkg/time_parse"
)

type Option func(*option)

type option struct {
	TraceRedis *Redis
}

func newOption() *option {
	return &option{}
}

var _ Repo = (*cacheRepo)(nil)

type Repo interface {
	i()
	Set(key, value string, ttl time.Duration, options ...Option) error
	Get(key string, options ...Option) (string, error)
	TTL(key string) (time.Duration, error)
	Expire(key string, ttl time.Duration) bool
	ExpireAt(key string, ttl time.Time) bool
	Del(key string, options ...Option) bool
	Exists(keys ...string) bool
	Incr(key string, options ...Option) int64
	Close() error
}

type cacheRepo struct {
	client *redis.Client
}

func New() (Repo, error) {
	client, err := redisConnect()
	if err != nil {
		return nil, err
	}

	return &cacheRepo{
		client: client,
	}, nil
}

func (c *cacheRepo) i() {}

func redisConnect() (*redis.Client, error) {
	cfg := config.Get().Redis
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           cfg.Db,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrap(err, "ping redis err")
	}

	return client, nil
}

// Set set some <key,value> into redis
func (c *cacheRepo) Set(key, value string, ttl time.Duration, options ...Option) error {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "set"
			opt.TraceRedis.Key = key
			opt.TraceRedis.Value = value
			opt.TraceRedis.TTL = ttl.Minutes()
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if err := c.client.Set(key, value, ttl).Err(); err != nil {
		return errors.Wrapf(err, "redis set key: %s err", key)
	}

	return nil
}

// Get get some key from redis
func (c *cacheRepo) Get(key string, options ...Option) (string, error) {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "get"
			opt.TraceRedis.Key = key
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
		}
	}()

	for _, f := range options {
		f(opt)
	}

	value, err := c.client.Get(key).Result()
	if err != nil {
		return "", errors.Wrapf(err, "redis get key: %s err", key)
	}

	return value, nil
}

// TTL get some key from redis
func (c *cacheRepo) TTL(key string) (time.Duration, error) {
	ttl, err := c.client.TTL(key).Result()
	if err != nil {
		return -1, errors.Wrapf(err, "redis get key: %s err", key)
	}

	return ttl, nil
}

// Expire expire some key
func (c *cacheRepo) Expire(key string, ttl time.Duration) bool {
	ok, _ := c.client.Expire(key, ttl).Result()
	return ok
}

// ExpireAt expire some key at some time
func (c *cacheRepo) ExpireAt(key string, ttl time.Time) bool {
	ok, _ := c.client.ExpireAt(key, ttl).Result()
	return ok
}

func (c *cacheRepo) Exists(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	value, _ := c.client.Exists(keys...).Result()
	return value > 0
}

func (c *cacheRepo) Del(key string, options ...Option) bool {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "del"
			opt.TraceRedis.Key = key
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if key == "" {
		return true
	}

	value, _ := c.client.Del(key).Result()
	return value > 0
}

func (c *cacheRepo) Incr(key string, options ...Option) int64 {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "incr"
			opt.TraceRedis.Key = key
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
		}
	}()

	for _, f := range options {
		f(opt)
	}
	value, _ := c.client.Incr(key).Result()
	return value
}

// Close close redis client
func (c *cacheRepo) Close() error {
	return c.client.Close()
}

// WithTrace 设置trace信息
func WithTrace() Option {
	return func(opt *option) {
		opt.TraceRedis = new(Redis)
	}
}
