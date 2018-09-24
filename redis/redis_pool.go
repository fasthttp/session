package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func newRedisPool(config *Config) *redis.Pool {
	server := fmt.Sprintf("%s:%d", config.Host, config.Port)
	timeout := time.Duration(config.IdleTimeout) * time.Second

	return &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: timeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if config.Password != "" {
				if _, err := c.Do("AUTH", config.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", config.DbNumber); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
