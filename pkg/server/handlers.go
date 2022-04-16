package server

import "net/http"

type (
	Handlers interface {
		Register()
		Get() http.Handler
	}
)
