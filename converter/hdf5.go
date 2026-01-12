package converter

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/ngyewch/tdms-go"
	"github.com/scigolib/hdf5"
)

func ConvertToHDF5(inputFile string, outputFile string) error {
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

	hdf5File, err := hdf5.CreateForWrite(outputFile, hdf5.CreateTruncate)
	if err != nil {
		return err
	}
	defer func(hdf5File *hdf5.FileWriter) {
		_ = hdf5File.Close()
	}(hdf5File)

	err = createHDF5Groups(hdf5File, tdmsFile.Root())
	if err != nil {
		return err
	}

	for _, channel := range channels {
		values := datasetMap[channel.Path()]
		hdf5Path, err := convertTDMSPathToHDFS5Path(channel.Path())
		if err != nil {
			return err
		}
		dataset, err := hdf5File.CreateDataset(hdf5Path, hdf5.Float64, []uint64{uint64(len(values))})
		if err != nil {
			return err
		}
		for propertyName, propertyValue := range channel.Properties().All() {
			switch v := propertyValue.(type) {
			case time.Time:
				convertedPropertyValue := v.Format(time.RFC3339)
				err = dataset.WriteAttribute(propertyName, convertedPropertyValue)
				if err != nil {
					return err
				}
			default:
				err = dataset.WriteAttribute(propertyName, propertyValue)
				if err != nil {
					return err
				}
			}
		}
		err = dataset.Write(values)
		if err != nil {
			return err
		}
	}

	return nil
}

func createHDF5Groups(hdf5File *hdf5.FileWriter, node *tdms.Node) error {
	for _, childNode := range node.Children() {
		if len(childNode.Children()) > 0 {
			hdf5Path, err := convertTDMSPathToHDFS5Path(childNode.Path())
			if err != nil {
				return err
			}
			group, err := hdf5File.CreateGroup(hdf5Path)
			if err != nil {
				return err
			}
			for propertyName, propertyValue := range childNode.Properties().All() {
				err = group.WriteAttribute(propertyName, propertyValue)
				if err != nil {
					return err
				}
			}
			err = createHDF5Groups(hdf5File, childNode)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func convertTDMSPathToHDFS5Path(path string) (string, error) {
	tdmsPath, err := tdms.ObjectPathFromString(path)
	if err != nil {
		return "", err
	}

	if tdmsPath.IsRoot() {
		return "/", nil
	} else if tdmsPath.IsGroup() {
		return "/" + slug.Make(tdmsPath.Group), nil
	} else if tdmsPath.IsChannel() {
		return "/" + slug.Make(tdmsPath.Group) + "/" + slug.Make(tdmsPath.Channel), nil
	}

	return "", fmt.Errorf("invalid path")
}
