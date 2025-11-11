package main

import (
	"context"

	"github.com/goforj/godump"
	"github.com/ngyewch/tdms-go"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	inputFile := cmd.StringArg(inputFileArg.Name)

	tdmsFile, err := tdms.OpenFile(inputFile)
	if err != nil {
		return err
	}

	godump.Dump(tdmsFile.Segments())

	return nil
}
