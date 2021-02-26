package rediscluster

import (
	"errors"
	"fmt"
)

var errConfigAddrsEmpty = errors.New("Config Addrs must not be empty")

func errRedisConnection(err error) error {
	return fmt.Errorf("Redis connection error: %v", err)
}
