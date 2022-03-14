package cache

import (
	"context"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var (
	cache       Cache
	maxIdle     = 20
	idleTimeout = 120 * time.Second
	maxActive   = 100
)

type CacheOption struct {
	Host        string
	Auth        bool
	Username    string
	Password    string
	Database    string
	Table       string
	DB          int
	MaxIdle     int
	IdleTimeout time.Duration
	MaxActive   int
}

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) (interface{}, error)
	Overdue(ctx context.Context, key interface{}) bool
	SetEx(ctx context.Context, key string, value interface{}, sec int) (interface{}, error)
}

type CacheBaseRedis struct {
	Pool *redigo.Pool
}

func CacheFactory(name string, opt CacheOption) {
	switch name {
	case "redis":
		cache = NewCacheBaseRedis(opt)
	default:
		cache = NewCacheBaseRedis(opt)
	}
}

// new cache base redis
func NewCacheBaseRedis(opt CacheOption) Cache {
	applyOption(opt)
	pool := &redigo.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		MaxActive:   maxActive,
		Wait:        true,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", opt.Host)
			if err != nil {
				return nil, err
			}
			if opt.Auth {
				if _, err := c.Do("AUTH", opt.Password); err != nil {
					c.Close()
					return nil, err
				}
				if _, err := c.Do("SELECT", opt.DB); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &CacheBaseRedis{
		Pool: pool,
	}
}

func applyOption(opt CacheOption) {
	if opt.MaxIdle > 0 {
		maxIdle = opt.MaxIdle
	}
	if opt.IdleTimeout > 0 {
		idleTimeout = opt.IdleTimeout
	}
	if opt.MaxActive > 0 {
		maxActive = opt.MaxActive
	}
}

func (c *CacheBaseRedis) Get(ctx context.Context, key string) (interface{}, error) {
	r := c.Pool.Get()
	defer r.Close()
	return r.Do("GET", key)
}

func (c *CacheBaseRedis) Set(ctx context.Context, key string, value interface{}) (interface{}, error) {
	r := c.Pool.Get()
	defer r.Close()
	return r.Do("SET", key, value)
}

func (c *CacheBaseRedis) SetEx(ctx context.Context, key string, value interface{}, sec int) (interface{}, error) {
	r := c.Pool.Get()
	defer r.Close()
	return r.Do("SETEX", key, sec, value)
}

func (c *CacheBaseRedis) Overdue(ctx context.Context, key interface{}) bool {
	r := c.Pool.Get()
	defer r.Close()
	res, err := r.Do("TTL", key)
	if err != nil {
		return false
	}
	return res.(int64) > 1
}

func Get(ctx context.Context, key string) (interface{}, error) {
	return cache.Get(ctx, key)
}

func Set(ctx context.Context, key string, value interface{}) (interface{}, error) {
	return cache.Set(ctx, key, value)
}

func SetEx(ctx context.Context, key string, value interface{}, sec int) (interface{}, error) {
	return cache.SetEx(ctx, key, value, sec)
}

func Overdue(ctx context.Context, key interface{}) bool {
	return cache.Overdue(ctx, key)
}
