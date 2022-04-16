package http

import (
	"net/http"

	"github.com/jkrus/Test_Simple_Multuplexor/pkg/server"
)

type (
	handlers struct {
		r *http.ServeMux
	}
)

func NewHandlers() server.Handlers {
	return &handlers{r: http.NewServeMux()}
}

func (h *handlers) Get() http.Handler {
	return h.r
}

func (h *handlers) Register() {
	r := http.NewServeMux()

	r.HandleFunc("/api/v1", hello)

	h.r.Handle("/", r)

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from HTTP"))
}
