package tdms

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
	isRawDataIndex() bool
}

type DefaultRawDataIndex struct {
	DataType         DataType
	ArrayDimension   uint32
	ChunkSize        uint64
	TotalSizeInBytes uint64
}

func (index DefaultRawDataIndex) isRawDataIndex() bool {
	return true
}

func (index DefaultRawDataIndex) GetDataType() DataType {
	return index.DataType
}

func (index DefaultRawDataIndex) GetArrayDimension() uint32 {
	return index.ArrayDimension
}

func (index DefaultRawDataIndex) GetChunkSize() uint64 {
	return index.ChunkSize
}
