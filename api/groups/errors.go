package groups

import "errors"

// ErrNilDatabaseHandler signals that a nil database handler has been provided
var ErrNilDatabaseHandler = errors.New("nil database handler")
