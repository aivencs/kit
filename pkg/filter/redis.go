package filter

import (
	"context"
	"time"

	redisbloom "github.com/RedisBloom/redisbloom-go"
	redigo "github.com/gomodule/redigo/redis"
)

var (
	filter      Filter
	maxIdle     = 20
	idleTimeout = 240 * time.Second
	maxActive   = 100
)

type Filter interface {
	Exist(ctx context.Context, val string) (bool, error)
	Add(ctx context.Context, val string) (bool, error)
}

type LinkFilterBaseRedis struct {
	Client *redisbloom.Client
	Pool   *redigo.Pool
	Key    string
}

type FilterOption struct {
	Host        string
	Auth        bool
	Username    string
	Password    string
	Database    string
	Table       string
	DB          int
	Key         string
	MaxIdle     int
	IdleTimeout time.Duration
	MaxActive   int
}

func applyOption(opt FilterOption) {
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

func FilterFactory(name string, opt FilterOption) {
	switch name {
	case "redis":
		filter = NewRedisFilter(opt)
	default:
		filter = NewRedisFilter(opt)
	}
}

// new filter base redis
func NewRedisFilter(opt FilterOption) Filter {
	applyOption(opt)
	rdp := &redigo.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		MaxActive:   maxActive,
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
	rbc := redisbloom.NewClientFromPool(rdp, opt.Key)
	return &LinkFilterBaseRedis{
		Pool:   rdp,
		Client: rbc,
		Key:    opt.Key,
	}
}

/*
function of redis filter
*/

func (c *LinkFilterBaseRedis) Exist(ctx context.Context, val string) (bool, error) {
	return c.Client.Exists(c.Key, val)
}

func (c *LinkFilterBaseRedis) Add(ctx context.Context, val string) (bool, error) {
	return c.Client.Add(c.Key, val)
}

/*
for caller
*/

func Exist(ctx context.Context, val string) (bool, error) {
	return filter.Exist(ctx, val)
}

func Add(ctx context.Context, val string) (bool, error) {
	return filter.Add(ctx, val)
}
