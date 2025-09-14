package tdms

import (
	"io"
)

type DAQmxFormatChangingScaler struct {
	DataType                     DataType
	RawBufferIndex               uint32
	RawByteOffsetWithinTheStride uint32
	SampleFormatBitmap           uint32
	ScaleId                      uint32
}

func readDAQmxFormatChangingScaler(r io.Reader, valueReader *ValueReader) (*DAQmxFormatChangingScaler, error) {
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
