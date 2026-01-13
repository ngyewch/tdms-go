//go:build !(linux && amd64)

package converter

import "fmt"

func ConvertToNetCDF4(inputFile string, outputFile string) error {
	return fmt.Errorf("NetCDF converter not supported on this OS/platform")
}
