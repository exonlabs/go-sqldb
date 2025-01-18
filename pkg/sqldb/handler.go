// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

// const (
// 	defaultConnectTimeout = float64(5)
// 	defaultRetryInterval  = float64(0.2)
// )

type Handler struct {
	log *Logger

	// database backend engine handler
	engine Engine

	// database connection params
	DbInfo *DbInfo
	// // session connection params
	// ConnectTimeout float64
	// RetryInterval  float64
}

func NewHandler(dbe Engine, log *Logger, opts Options) *Handler {
	h := &Handler{
		log:    log,
		engine: dbe,
		// ConnectTimeout: opts.GetFloat64(
		// 	"connect_timeout", defaultConnectTimeout),
		// RetryInterval: opts.GetFloat64(
		// 	"retry_interval", defaultRetryInterval),
	}
	// h.SetConfig(opts)
	return h
}

func (dbh *Handler) Engine() Engine {
	return dbh.engine
}

func (dbh *Handler) Session() *Session {
	return newSession(dbh)
}

// func (dbh *Handler) SetConfig(opts Options) error {

// 	return nil
// }
