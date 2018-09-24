package memcache

import "errors"

var errInvalidProviderConfig = errors.New("Invalid provider config")
var errConfigServerListEmpty = errors.New("Config ServerList must not be empty")
var errConfigMaxIdleConnsZero = errors.New("Config MaxIdleConns must be more than 0")
