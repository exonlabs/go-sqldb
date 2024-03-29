package sqldb

import (
	"errors"
	"fmt"
)

// Define common errors
var (
	ErrError     = errors.New("")
	ErrOpen      = fmt.Errorf("%wopen failed", ErrError)
	ErrClosed    = fmt.Errorf("%wconnection closed", ErrError)
	ErrBreak     = fmt.Errorf("%woperation break", ErrError)
	ErrTimeout   = fmt.Errorf("%woperation timeout", ErrError)
	ErrOperation = fmt.Errorf("%woperation error", ErrError)
)
