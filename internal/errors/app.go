package errors

import "github.com/pkg/errors"

const (
	errLoadConfigMsg       = "load config"
	errStartHTTPServerMsg  = "start http Server"
	errStartInfoServiceMsg = "start info service"
)

func ErrLoadConfig(w error) error {
	return errors.Wrap(w, errLoadConfigMsg)
}

func ErrStartHTTPServer(w error) error {
	return errors.Wrap(w, errStartHTTPServerMsg)
}

func ErrStartInfoService(w error) error {
	return errors.Wrap(w, errStartInfoServiceMsg)
}
