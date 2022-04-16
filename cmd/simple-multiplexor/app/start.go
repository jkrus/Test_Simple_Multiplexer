package app

import (
	"context"
	"flag"
	"fmt"
	"sync"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
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
	fmt.Println("START")
	return nil
}

func (s startCmd) Name() string {
	return s.name
}
