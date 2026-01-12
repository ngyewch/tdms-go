package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ngyewch/tdms-go/converter"
	"github.com/urfave/cli/v3"
)

func doConvert(ctx context.Context, cmd *cli.Command) error {
	inputFile := cmd.StringArg(inputFileArg.Name)
	outputFile := cmd.StringArg(outputFileArg.Name)

	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}
	if outputFile == "" {
		return fmt.Errorf("output file is required")
	}

	outputExtension := filepath.Ext(outputFile)
	switch outputExtension {
	case ".mat":
		return doConvertToMAT(ctx, inputFile, outputFile)
	case ".h5", ".hdf5":
		return doConvertToHDF5(ctx, inputFile, outputFile)
	case ".cdl":
		return doConvertToCDL(ctx, inputFile, outputFile)
	default:
		return fmt.Errorf("unsupported output file extension")
	}
}

func doConvertToHDF5(ctx context.Context, inputFile string, outputFile string) error {
	return converter.ConvertToHDF5(inputFile, outputFile)
}

func doConvertToMAT(ctx context.Context, inputFile string, outputFile string) error {
	return converter.ConvertToMAT(inputFile, outputFile)
}

func doConvertToCDL(ctx context.Context, inputFile string, outputFile string) error {
	return converter.ConvertToCDL(inputFile, outputFile)
}
