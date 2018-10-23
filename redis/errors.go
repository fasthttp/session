package redis

import "errors"

var errConfigHostEmpty = errors.New("Config Host must not be empty")
var errConfigPortZero = errors.New("Config Port must not be more than 0")
var errConfigPoolSizeZero = errors.New("Config PoolSize must be more than 0")
var errConfigIdleTimeoutZero = errors.New("Config IdleTimeout must be more than 0")

func errRedisConnection(err error) error {
	return errors.New("Redis connection error: " + err.Error())
}
