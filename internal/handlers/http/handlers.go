package http

import (
	"net/http"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/info"
	"github.com/jkrus/Test_Simple_Multuplexor/internal/services"
	"github.com/jkrus/Test_Simple_Multuplexor/pkg/server"
)

const (
	apiV1 = "/api/v1"
)

type (
	handlers struct {
		r *http.ServeMux

		info *info.Handler
	}
)

func NewHandlers(service *services.Services) server.Handlers {
	return &handlers{
		r:    http.NewServeMux(),
		info: info.NewHandler(service.Info),
	}
}

func (h *handlers) Get() http.Handler {
	return h.r
}

func (h *handlers) Register() {
	r := http.NewServeMux()

	r.HandleFunc(apiV1, hello)

	h.info.Register(apiV1, r)

	h.r.Handle("/", r)

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from HTTP"))
}
