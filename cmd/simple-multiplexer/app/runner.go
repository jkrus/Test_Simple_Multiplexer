package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
)

type (
	Runner interface {
		Init([]string) error
		Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) error
		Name() string
	}
)

func Start(ctx context.Context, args []string, wg *sync.WaitGroup, cfg *config.Config) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Runner{
		newHelpCmd(),
		newInitCmd(),
		newStartCmd(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			err := cmd.Init(os.Args[2:])
			if err != nil {
				return err
			}
			return cmd.Run(ctx, wg, cfg)
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}
