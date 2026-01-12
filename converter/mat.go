package converter

import (
	"github.com/gosimple/slug"
	"github.com/ngyewch/tdms-go"
	"github.com/scigolib/matlab"
	"github.com/scigolib/matlab/types"
)

func ConvertToMAT(inputFile string, outputFile string) error {
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

	matFile, err := matlab.Create(outputFile, matlab.Version73)
	if err != nil {
		return err
	}
	defer func(hdf5File *matlab.MatFileWriter) {
		_ = hdf5File.Close()
	}(matFile)

	for _, channel := range channels {
		values := datasetMap[channel.Path()]
		attributes := make(map[string]any)
		for propertyName, propertyValue := range channel.Properties().All() {
			/*
				switch v := propertyValue.(type) {
					case time.Time:
						convertedPropertyValue := v.Format(time.RFC3339)
						attributes[propertyName] = convertedPropertyValue
				default:
					attributes[propertyName] = propertyValue
				}
			*/
			attributes[propertyName] = propertyValue
		}
		err = matFile.WriteVariable(&types.Variable{
			Name:       slug.Make(channel.Name()),
			Dimensions: []int{len(values)},
			DataType:   types.Double,
			Data:       values,
			Attributes: attributes,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
