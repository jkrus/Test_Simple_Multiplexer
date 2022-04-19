package info

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
	"github.com/jkrus/Test_Simple_Multuplexor/pkg/models"
	"github.com/jkrus/Test_Simple_Multuplexor/pkg/service"
)

type (
	Service interface {
		service.Service
		AddToWaitGroup()
		DellFromWaitGroup()
		GetInfoByURLS([]string) []models.Info
	}

	infoService struct {
		ctx    context.Context
		mainWG *sync.WaitGroup
		cfg    *config.Config

		wg *sync.WaitGroup
	}
)

func NewInfoService(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) Service {

	return &infoService{ctx: ctx, mainWG: wg, cfg: cfg, wg: &sync.WaitGroup{}}
}

type response struct {
	url  string
	body []byte
	err  error
}

func (is *infoService) AddToWaitGroup() {
	is.wg.Add(1)
}

func (is *infoService) DellFromWaitGroup() {
	is.wg.Done()
}

func (is *infoService) GetInfoByURLS(urls []string) []models.Info {
	info := make([]models.Info, 0, len(urls))
	ch := make(chan response, 20)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		exit := false
		for {
			select {
			case b, ok := <-ch:
				if !ok {
					exit = true
					break
				}
				if b.err != nil {
					info = append(info, models.Info{URL: b.url, Body: []byte(b.err.Error())})
					break
				}

				info = append(info, models.Info{URL: b.url, Body: b.body})
			}
			if exit {
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go send(is.ctx, &wg, urls, ch)
	wg.Wait()

	return info
}

func (is *infoService) Start() error {
	log.Println("Start Info service...")

	is.createHandlerContext()

	log.Println("Info service start success.")

	return nil
}

func (is *infoService) Stop() error {
	log.Println("Stop Info Service...")

	is.wg.Wait()

	log.Println("Info service stopped.")

	return nil
}

func (is *infoService) createHandlerContext() {
	is.mainWG.Add(1)
	go func() {
		for {
			<-is.ctx.Done()
			_ = is.Stop()
			is.mainWG.Done()
			return
		}
	}()

}

func send(ctx context.Context, waitGroup *sync.WaitGroup, urls []string, ch chan<- response) {
	c, cancel := context.WithTimeout(ctx, 1*time.Second)

	wg := sync.WaitGroup{}

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			defer cancel()

			request, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				ch <- response{
					url:  url,
					body: nil,
					err:  err,
				}
			}

			request.WithContext(c)
			client := http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				ch <- response{
					url:  url,
					body: nil,
					err:  err,
				}
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ch <- response{
					url:  url,
					body: nil,
					err:  err,
				}
			}

			defer resp.Body.Close()

			ch <- response{
				url:  url,
				body: body,
				err:  nil,
			}
		}(u)
	}

	wg.Wait()
	close(ch)
	cancel()
	waitGroup.Done()
}
