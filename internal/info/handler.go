package info

import (
	"encoding/json"
	"net/http"

	"github.com/jkrus/Test_Simple_Multuplexor/pkg/models"
)

const (
	infoURL = "/info"
)

type (
	Handler struct {
		info    Service
		chanReq chan requestData
	}

	requestData struct {
		req          *http.Request
		resp         http.ResponseWriter
		responseData chan responseData
	}

	responseData struct {
		err  error
		data []models.Info
	}
)

func NewHandler(service Service) *Handler {
	return &Handler{
		info:    service,
		chanReq: make(chan requestData, 100),
	}
}

func (h *Handler) Register(url string, mux *http.ServeMux) {
	r := http.NewServeMux()

	r.HandleFunc(url+infoURL, h.getInfo)

	mux.Handle("/", r)

	go h.getInfoWorker()
}

func (h *Handler) getInfo(w http.ResponseWriter, r *http.Request) {
	h.info.AddToWaitGroup()
	defer h.info.DellFromWaitGroup()

	info := make(chan responseData)
	h.chanReq <- requestData{resp: w, req: r, responseData: info}
	resp := <-info

	if resp.err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp.err.Error()))
	}

	marshal, err := json.Marshal(resp.data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}

// Наверное не совсем такое решение ожидалось для задачи одновременного
// обслуживания не более 100 входящих запросов.
// Это воркер, который стартует при регистрации обработчика
// При регистрации обработчика создан буферизированный канал из 100 элементов.
// Элементом канала является структура в которую я помещаю ResponseWriter, Request
// и канал для данных.
// Поскольку при graceful shutdown сервис будет ожидать дообработки всех запросов, то
// данная горутина завершится без потери данных и завершится вместе с основной горутиной.
func (h *Handler) getInfoWorker() {
	for {
		select {
		case r := <-h.chanReq:
			urls, err := prepareRequest(r.req)
			if err != nil {
				r.responseData <- responseData{err: err, data: nil}
				close(r.responseData)
			}

			resp := h.info.GetInfoByURLS(urls)
			r.responseData <- responseData{err: nil, data: resp}
			close(r.responseData)
		}
	}
}
