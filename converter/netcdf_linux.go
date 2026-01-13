//go:build linux && amd64

package converter

import (
	"github.com/fhs/go-netcdf/netcdf"
	"github.com/ngyewch/tdms-go"
)

func ConvertToNetCDF4(inputFile string, outputFile string) error {
	tdmsFile, err := tdms.OpenFile(inputFile)
	if err != nil {
		return err
	}
	defer func(tdmsFile *tdms.File) {
		_ = tdmsFile.Close()
	}(tdmsFile)

	datasetMap := make(map[string][]float64)
	channels := make([]*tdms.Node, 0)

	err = tdmsFile.ReadData(func(chunk tdms.Chunk) error {
		for _, channel := range chunk.Channels {
			values, exists := datasetMap[channel.Path]
			if !exists {
				channels = append(channels, channel.Node)
			}
			datasetMap[channel.Path] = append(values, channel.Samples...)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ncFile, err := netcdf.CreateFile(outputFile, netcdf.CLOBBER|netcdf.NETCDF4)
	if err != nil {
		return err
	}
	defer func(ncFile netcdf.Dataset) {
		_ = ncFile.Close()
	}(ncFile)

	for _, channel := range channels {
		variableName := normalizeNetCDFIdentifier(channel.Name())
		values := datasetMap[channel.Path()]
		dims := make([]netcdf.Dim, 1)
		dims[0], err = ncFile.AddDim(variableName, uint64(len(values)))
		if err != nil {
			return err
		}
		v, err := ncFile.AddVar(variableName, netcdf.DOUBLE, dims)
		if err != nil {
			return err
		}
		err = v.SetCompression(true, true, 9)
		if err != nil {
			return err
		}
		/*
			for propertyName, propertyValue := range channel.Properties().All() {
				convertedPropertyValue := func(propertyValue any) string {
					switch v := propertyValue.(type) {
					case time.Time:
						return v.Format(time.RFC3339)
					default:
						return fmt.Sprintf("%v", propertyValue)
					}
				}(propertyValue)
				attributeName := normalizeNetCDFIdentifier(propertyName)
				attribute := v.Attr(attributeName)
			}
		*/
		err = v.WriteFloat64s(values)
		if err != nil {
			return err
		}
	}

	return nil
}
