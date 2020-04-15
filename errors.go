package session

import (
	"errors"
)

var errNotSetProvider = errors.New("Not setted a session provider")
var errEmptySessionID = errors.New("Empty session id")
