package memcache

import "errors"

var errConfigServerListEmpty = errors.New("Config ServerList must not be empty")
var errConfigMaxIdleConnsZero = errors.New("Config MaxIdleConns must be more than 0")
var errExpirationIsTooBig = errors.New("Expiration duration should be a 32 bits value")
