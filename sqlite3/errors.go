package sqlite3

import "errors"

var errInvalidProviderConfig = errors.New("Invalid provider config")
var errConfigDBPathEmpty = errors.New("Config DBPath must not be empty")
