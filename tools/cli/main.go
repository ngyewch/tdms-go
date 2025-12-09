package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
)

var (
	version string

	inputFileArg = &cli.StringArg{
		Name:      "input-file",
		UsageText: "(input file)",
	}

	app = &cli.Command{
		Name:    "tdms-cli",
		Usage:   "TDMS CLI",
		Version: version,
		Action:  nil,
		Commands: []*cli.Command{
			{
				Name:  "test",
				Usage: "test",
				Arguments: []cli.Argument{
					inputFileArg,
				},
				Action: doTest,
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		oopsError, ok := oops.AsOops(err)
		if ok {
			fmt.Printf("%+v", oopsError)
		} else {
			log.Fatal(err)
		}
	}
}
