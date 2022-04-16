package errors

import "github.com/pkg/errors"

const (
	errLoadConfigMsg      = "load config"
	errStartHTTPServerMsg = "start http Server"
)

func ErrLoadConfig(w error) error {
	return errors.Wrap(w, errLoadConfigMsg)
}

func ErrStartHTTPServer(w error) error {
	return errors.Wrap(w, errStartHTTPServerMsg)
}
