package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/tdms-go"
	"github.com/urfave/cli/v3"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	inputFile := cmd.StringArg(inputFileArg.Name)

	tdmsFile, err := tdms.OpenFile(inputFile)
	if err != nil {
		return err
	}
	defer func(tdmsFile *tdms.File) {
		_ = tdmsFile.Close()
	}(tdmsFile)

	if false {
		fmt.Println("/")
		tdmsFile.Root().Properties().All()(func(name string, value any) bool {
			fmt.Printf("  * %s: %v [%T]\n", name, value, value)
			return true
		})
		for _, group := range tdmsFile.Root().Children() {
			fmt.Printf("- %s\n", group.Name())
			group.Properties().All()(func(name string, value any) bool {
				fmt.Printf("    * %s: %v [%T]\n", name, value, value)
				return true
			})
			for _, channel := range group.Children() {
				fmt.Printf("  - %s\n", channel.Name())
				channel.Properties().All()(func(name string, value any) bool {
					fmt.Printf("      * %s: %v [%T]\n", name, value, value)
					return true
				})
			}
		}
	}

	err = tdmsFile.ReadData()
	if err != nil {
		return err
	}

	return nil
}
