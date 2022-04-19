package info

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	infoURL = "/info"
)

type (
	Handler struct {
		info Service
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{
		info: service,
	}
}

func (h *Handler) Register(url string, mux *http.ServeMux) {
	r := http.NewServeMux()

	r.HandleFunc(url+infoURL, h.getInfo)

	mux.Handle("/", r)
}

func (h *Handler) getInfo(w http.ResponseWriter, r *http.Request) {
	h.info.AddToWaitGroup()
	defer h.info.DellFromWaitGroup()

	urls, err := prepareRequest(r)
	log.Println(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	resp := h.info.GetInfoByURLS(urls)
	marshal, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}
