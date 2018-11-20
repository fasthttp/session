package session

import (
	"errors"
	"fmt"
)

var errNotSetProvider = errors.New("Not setted a session provider")
var errEmptySessionID = errors.New("Empty session id")

func errRegisterNilProvider(providerName string) error {
	return fmt.Errorf("The provider %s can not be nil", providerName)
}

func errProviderAlreadyRegisted(providerName string) error {
	return fmt.Errorf("The provider %s is already registered", providerName)
}
