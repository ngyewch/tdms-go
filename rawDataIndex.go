package tdms

import (
	"io"
)

const (
	RawDataIndexTypeSameAsPreviousSegment         = 0x00000000
	RawDataIndexTypeNoRawData                     = 0xffffffff
	RawDataIndexTypeDAQmxFormatChangingScalerType = 0x00001269
	RawDataIndexTypeDAQmxDigitalLineScalerType    = 0x0000126a
)

type RawDataIndex interface {
	GetDataType() DataType
	GetArrayDimension() uint32
	GetChunkSize() uint64
	GetTotalSizeInBytes() uint64
	IsCompatibleWith(rawDataIndex RawDataIndex) bool
	PopulateScalers(scalers []Scaler)
}

type DefaultRawDataIndex struct {
	DataType         DataType
	ArrayDimension   uint32
	ChunkSize        uint64
	TotalSizeInBytes uint64
}

func (index *DefaultRawDataIndex) GetDataType() DataType {
	return index.DataType
}

func (index *DefaultRawDataIndex) GetArrayDimension() uint32 {
	return index.ArrayDimension
}

func (index *DefaultRawDataIndex) GetChunkSize() uint64 {
	return index.ChunkSize
}

func (index *DefaultRawDataIndex) GetTotalSizeInBytes() uint64 {
	if index.TotalSizeInBytes > 0 {
		return index.TotalSizeInBytes
	}
	return uint64(index.DataType.SizeOf()) * uint64(index.ArrayDimension) * index.ChunkSize
}

func (index *DefaultRawDataIndex) PopulateScalers(scalers []Scaler) {
	// do nothing
}

func (index *DefaultRawDataIndex) IsCompatibleWith(rawDataIndex RawDataIndex) bool {
	switch rawDataIndex.(type) {
	case *DefaultRawDataIndex:
		return true
	default:
		return false
	}
}

func ReadDefaultRawDataIndex(r io.Reader, valueReader *ValueReader) (*DefaultRawDataIndex, error) {
	var defaultRawDataIndex DefaultRawDataIndex
	var err error
	defaultRawDataIndex.DataType, err = valueReader.ReadDataType(r)
	if err != nil {
		return nil, err
	}
	defaultRawDataIndex.ArrayDimension, err = valueReader.ReadU32(r)
	if err != nil {
		return nil, err
	}
	defaultRawDataIndex.ChunkSize, err = valueReader.ReadU64(r)
	if err != nil {
		return nil, err
	}
	// TODO confirm
	if defaultRawDataIndex.DataType.SizeOf() <= 0 {
		defaultRawDataIndex.TotalSizeInBytes, err = valueReader.ReadU64(r)
		if err != nil {
			return nil, err
		}
	}
	return &defaultRawDataIndex, nil
}
