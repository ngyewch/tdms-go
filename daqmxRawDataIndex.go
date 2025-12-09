package tdms

import (
	"fmt"
	"github.com/samber/oops"
	"io"
)

type DAQmxRawDataIndex struct {
	DataType       DataType
	ArrayDimension uint32
	ChunkSize      uint64
	Scalers        []*DAQmxFormatChangingScaler
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
	for _, indexScaler := range index.Scalers {
		scalers[indexScaler.ScaleId] = indexScaler
	}
}

func (index *DAQmxRawDataIndex) IsCompatibleWith(rawDataIndex RawDataIndex) bool {
	switch otherIndex := rawDataIndex.(type) {
	case *DAQmxRawDataIndex:
		if len(index.Scalers) == 0 {

		}
		if len(otherIndex.Scalers) == 0 {

		}

		return true
	default:
		return false
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
		daqmxRawDataIndex.Scalers = append(daqmxRawDataIndex.Scalers, scaler)
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
