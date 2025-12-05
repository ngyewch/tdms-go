package tdms

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
