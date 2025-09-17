package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/goforj/godump"
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
	for {
		segment, err := reader.NextSegment()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if segment == nil {
			break
		}
		godump.Dump(segment)
		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		fmt.Printf("offset: %d [0x%x]\n", pos, pos)
	}

	return nil
}
