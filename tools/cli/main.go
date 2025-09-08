package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	version string

	app = &cli.Command{
		Name:    "tdms-cli",
		Usage:   "TDMS CLI",
		Version: version,
		Action:  nil,
		Commands: []*cli.Command{
			{
				Name:   "test",
				Usage:  "test",
				Action: doTest,
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
