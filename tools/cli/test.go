package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ngyewch/tdms-go"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	f, err := os.Open(cmd.Args().First())
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	reader := tdms.NewReader(f)
	leadIn, err := reader.ReadLeadIn()
	if err != nil {
		return err
	}
	fmt.Printf("leadIn: %+v\n", leadIn)
	err = reader.ReadMetaData(leadIn)
	if err != nil {
		return err
	}

	return nil
}
