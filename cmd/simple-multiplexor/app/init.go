package app

import (
	"context"
	"flag"
	"sync"

	"github.com/pkg/errors"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
)

type (
	initCmd struct {
		fs   *flag.FlagSet
		name string
	}
)

func newInitCmd() Runner {
	sc := &initCmd{
		fs: flag.NewFlagSet("init", flag.ContinueOnError),
	}

	sc.fs.StringVar(&sc.name, "Init", "init", "use for init config params & config dir")

	return sc
}

func (s initCmd) Init(args []string) error {
	return s.fs.Parse(args)

}

func (s initCmd) Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) error {
	if err := cfg.Init(); err != nil {
		return errors.Wrap(err, "config init")
	}

	return nil
}

func (s initCmd) Name() string {
	return s.name
}
