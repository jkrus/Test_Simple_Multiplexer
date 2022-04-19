package app

import (
	"context"
	"flag"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
)

type (
	helpCmd struct {
		fs   *flag.FlagSet
		name string
	}
)

func newHelpCmd() Runner {
	sc := &helpCmd{
		fs: flag.NewFlagSet("help", flag.ContinueOnError),
	}

	sc.fs.StringVar(&sc.name, "Help", "help", "use for help information")

	return sc
}

func (s *helpCmd) Init(args []string) error {
	return s.fs.Parse(args)

}

func (s *helpCmd) Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) error {
	fmt.Println("NAME:")
	fmt.Printf("	%v", config.AppUsage)
	fmt.Println()

	fmt.Println("USAGE:")
	fmt.Printf("	%v  command", config.AppName)
	fmt.Println()

	fmt.Println("VERSION:")
	fmt.Printf("	%v", version())
	fmt.Println()

	fmt.Println("COMMANDS::")
	fmt.Printf("	init	Initialize %v", config.AppUsage)
	fmt.Println()
	fmt.Printf("	start	Starts %v", config.AppUsage)
	fmt.Println()
	fmt.Println("	help	Shows a list of commands or help for one command")
	if err := cfg.Init(); err != nil {
		return errors.Wrap(err, "config init")
	}

	return nil
}

func (s *helpCmd) Name() string {
	return s.name
}
