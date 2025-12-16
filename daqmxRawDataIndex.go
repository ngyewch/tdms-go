package tdms

import (
	"fmt"
	"io"

	"github.com/samber/oops"
)

type DAQmxRawDataIndex struct {
	DataType       DataType
	ArrayDimension uint32
	ChunkSize      uint64
	Scalers        []Scaler
	RawDataWidths  []uint32
}

func (index *DAQmxRawDataIndex) GetDataType() DataType {
	return index.DataType
}

func (index *DAQmxRawDataIndex) GetArrayDimension() uint32 {
	return index.ArrayDimension
}

func (index *DAQmxRawDataIndex) GetChunkSize() uint64 {
	return index.ChunkSize
}

func (index *DAQmxRawDataIndex) GetTotalSizeInBytes() uint64 {
	var totalRawDataWidth uint64
	for _, rawDataWidth := range index.RawDataWidths {
		totalRawDataWidth += uint64(rawDataWidth)
	}
	return totalRawDataWidth * uint64(index.ArrayDimension) * index.ChunkSize
}

func (index *DAQmxRawDataIndex) PopulateScalers(scalers []Scaler) {
	for _, scaler := range scalers {
		for len(index.Scalers) <= int(scaler.ScaleId()) {
			index.Scalers = append(index.Scalers, nil)
		}
		index.Scalers[int(scaler.ScaleId())] = scaler
	}
}

func (index *DAQmxRawDataIndex) CheckCompatibility(rawDataIndex RawDataIndex) error {
	switch otherIndex := rawDataIndex.(type) {
	case *DAQmxRawDataIndex:
		if index.ChunkSize != otherIndex.ChunkSize {
			return fmt.Errorf("chunk size mismatch")
		}
		if len(index.RawDataWidths) != len(otherIndex.RawDataWidths) {
			return fmt.Errorf("raw data widths mismatch")
		}
		for i, rawDataWidth := range index.RawDataWidths {
			if rawDataWidth != otherIndex.RawDataWidths[i] {
				return fmt.Errorf("raw data widths mismatch")
			}
		}
		// TODO
		return nil
	default:
		return fmt.Errorf("incompatible raw data indexes [%T vs %T]", index, rawDataIndex)
	}
}

func ReadDAQmxRawDataIndex(r io.Reader, valueReader *ValueReader) (*DAQmxRawDataIndex, error) {
	var daqmxRawDataIndex DAQmxRawDataIndex
	var err error
	daqmxRawDataIndex.DataType, err = valueReader.ReadDataType(r)
	if err != nil {
		return nil, err
	}
	daqmxRawDataIndex.ArrayDimension, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	if daqmxRawDataIndex.ArrayDimension != 1 {
		return nil, oops.
			In("DAQmxRawDataIndex").
			With("arrayDimension", daqmxRawDataIndex.ArrayDimension).
			Errorf("invalid arrayDimension")
	}
	daqmxRawDataIndex.ChunkSize, err = valueReader.ReadU64(r)
	if err != nil {
		return nil, err
	}
	scalerVectorSize, err := valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	if scalerVectorSize == 0 {
		return nil, fmt.Errorf("no scalers specified")
	}
	for i := 0; i < int(scalerVectorSize); i++ {
		scaler, err := ReadDAQmxFormatChangingScaler(r, valueReader)
		if err != nil {
			return nil, err
		}
		for len(daqmxRawDataIndex.Scalers) <= int(scaler.ScaleId()) {
			daqmxRawDataIndex.Scalers = append(daqmxRawDataIndex.Scalers, nil)
		}
		daqmxRawDataIndex.Scalers[int(scaler.ScaleId())] = scaler
	}
	rawDataWidthVectorSize, err := valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(rawDataWidthVectorSize); i++ {
		rawDataWidth, err := valueReader.ReadU32(r)
		if err != nil {
			return nil, err
		}
		daqmxRawDataIndex.RawDataWidths = append(daqmxRawDataIndex.RawDataWidths, rawDataWidth)
	}
	return &daqmxRawDataIndex, nil
}
