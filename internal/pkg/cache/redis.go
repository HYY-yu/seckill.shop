package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

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
	Set(ctx context.Context, key, value string, ttl time.Duration, options ...Option) error
	Get(ctx context.Context, key string, options ...Option) (string, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Expire(ctx context.Context, key string, ttl time.Duration) bool
	ExpireAt(ctx context.Context, key string, ttl time.Time) bool
	Del(ctx context.Context, key string, options ...Option) bool
	Exists(ctx context.Context, keys ...string) bool
	Incr(ctx context.Context, key string, options ...Option) int64
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "ping redis err")
	}

	return client, nil
}

// Set set some <key,value> into redis
func (c *cacheRepo) Set(ctx context.Context, key, value string, ttl time.Duration, options ...Option) error {
	var err error
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "set"
			opt.TraceRedis.Key = key
			opt.TraceRedis.TTL = ttl.Minutes()
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
			opt.TraceRedis.Err = err

			addTracing(ctx, opt.TraceRedis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if err = c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		err = errors.Wrapf(err, "redis set key: %s err", key)
	}
	return err
}

// Get get some key from redis
func (c *cacheRepo) Get(ctx context.Context, key string, options ...Option) (string, error) {
	var err error
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "get"
			opt.TraceRedis.Key = key
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
			opt.TraceRedis.Err = err

			addTracing(ctx, opt.TraceRedis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		err = errors.Wrapf(err, "redis get key: %s err", key)
	}
	return value, err
}

// TTL get some key from redis
func (c *cacheRepo) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return -1, errors.Wrapf(err, "redis get key: %s err", key)
	}

	return ttl, nil
}

// Expire expire some key
func (c *cacheRepo) Expire(ctx context.Context, key string, ttl time.Duration) bool {
	ok, _ := c.client.Expire(ctx, key, ttl).Result()
	return ok
}

// ExpireAt expire some key at some time
func (c *cacheRepo) ExpireAt(ctx context.Context, key string, ttl time.Time) bool {
	ok, _ := c.client.ExpireAt(ctx, key, ttl).Result()
	return ok
}

func (c *cacheRepo) Exists(ctx context.Context, keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	value, _ := c.client.Exists(ctx, keys...).Result()
	return value > 0
}

func (c *cacheRepo) Del(ctx context.Context, key string, options ...Option) bool {
	var err error
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "del"
			opt.TraceRedis.Key = key
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
			opt.TraceRedis.Err = err

			addTracing(ctx, opt.TraceRedis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if key == "" {
		return true
	}

	value, err := c.client.Del(ctx, key).Result()
	if err != nil {
		err = errors.Wrapf(err, "redis del key: %s err", key)
	}
	return value > 0
}

func (c *cacheRepo) Incr(ctx context.Context, key string, options ...Option) int64 {
	var err error

	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.TraceRedis != nil {
			opt.TraceRedis.Timestamp = time_parse.CSTLayoutString()
			opt.TraceRedis.Handle = "incr"
			opt.TraceRedis.Key = key
			opt.TraceRedis.CostSeconds = time.Since(ts).Seconds()
			opt.TraceRedis.Err = err

			addTracing(ctx, opt.TraceRedis)
		}
	}()

	for _, f := range options {
		f(opt)
	}
	value, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		err = errors.Wrapf(err, "redis Incr key: %s err", key)
	}
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
