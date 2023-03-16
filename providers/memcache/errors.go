package memcache

import "errors"

var (
	ErrConfigServerListEmpty  = errors.New("Config ServerList must not be empty")
	ErrConfigMaxIdleConnsZero = errors.New("Config MaxIdleConns must be more than 0")
	ErrExpirationIsTooBig     = errors.New("Expiration duration should be a 32 bits value")
)
