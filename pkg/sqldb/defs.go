package sqldb

import (
	"errors"
	"fmt"

	"github.com/exonlabs/go-utils/pkg/types"
	"github.com/exonlabs/go-utils/pkg/xlog"
)

const SQL_PLACEHOLDER = "$?"

type Logger = xlog.Logger
type Options = types.Dict

type Data = types.Dict
type DataAdapter = func(any) (any, error)

// common errors
var (
	ErrError     = errors.New("")
	ErrOpen      = fmt.Errorf("%wopen failed", ErrError)
	ErrClosed    = fmt.Errorf("%wconnection closed", ErrError)
	ErrBreak     = fmt.Errorf("%woperation break", ErrError)
	ErrTimeout   = fmt.Errorf("%woperation timeout", ErrError)
	ErrOperation = fmt.Errorf("%woperation error", ErrError)
)
