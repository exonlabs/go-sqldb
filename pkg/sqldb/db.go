package sqldb

import (
	"github.com/exonlabs/go-utils/pkg/types"
	"github.com/exonlabs/go-utils/pkg/xlog"
)

const (
	defaultConnectTimeout = float64(5)
	defaultRetryInterval  = float64(0.2)
)

type Options = types.NDict

type Handler struct {
	Logger *xlog.Logger

	// database backend engine handler
	engine Engine

	// session connection params
	ConnectTimeout float64
	RetryInterval  float64
}

func NewHandler(dbe Engine, opts Options, logger *xlog.Logger) *Handler {
	return &Handler{
		Logger: logger,
		engine: dbe,
		ConnectTimeout: opts.GetFloat64(
			"connect_timeout", defaultConnectTimeout),
		RetryInterval: opts.GetFloat64(
			"retry_interval", defaultRetryInterval),
	}
}

func (dbh *Handler) Session() *Session {
	return newSession(dbh)
}
