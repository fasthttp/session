package redis

import (
	"errors"
	"fmt"
)

var errConfigAddrEmpty = errors.New("Config Addr must not be empty")

func errRedisConnection(err error) error {
	return fmt.Errorf("Redis connection error: %v", err)
}
