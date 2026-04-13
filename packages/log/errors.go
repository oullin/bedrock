package log

import "errors"

var (
	ErrInvalidLevel      = errors.New("log: invalid log level")
	ErrUnsupportedDriver = errors.New("log: unsupported driver")
	ErrChannelNotFound   = errors.New("log: channel not found in configuration")
	ErrMissingPath       = errors.New("log: file path is required")
	ErrHandlerClosed     = errors.New("log: handler is closed")
	ErrNoDispatcher      = errors.New("log: no event dispatcher set")
)
