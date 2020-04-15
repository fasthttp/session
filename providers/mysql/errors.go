package mysql

import "errors"

var errConfigHostEmpty = errors.New("Config Host must not be empty")
var errConfigPortZero = errors.New("Config Port must be more than 0")
