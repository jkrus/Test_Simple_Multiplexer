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
	// Запускаю горутину для прослушивания канала в котором находятся ответы на запросы
	// по указанным urls
	go func() {
		exit := false
		for {
			select {
			case b, ok := <-ch:
				if !ok {
					exit = true
					break
				}
				// В случае, если на этапе запроса произошла ошибка,
				// возвращаю в ответе только url, при обработке которого она произошла и текст ошибки
				if b.err != nil {
					infoError := make([]models.Info, 0, len(urls))
					infoError = append(info, models.Info{URL: b.url, Body: []byte(b.err.Error())})
					info = infoError
					wg.Done()
					return
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
	// Запускаю горутину для отправки запросов
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
	// Добавлю группу ожидания для серии запросов
	wg := sync.WaitGroup{}
	// Создаю множество в котором сгруппирую входящие urls по 4
	set := make(map[int][]string)
	// Создам канал для отслеживания ошибки в каком - либо запросе.
	// Если при запросе произошла ошибка - закрою канал
	// Это будет сигналом к прекращению обработки списка urls
	errorChan := make(chan struct{})
	for idx, u := range urls {
		s := idx - idx%4
		set[s] = append(set[s], u)
	}

	errorFlag := false
	// Итерируюсь по множеству групп urls
	for _, valUrls := range set {
		// Для каждого url в группе выполню запрос и дождусь, пока не вернуться все запросы для данной группы.
		// После чего возьму следующую группу. Так я обеспечу не более 4-х исходящих запросов для каждого входящего
		for _, url := range valUrls {
			select {
			case _, ok := <-errorChan:
				if !ok {
					errorFlag = true
					break
				}
			default:
				wg.Add(1)
				go func(url string) {
					// Создам контекст с таймаутом для каждого запроса
					// Честно признаться не совсем уверен, что именно такой подход это go way.
					c, cancel := context.WithTimeout(ctx, 1*time.Second)
					defer wg.Done()
					defer cancel()

					request, err := http.NewRequest(http.MethodGet, url, nil)
					if err != nil {
						ch <- response{
							url:  url,
							body: nil,
							err:  err,
						}
						close(errorChan)
						return
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
						close(errorChan)
						return
					}

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						ch <- response{
							url:  url,
							body: nil,
							err:  err,
						}
						close(errorChan)
						return
					}

					defer resp.Body.Close()

					ch <- response{
						url:  url,
						body: body,
						err:  nil,
					}
				}(url)
			}
			if errorFlag {
				break
			}
		}

		wg.Wait()
		if errorFlag {
			break
		}
	}
	if !errorFlag {
		close(errorChan)
	}
	// Закрою канал, сообщив вызывающей горутине, что больше данных ожидать не стОит.
	close(ch)
	// После того как все urls будут отработаны - сообщу об этом вызывающей горутине
	waitGroup.Done()
}
