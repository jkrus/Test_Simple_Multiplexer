package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/pkg/errors"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
	app_err "github.com/jkrus/Test_Simple_Multuplexor/internal/errors"
	app "github.com/jkrus/Test_Simple_Multuplexor/internal/handlers/http"
	"github.com/jkrus/Test_Simple_Multuplexor/pkg/server"
)

type (
	startCmd struct {
		fs   *flag.FlagSet
		name string
	}
)

func newStartCmd() Runner {
	sc := &startCmd{
		fs: flag.NewFlagSet("start", flag.ContinueOnError),
	}

	sc.fs.StringVar(&sc.name, "Start", "start", "use for start app")

	return sc
}

func (s startCmd) Init(args []string) error {
	return s.fs.Parse(args)
}

func (s startCmd) Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) error {
	err := cfg.Load()
	if err != nil {
		return app_err.ErrLoadConfig(err)
	}

	handlers := app.NewHandlers()

	newHTTP := server.NewHTTP(ctx, wg, cfg, handlers)

	handlers.Register()

	if err = newHTTP.Start(); err != nil {
		if !errors.Is(http.ErrServerClosed, err) {
			return app_err.ErrStartHTTPServer(err)
		}
	}

	<-ctx.Done()
	wg.Wait()
	log.Println("Application shutdown complete.")

	return nil
}

func (s startCmd) Name() string {
	return s.name
}
