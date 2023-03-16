package postgre

import "errors"

var (
	ErrConfigHostEmpty = errors.New("Config Host must not be empty")
	ErrConfigPortZero  = errors.New("Config Port must be more than 0")
)
