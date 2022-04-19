package services

import (
	"context"
	"sync"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
	"github.com/jkrus/Test_Simple_Multuplexor/internal/errors"
	"github.com/jkrus/Test_Simple_Multuplexor/internal/info"
)

type (
	Services struct {
		Info info.Service
	}
)

func NewServices(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) (*Services, error) {
	// provide Info Service.
	infoService := info.NewInfoService(ctx, wg, cfg)
	if err := infoService.Start(); err != nil {
		return nil, errors.ErrStartInfoService(err)
	}
	return &Services{
		Info: infoService,
	}, nil

}
