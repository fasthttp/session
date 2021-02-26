package redisfailover

import (
	"errors"
	"fmt"
)

var errConfigMasterNameEmpty = errors.New("Config MasterName must not be empty")

func errRedisConnection(err error) error {
	return fmt.Errorf("Redis connection error: %v", err)
}
