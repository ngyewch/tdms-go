package tdms

import (
	"io"
)

type DAQmxRawDataIndex struct {
	DataType       DataType
	ArrayDimension uint32
	ChunkSize      uint64
	Scalers        []*DAQmxFormatChangingScaler
	RawDataWidths  []uint32
}

func (index DAQmxRawDataIndex) isRawDataIndex() bool {
	return true
}

func (index DAQmxRawDataIndex) GetDataType() DataType {
	return index.DataType
}

func (index DAQmxRawDataIndex) GetArrayDimension() uint32 {
	return index.ArrayDimension
}

func (index DAQmxRawDataIndex) GetChunkSize() uint64 {
	return index.ChunkSize
}

func (index DAQmxRawDataIndex) GetTotalSizeInBytes() uint64 {
	var totalRawDataWidth uint64
	for _, rawDataWidth := range index.RawDataWidths {
		totalRawDataWidth += uint64(rawDataWidth)
	}
	return totalRawDataWidth * uint64(index.ArrayDimension) * index.ChunkSize
}

type DAQmxFormatChangingScaler struct {
	DataType                     DataType
	RawBufferIndex               uint32
	RawByteOffsetWithinTheStride uint32
	SampleFormatBitmap           uint32
	ScaleId                      uint32
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
	daqmxRawDataIndex.ChunkSize, err = valueReader.ReadU64(r)
	if err != nil {
		return nil, err
	}
	scalerVectorSize, err := valueReader.ReadU32(r)
	if err != nil {
		return nil, err
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

func ReadDAQmxFormatChangingScaler(r io.Reader, valueReader *ValueReader) (*DAQmxFormatChangingScaler, error) {
	var scaler DAQmxFormatChangingScaler
	var err error
	scaler.DataType, err = valueReader.ReadDataType(r)
	if err != nil {
		return nil, err
	}
	scaler.RawBufferIndex, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.RawByteOffsetWithinTheStride, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.SampleFormatBitmap, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	scaler.ScaleId, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	return &scaler, nil
}
