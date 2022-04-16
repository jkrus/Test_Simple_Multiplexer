package main

import (
	"log"
	"os"

	"github.com/jkrus/Test_Simple_Multuplexor/cmd/simple-multiplexor/app"
	"github.com/jkrus/Test_Simple_Multuplexor/internal/config"
)

func main() {
	args := os.Args[1:]
	ctx := app.NewContext()
	wg := app.NewWaitGroup()
	cfg := config.NewConfig()

	err := app.Start(ctx, args, wg, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
