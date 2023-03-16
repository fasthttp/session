package session

import (
	"errors"
)

var (
	ErrNotSetProvider = errors.New("Not setted a session provider")
	ErrEmptySessionID = errors.New("Empty session id")
)
